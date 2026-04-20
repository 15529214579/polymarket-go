package tickrec

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

func newRec(t *testing.T) (*Recorder, string) {
	t.Helper()
	dir := t.TempDir()
	r, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return r, dir
}

func mkTick(mid float64, sec int64) feed.Tick {
	return feed.Tick{
		AssetID: "A",
		Market:  "M",
		Time:    time.Unix(sec, 0).UTC(),
		BestBid: mid - 0.005,
		BestAsk: mid + 0.005,
		Mid:     mid,
		Trades:  1,
		BuyVol:  10,
		SellVol: 5,
	}
}

func readRows(t *testing.T, path string) []TickRow {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open jsonl: %v", err)
	}
	defer f.Close()
	var out []TickRow
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var r TickRow
		if err := json.Unmarshal(sc.Bytes(), &r); err != nil {
			t.Fatalf("decode: %v", err)
		}
		out = append(out, r)
	}
	return out
}

func TestRecorder_RecordsAfterStart(t *testing.T) {
	r, _ := newRec(t)
	if err := r.Start("p1", "asset-x"); err != nil {
		t.Fatalf("Start: %v", err)
	}
	for i, mid := range []float64{0.50, 0.51, 0.52} {
		if err := r.Record("p1", mkTick(mid, int64(1000+i))); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}
	if err := r.Stop("p1"); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	rows := readRows(t, r.Path("p1"))
	if len(rows) != 3 {
		t.Fatalf("want 3 rows, got %d", len(rows))
	}
	if rows[0].AssetID != "asset-x" || rows[2].Mid != 0.52 {
		t.Fatalf("payload mismatch: %+v", rows)
	}
}

func TestRecorder_DedupesWithinSameSecond(t *testing.T) {
	r, _ := newRec(t)
	_ = r.Start("p1", "asset-x")
	// Three samples all at unix second 5000 — only the first should land.
	_ = r.Record("p1", mkTick(0.40, 5000))
	_ = r.Record("p1", mkTick(0.41, 5000))
	_ = r.Record("p1", mkTick(0.42, 5000))
	// New second now writes.
	_ = r.Record("p1", mkTick(0.50, 5001))
	_ = r.Stop("p1")
	rows := readRows(t, r.Path("p1"))
	if len(rows) != 2 {
		t.Fatalf("want 2 rows (dedup), got %d", len(rows))
	}
	if rows[0].Mid != 0.40 || rows[1].Mid != 0.50 {
		t.Fatalf("dedup picked wrong rows: %+v", rows)
	}
}

func TestRecorder_UnknownPosIsNoop(t *testing.T) {
	r, _ := newRec(t)
	if err := r.Record("ghost", mkTick(0.5, 1)); err != nil {
		t.Fatalf("Record on unknown should be nil, got %v", err)
	}
	if err := r.Stop("ghost"); err != nil {
		t.Fatalf("Stop on unknown should be nil, got %v", err)
	}
}

func TestRecorder_StartIdempotent(t *testing.T) {
	r, _ := newRec(t)
	_ = r.Start("p1", "asset-x")
	_ = r.Start("p1", "asset-y") // second Start should not clobber the recording
	_ = r.Record("p1", mkTick(0.6, 9000))
	_ = r.Stop("p1")
	rows := readRows(t, r.Path("p1"))
	if len(rows) != 1 || rows[0].AssetID != "asset-x" {
		t.Fatalf("idempotency violated: %+v", rows)
	}
}

func TestRecorder_AppendsAcrossStartStopCycles(t *testing.T) {
	r, _ := newRec(t)
	_ = r.Start("p1", "asset-x")
	_ = r.Record("p1", mkTick(0.30, 1))
	_ = r.Stop("p1")

	// Re-Start same posID; new lastSec resets, but file is opened with O_APPEND
	// so prior rows survive. (lastSec only dedupes within a single Start cycle.)
	_ = r.Start("p1", "asset-x")
	_ = r.Record("p1", mkTick(0.31, 2))
	_ = r.Stop("p1")
	rows := readRows(t, r.Path("p1"))
	if len(rows) != 2 {
		t.Fatalf("want 2 rows across cycles, got %d", len(rows))
	}
}

func TestRecorder_Snapshot(t *testing.T) {
	r, _ := newRec(t)
	_ = r.Start("p1", "A1")
	_ = r.Start("p2", "A2")
	snap := r.Snapshot()
	if len(snap) != 2 || snap["p1"] != "A1" || snap["p2"] != "A2" {
		t.Fatalf("snapshot wrong: %+v", snap)
	}
	_ = r.Stop("p1")
	snap = r.Snapshot()
	if _, ok := snap["p1"]; ok || snap["p2"] != "A2" {
		t.Fatalf("snapshot after stop wrong: %+v", snap)
	}
}

func TestRecorder_PathInsideDir(t *testing.T) {
	r, dir := newRec(t)
	if got := r.Path("p1"); filepath.Dir(got) != dir {
		t.Fatalf("Path outside dir: %s vs %s", got, dir)
	}
}
