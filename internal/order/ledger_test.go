package order

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLedgerClientPersistsFillUntilApplied(t *testing.T) {
	path := filepath.Join(t.TempDir(), "orders.sqlite")
	ledger, err := OpenExecutionLedger(path)
	if err != nil {
		t.Fatal(err)
	}
	defer ledger.Close()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("ledger mode=%#o", info.Mode().Perm())
	}
	client, err := NewLedgerClient(NewPaperClient(0), ledger, "paper")
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.Submit(context.Background(), Intent{
		ClientID: "p1", AssetID: "asset", Side: Buy, SizeUSD: 20, LimitPx: 0.5, Type: FAK,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.ExecutionID == "" || result.Status != StatusFilled {
		t.Fatalf("result=%+v", result)
	}
	records, err := ledger.UnappliedFills("paper")
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].Intent.ClientID != "p1" {
		t.Fatalf("records=%+v", records)
	}
	if err := ledger.MarkApplied(result.ExecutionID); err != nil {
		t.Fatal(err)
	}
	records, err = ledger.UnappliedFills("paper")
	if err != nil || len(records) != 0 {
		t.Fatalf("records=%+v err=%v", records, err)
	}
}

func TestLedgerMarkAppliedRejectsUnknownExecution(t *testing.T) {
	ledger, err := OpenExecutionLedger(filepath.Join(t.TempDir(), "orders.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer ledger.Close()
	if err := ledger.MarkApplied("missing"); err == nil {
		t.Fatal("expected missing execution error")
	}
}

func TestOpenExecutionLedgerRejectsSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target.sqlite")
	if err := os.WriteFile(target, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "orders.sqlite")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if ledger, err := OpenExecutionLedger(link); err == nil {
		_ = ledger.Close()
		t.Fatal("expected ledger symlink rejection")
	}
}

type preparedTestClient struct{}

func (preparedTestClient) Name() string { return "prepared-test" }

func (preparedTestClient) Submit(ctx context.Context, _ Intent) (Result, error) {
	if err := notifyPreparedOrder(ctx, "0xprepared"); err != nil {
		return Result{}, err
	}
	return Result{OrderID: "0xprepared", Status: StatusPending}, errors.New("network outcome unknown")
}

type emptyResultClient struct{}

func (emptyResultClient) Name() string { return "empty-result-test" }

func (emptyResultClient) Submit(context.Context, Intent) (Result, error) {
	return Result{}, nil
}

func TestLedgerClientRejectsEmptyInnerResult(t *testing.T) {
	ledger, err := OpenExecutionLedger(filepath.Join(t.TempDir(), "orders.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer ledger.Close()
	client, err := NewLedgerClient(emptyResultClient{}, ledger, "paper")
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.Submit(context.Background(), Intent{
		AssetID: "asset", Side: Buy, SizeUSD: 20, LimitPx: 0.5, Type: FAK,
	})
	if err == nil || result.Status != StatusRejected {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}

func TestLedgerClientStoresPreparedOrderBeforeUnknownResult(t *testing.T) {
	ledger, err := OpenExecutionLedger(filepath.Join(t.TempDir(), "orders.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer ledger.Close()
	client, err := NewLedgerClient(preparedTestClient{}, ledger, "live")
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.Submit(context.Background(), Intent{
		ClientID: "p1", AssetID: "123", Side: Buy, SizeUSD: 20, LimitPx: 0.5, Type: FAK,
	})
	if err == nil || result.Status != StatusPending {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	records, err := ledger.pending("live")
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].OrderID != "0xprepared" {
		t.Fatalf("records=%+v", records)
	}
	second, err := client.Submit(context.Background(), Intent{
		ClientID: "p2", AssetID: "456", Side: Buy, SizeUSD: 20, LimitPx: 0.5, Type: FAK,
	})
	if err == nil || second.Status != StatusRejected || !strings.Contains(err.Error(), "outstanding") {
		t.Fatalf("second=%+v err=%v", second, err)
	}
	records, err = ledger.pending("live")
	if err != nil || len(records) != 1 {
		t.Fatalf("records=%+v err=%v", records, err)
	}
}

func TestLedgerClientBlocksNextOrderUntilFillIsApplied(t *testing.T) {
	path := filepath.Join(t.TempDir(), "orders.sqlite")
	firstLedger, err := OpenExecutionLedger(path)
	if err != nil {
		t.Fatal(err)
	}
	defer firstLedger.Close()
	secondLedger, err := OpenExecutionLedger(path)
	if err != nil {
		t.Fatal(err)
	}
	defer secondLedger.Close()

	first, err := NewLedgerClient(NewPaperClient(0), firstLedger, "live")
	if err != nil {
		t.Fatal(err)
	}
	fill, err := first.Submit(context.Background(), Intent{
		ClientID: "p1", AssetID: "asset", Side: Buy, SizeUSD: 5, LimitPx: 0.5, Type: FAK,
	})
	if err != nil || fill.Status != StatusFilled {
		t.Fatalf("fill=%+v err=%v", fill, err)
	}

	second, err := NewLedgerClient(NewPaperClient(0), secondLedger, "live")
	if err != nil {
		t.Fatal(err)
	}
	rejected, err := second.Submit(context.Background(), Intent{
		ClientID: "p2", AssetID: "asset-2", Side: Buy, SizeUSD: 5, LimitPx: 0.5, Type: FAK,
	})
	if err == nil || rejected.Status != StatusRejected || !strings.Contains(err.Error(), "outstanding") {
		t.Fatalf("rejected=%+v err=%v", rejected, err)
	}
	if err := firstLedger.MarkApplied(fill.ExecutionID); err != nil {
		t.Fatal(err)
	}
	accepted, err := second.Submit(context.Background(), Intent{
		ClientID: "p2", AssetID: "asset-2", Side: Buy, SizeUSD: 5, LimitPx: 0.5, Type: FAK,
	})
	if err != nil || accepted.Status != StatusFilled {
		t.Fatalf("accepted=%+v err=%v", accepted, err)
	}
}

func TestExecutionLedgerStaleHealthCounts(t *testing.T) {
	ledger, err := OpenExecutionLedger(filepath.Join(t.TempDir(), "orders.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer ledger.Close()

	executionID, err := ledger.begin("live", Intent{ClientID: "p1", AssetID: "asset", Side: Sell})
	if err != nil {
		t.Fatal(err)
	}
	if count, err := ledger.StaleUnresolvedCount("live", time.Now().Add(-time.Minute)); err != nil || count != 0 {
		t.Fatalf("fresh pending count=%d err=%v", count, err)
	}
	if count, err := ledger.StaleUnresolvedCount("live", time.Now().Add(time.Second)); err != nil || count != 1 {
		t.Fatalf("stale pending count=%d err=%v", count, err)
	}

	if err := ledger.completeUnapplied(executionID, Result{Status: StatusExpired}, nil); err != nil {
		t.Fatal(err)
	}
	if count, err := ledger.StaleUnresolvedCount("live", time.Now().Add(time.Second)); err != nil || count != 0 {
		t.Fatalf("terminal pending count=%d err=%v", count, err)
	}
	if count, err := ledger.StaleUnappliedCount("live", time.Now().Add(time.Second)); err != nil || count != 1 {
		t.Fatalf("terminal unapplied count=%d err=%v", count, err)
	}
	records, err := ledger.UnappliedNonFills("live")
	if err != nil || len(records) != 1 || records[0].ID != executionID {
		t.Fatalf("non-fills=%+v err=%v", records, err)
	}
	if err := ledger.MarkApplied(executionID); err != nil {
		t.Fatal(err)
	}
	if count, err := ledger.StaleUnappliedCount("live", time.Now().Add(time.Second)); err != nil || count != 0 {
		t.Fatalf("applied count=%d err=%v", count, err)
	}
}

func TestExecutionLedgerAllowsOnlyOnePendingOrderAcrossConnections(t *testing.T) {
	path := filepath.Join(t.TempDir(), "orders.sqlite")
	first, err := OpenExecutionLedger(path)
	if err != nil {
		t.Fatal(err)
	}
	defer first.Close()
	second, err := OpenExecutionLedger(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()

	if _, err := first.begin("live", Intent{ClientID: "p1", Side: Buy}); err != nil {
		t.Fatal(err)
	}
	if _, err := second.begin("live", Intent{ClientID: "p2", Side: Buy}); err == nil {
		t.Fatal("expected unique pending execution rejection")
	}
}
