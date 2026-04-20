// Package journal persists closed paper/real trades as one JSONL line per trade.
// Survives bot restarts so the SGT 00:00 daily-report cron can compute realized
// PnL even if the bot crashed mid-day.
package journal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SGT is Asia/Singapore (UTC+8); all journal partitioning + day rolls use it
// because the boss reads the report on SGT clock.
var SGT = time.FixedZone("SGT", 8*3600)

// TradeRecord is a single closed position. Field names are stable JSON contract;
// the daily-report mode reads them by name.
type TradeRecord struct {
	ID           string    `json:"id"`
	AssetID      string    `json:"asset_id"`
	Market       string    `json:"market"`
	Question     string    `json:"question"`
	Outcome      string    `json:"outcome"`
	Side         string    `json:"side"` // always "buy" for now (long-only)
	SizeUSD      float64   `json:"size_usd"`
	Units        float64   `json:"units"`
	EntryMid     float64   `json:"entry_mid"`
	EntryTime    time.Time `json:"entry_time"`
	ExitMid      float64   `json:"exit_mid"`
	ExitTime     time.Time `json:"exit_time"`
	ExitReason   string    `json:"exit_reason"`
	HeldSec      int       `json:"held_sec"`
	PnLUSD       float64   `json:"pnl_usd"`
	OpenOrderID  string    `json:"open_order_id"`
	CloseOrderID string    `json:"close_order_id"`
	Mode         string    `json:"mode"`          // "paper" / "live"
	SignalSource string    `json:"signal_source"` // "auto" / "manual"
}

// Journal appends TradeRecords into ./db/trades-YYYY-MM-DD.jsonl partitioned
// by SGT entry date. Concurrent-safe across goroutines (mu guards the open file
// handle, not the file itself — multiple processes would race, but this bot
// is single-instance).
type Journal struct {
	dir string

	mu     sync.Mutex
	day    string // "YYYY-MM-DD" SGT
	f      *os.File
	writer *bufio.Writer
}

func New(dir string) (*Journal, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("journal mkdir: %w", err)
	}
	return &Journal{dir: dir}, nil
}

// Append writes one trade record. Rotates the underlying file at SGT midnight.
func (j *Journal) Append(rec TradeRecord) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	day := rec.EntryTime.In(SGT).Format("2006-01-02")
	if rec.EntryTime.IsZero() {
		day = time.Now().In(SGT).Format("2006-01-02")
	}
	if j.f == nil || j.day != day {
		if err := j.rotateLocked(day); err != nil {
			return err
		}
	}
	line, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("journal marshal: %w", err)
	}
	if _, err := j.writer.Write(line); err != nil {
		return err
	}
	if err := j.writer.WriteByte('\n'); err != nil {
		return err
	}
	return j.writer.Flush()
}

// Close flushes & closes the current file. Safe to call multiple times.
func (j *Journal) Close() error {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.closeLocked()
}

func (j *Journal) rotateLocked(day string) error {
	if err := j.closeLocked(); err != nil {
		return err
	}
	path := filepath.Join(j.dir, "trades-"+day+".jsonl")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("journal open %s: %w", path, err)
	}
	j.f = f
	j.writer = bufio.NewWriter(f)
	j.day = day
	return nil
}

func (j *Journal) closeLocked() error {
	if j.writer != nil {
		_ = j.writer.Flush()
		j.writer = nil
	}
	if j.f != nil {
		err := j.f.Close()
		j.f = nil
		return err
	}
	return nil
}

// Path returns the file path for a given SGT day. Useful for the report cron.
func Path(dir, day string) string {
	return filepath.Join(dir, "trades-"+day+".jsonl")
}

// Read returns all records for one SGT day. Missing file → nil slice, nil err.
func Read(dir, day string) ([]TradeRecord, error) {
	path := Path(dir, day)
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var out []TradeRecord
	dec := json.NewDecoder(f)
	for {
		var rec TradeRecord
		if err := dec.Decode(&rec); err != nil {
			if err == io.EOF {
				return out, nil
			}
			return out, fmt.Errorf("journal decode %s: %w", path, err)
		}
		out = append(out, rec)
	}
}
