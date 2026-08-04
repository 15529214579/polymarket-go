package order

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type preparedOrderObserver func(string) error

type preparedOrderObserverKey struct{}

func withPreparedOrderObserver(ctx context.Context, observer preparedOrderObserver) context.Context {
	return context.WithValue(ctx, preparedOrderObserverKey{}, observer)
}

func notifyPreparedOrder(ctx context.Context, orderID string) error {
	observer, _ := ctx.Value(preparedOrderObserverKey{}).(preparedOrderObserver)
	if observer == nil {
		return nil
	}
	return observer(orderID)
}

type ExecutionRecord struct {
	ID        string
	Mode      string
	Status    Status
	OrderID   string
	Intent    Intent
	Result    Result
	Applied   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ExecutionLedger struct {
	db *sql.DB
}

func OpenExecutionLedger(path string) (*ExecutionLedger, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("order ledger path is required")
	}
	if info, err := os.Lstat(path); err == nil {
		if !info.Mode().IsRegular() {
			return nil, errors.New("order ledger must be a regular file")
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("inspect order ledger: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("create order ledger directory: %w", err)
	}
	db, err := sql.Open("sqlite", path+"?_journal_mode=WAL&_busy_timeout=5000&_synchronous=FULL")
	if err != nil {
		return nil, fmt.Errorf("open order ledger: %w", err)
	}
	db.SetMaxOpenConns(1)
	ledger := &ExecutionLedger{db: db}
	if err := ledger.init(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := os.Chmod(path, 0o600); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("secure order ledger: %w", err)
	}
	return ledger, nil
}

func (l *ExecutionLedger) init() error {
	const ddl = `
CREATE TABLE IF NOT EXISTS order_executions (
  id TEXT PRIMARY KEY,
  mode TEXT NOT NULL,
  status TEXT NOT NULL,
  order_id TEXT NOT NULL DEFAULT '',
  intent_json BLOB NOT NULL,
  result_json BLOB,
  applied INTEGER NOT NULL DEFAULT 0,
  created_at_ns INTEGER NOT NULL,
  updated_at_ns INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_order_executions_pending
  ON order_executions(mode, status, applied);
CREATE INDEX IF NOT EXISTS idx_order_executions_order_id
  ON order_executions(order_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_order_executions_single_pending
  ON order_executions(mode) WHERE status = 'pending';
`
	if _, err := l.db.Exec(ddl); err != nil {
		return fmt.Errorf("initialize order ledger: %w", err)
	}
	return nil
}

func (l *ExecutionLedger) Close() error {
	if l == nil || l.db == nil {
		return nil
	}
	return l.db.Close()
}

func (l *ExecutionLedger) begin(mode string, intent Intent) (string, error) {
	id, err := executionID()
	if err != nil {
		return "", err
	}
	raw, err := json.Marshal(intent)
	if err != nil {
		return "", err
	}
	now := time.Now().UTC().UnixNano()
	_, err = l.db.Exec(`INSERT INTO order_executions
    (id, mode, status, intent_json, created_at_ns, updated_at_ns)
    VALUES (?, ?, ?, ?, ?, ?)`, id, mode, StatusPending, raw, now, now)
	if err != nil {
		return "", fmt.Errorf("record order intent: %w", err)
	}
	return id, nil
}

func (l *ExecutionLedger) prepared(executionID, orderID string) error {
	result, err := l.db.Exec(`UPDATE order_executions SET order_id = ?, updated_at_ns = ? WHERE id = ?`,
		orderID, time.Now().UTC().UnixNano(), executionID)
	if err != nil {
		return fmt.Errorf("record prepared order: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("record prepared order affected rows: %w", err)
	}
	if rows != 1 {
		return fmt.Errorf("record prepared order: execution %s not found", executionID)
	}
	return nil
}

func (l *ExecutionLedger) complete(executionID string, result Result, submitErr error) error {
	return l.completeWithApply(executionID, result, submitErr, true)
}

func (l *ExecutionLedger) completeUnapplied(executionID string, result Result, submitErr error) error {
	return l.completeWithApply(executionID, result, submitErr, false)
}

func (l *ExecutionLedger) completeWithApply(executionID string, result Result, submitErr error, applyTerminalNonFill bool) error {
	result.ExecutionID = executionID
	status := result.Status
	if status == "" {
		status = StatusRejected
		if submitErr != nil {
			result.Error = submitErr.Error()
		}
	}
	raw, err := json.Marshal(result)
	if err != nil {
		return err
	}
	applied := 0
	if applyTerminalNonFill && status != StatusFilled && status != StatusPending {
		applied = 1
	}
	dbResult, err := l.db.Exec(`UPDATE order_executions
    SET status = ?, order_id = CASE WHEN ? <> '' THEN ? ELSE order_id END,
        result_json = ?, applied = ?, updated_at_ns = ?
    WHERE id = ?`, status, result.OrderID, result.OrderID, raw, applied,
		time.Now().UTC().UnixNano(), executionID)
	if err != nil {
		return fmt.Errorf("complete order execution: %w", err)
	}
	if rows, rowsErr := dbResult.RowsAffected(); rowsErr != nil || rows != 1 {
		return fmt.Errorf("complete order execution %s: affected=%d err=%v", executionID, rows, rowsErr)
	}
	return nil
}

func (l *ExecutionLedger) MarkApplied(executionID string) error {
	if strings.TrimSpace(executionID) == "" {
		return nil
	}
	result, err := l.db.Exec(`UPDATE order_executions SET applied = 1, updated_at_ns = ? WHERE id = ?`,
		time.Now().UTC().UnixNano(), executionID)
	if err != nil {
		return err
	}
	if rows, rowsErr := result.RowsAffected(); rowsErr != nil || rows != 1 {
		return fmt.Errorf("mark execution %s applied: affected=%d err=%v", executionID, rows, rowsErr)
	}
	return nil
}

func (l *ExecutionLedger) ResolveInterruptedPaper() error {
	_, err := l.db.Exec(`UPDATE order_executions
    SET status = ?, applied = 1, updated_at_ns = ?
    WHERE mode = 'paper' AND status = ?`, StatusRejected, time.Now().UTC().UnixNano(), StatusPending)
	return err
}

func (l *ExecutionLedger) UnresolvedCount(mode string) (int, error) {
	var count int
	err := l.db.QueryRow(`SELECT COUNT(*) FROM order_executions
    WHERE mode = ? AND status = ?`, mode, StatusPending).Scan(&count)
	return count, err
}

func (l *ExecutionLedger) OutstandingCount(mode string) (int, error) {
	var count int
	err := l.db.QueryRow(`SELECT COUNT(*) FROM order_executions
    WHERE mode = ? AND (status = ? OR applied = 0)`, mode, StatusPending).Scan(&count)
	return count, err
}

func (l *ExecutionLedger) StaleUnresolvedCount(mode string, before time.Time) (int, error) {
	var count int
	err := l.db.QueryRow(`SELECT COUNT(*) FROM order_executions
    WHERE mode = ? AND status = ? AND updated_at_ns <= ?`,
		mode, StatusPending, before.UTC().UnixNano()).Scan(&count)
	return count, err
}

func (l *ExecutionLedger) StaleUnappliedCount(mode string, before time.Time) (int, error) {
	var count int
	err := l.db.QueryRow(`SELECT COUNT(*) FROM order_executions
    WHERE mode = ? AND status <> ? AND applied = 0 AND updated_at_ns <= ?`,
		mode, StatusPending, before.UTC().UnixNano()).Scan(&count)
	return count, err
}

func (l *ExecutionLedger) UnappliedFills(mode string) ([]ExecutionRecord, error) {
	return l.records(`SELECT id, mode, status, order_id, intent_json, result_json, applied, created_at_ns, updated_at_ns
    FROM order_executions WHERE mode = ? AND status = ? AND applied = 0 ORDER BY created_at_ns`, mode, StatusFilled)
}

func (l *ExecutionLedger) UnappliedNonFills(mode string) ([]ExecutionRecord, error) {
	return l.records(`SELECT id, mode, status, order_id, intent_json, result_json, applied, created_at_ns, updated_at_ns
    FROM order_executions WHERE mode = ? AND status <> ? AND status <> ? AND applied = 0 ORDER BY created_at_ns`,
		mode, StatusFilled, StatusPending)
}

func (l *ExecutionLedger) pending(mode string) ([]ExecutionRecord, error) {
	return l.records(`SELECT id, mode, status, order_id, intent_json, result_json, applied, created_at_ns, updated_at_ns
    FROM order_executions WHERE mode = ? AND status = ? ORDER BY created_at_ns`, mode, StatusPending)
}

func (l *ExecutionLedger) records(query string, args ...any) ([]ExecutionRecord, error) {
	rows, err := l.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []ExecutionRecord
	for rows.Next() {
		var record ExecutionRecord
		var intentRaw, resultRaw []byte
		var applied int
		var createdNS, updatedNS int64
		if err := rows.Scan(&record.ID, &record.Mode, &record.Status, &record.OrderID,
			&intentRaw, &resultRaw, &applied, &createdNS, &updatedNS); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(intentRaw, &record.Intent); err != nil {
			return nil, fmt.Errorf("decode ledger intent %s: %w", record.ID, err)
		}
		if len(resultRaw) > 0 {
			if err := json.Unmarshal(resultRaw, &record.Result); err != nil {
				return nil, fmt.Errorf("decode ledger result %s: %w", record.ID, err)
			}
		}
		record.Result.ExecutionID = record.ID
		record.Applied = applied != 0
		record.CreatedAt = time.Unix(0, createdNS).UTC()
		record.UpdatedAt = time.Unix(0, updatedNS).UTC()
		records = append(records, record)
	}
	return records, rows.Err()
}

func (l *ExecutionLedger) ReconcileLive(ctx context.Context, client *V2Client) error {
	records, err := l.pending("live")
	if err != nil {
		return err
	}
	var errs []error
	for _, record := range records {
		if record.OrderID == "" {
			errs = append(errs, fmt.Errorf("execution %s has no prepared order id", record.ID))
			continue
		}
		status, getErr := client.GetOrder(ctx, record.OrderID)
		if getErr != nil {
			errs = append(errs, fmt.Errorf("execution %s lookup: %w", record.ID, getErr))
			continue
		}
		normalized := strings.ToUpper(status.Status)
		var result Result
		switch {
		case normalized == "MATCHED", normalized == "ORDER_STATUS_MATCHED":
			result = client.resultForFill(ctx, record.OrderID, record.Intent, record.CreatedAt, status)
		case normalized == "CANCELLED", normalized == "CANCELED", normalized == "ORDER_STATUS_CANCELED", normalized == "ORDER_STATUS_CANCELED_MARKET_RESOLVED":
			if status.SizeMatched > 0 {
				result = client.resultForPartialFill(ctx, record.OrderID, record.Intent, record.CreatedAt, status)
			} else {
				result = Result{OrderID: record.OrderID, Status: StatusExpired, SubmitAt: record.CreatedAt, Error: "reconciled canceled order"}
			}
		default:
			result, getErr = client.cancelAndResolve(record.OrderID, record.Intent, record.CreatedAt, "startup reconciliation")
		}
		if err := l.completeUnapplied(record.ID, result, getErr); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

type LedgerClient struct {
	inner  Client
	ledger *ExecutionLedger
	mode   string
	mu     sync.Mutex
}

func NewLedgerClient(inner Client, ledger *ExecutionLedger, mode string) (*LedgerClient, error) {
	if inner == nil || ledger == nil {
		return nil, errors.New("order: ledger client requires inner client and ledger")
	}
	if mode != "paper" && mode != "live" {
		return nil, fmt.Errorf("order: invalid ledger mode %q", mode)
	}
	return &LedgerClient{inner: inner, ledger: ledger, mode: mode}, nil
}

func (c *LedgerClient) Name() string { return c.inner.Name() + "-ledger" }

func (c *LedgerClient) Submit(ctx context.Context, intent Intent) (Result, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	outstanding, err := c.ledger.OutstandingCount(c.mode)
	if err != nil {
		return Result{Status: StatusRejected, Error: err.Error()}, fmt.Errorf("order: check outstanding executions: %w", err)
	}
	if outstanding > 0 {
		err := fmt.Errorf("order: %d outstanding %s execution(s); reconciliation or state persistence required", outstanding, c.mode)
		return Result{Status: StatusRejected, Error: err.Error()}, err
	}
	executionID, err := c.ledger.begin(c.mode, intent)
	if err != nil {
		if count, countErr := c.ledger.OutstandingCount(c.mode); countErr == nil && count > 0 {
			err = fmt.Errorf("order: concurrent %s execution won the durable reservation: %w", c.mode, err)
		}
		return Result{Status: StatusRejected, Error: err.Error()}, err
	}
	ctx = withPreparedOrderObserver(ctx, func(orderID string) error {
		return c.ledger.prepared(executionID, orderID)
	})
	result, submitErr := c.inner.Submit(ctx, intent)
	result.ExecutionID = executionID
	if result.Status == "" {
		result.Status = StatusRejected
		if submitErr != nil {
			result.Error = submitErr.Error()
		} else {
			result.Error = "inner order client returned no status"
			submitErr = errors.New(result.Error)
		}
	}
	if err := c.ledger.complete(executionID, result, submitErr); err != nil {
		message := fmt.Sprintf("execution result not durable: %v", err)
		return Result{
			ExecutionID: executionID,
			OrderID:     result.OrderID,
			Status:      StatusPending,
			SubmitAt:    result.SubmitAt,
			Error:       message,
		}, errors.Join(submitErr, errors.New(message))
	}
	return result, submitErr
}

func (c *LedgerClient) CancelAllOpen(ctx context.Context) error {
	admin, ok := c.inner.(interface{ CancelAllOpen(context.Context) error })
	if !ok {
		return errors.New("order: inner client does not support cancel-all")
	}
	return admin.CancelAllOpen(ctx)
}

func executionID() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(raw[:]), nil
}
