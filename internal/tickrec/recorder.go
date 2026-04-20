// Package tickrec persists per-position 1Hz tick paths to JSONL files so
// Phase 7.e can replay our own momentum opens against true price paths
// (instead of the python end-point approximation used in Phase 7.d).
//
// One file per position; rows are append-only; deduped to one row per second
// to match the sampler's 1-Hz cadence. Recorder is safe for concurrent use.
package tickrec

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

// TickRow is one persisted sample. Field names are short for replay-DB size.
type TickRow struct {
	PosID   string    `json:"pos_id"`
	AssetID string    `json:"asset_id"`
	Time    time.Time `json:"t"`
	BestBid float64   `json:"bid,omitempty"`
	BestAsk float64   `json:"ask,omitempty"`
	Mid     float64   `json:"mid"`
	Trades  int       `json:"trades,omitempty"`
	BuyVol  float64   `json:"buy_vol,omitempty"`
	SellVol float64   `json:"sell_vol,omitempty"`
}

// Recorder owns a directory of pos-keyed JSONL files.
type Recorder struct {
	dir string

	mu   sync.Mutex
	open map[string]*recording // posID → recording
}

type recording struct {
	posID   string
	assetID string
	f       *os.File
	enc     *json.Encoder
	lastSec int64 // last unix second written; dedupe within-second resamples
}

// New creates the directory if missing and returns a ready recorder.
func New(dir string) (*Recorder, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("tickrec mkdir: %w", err)
	}
	return &Recorder{dir: dir, open: map[string]*recording{}}, nil
}

// Start begins recording for posID. Idempotent: a second Start with the same
// posID is a no-op (the existing file/encoder is reused).
func (r *Recorder) Start(posID, assetID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.open[posID]; ok {
		return nil
	}
	path := filepath.Join(r.dir, posID+".jsonl")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("tickrec open %s: %w", path, err)
	}
	r.open[posID] = &recording{
		posID:   posID,
		assetID: assetID,
		f:       f,
		enc:     json.NewEncoder(f),
	}
	return nil
}

// Record appends one tick if posID is open and the tick second is newer than
// the last persisted row. Unknown posIDs are silently ignored.
func (r *Recorder) Record(posID string, tick feed.Tick) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, ok := r.open[posID]
	if !ok {
		return nil
	}
	sec := tick.Time.Unix()
	if sec <= rec.lastSec {
		return nil
	}
	row := TickRow{
		PosID:   posID,
		AssetID: rec.assetID,
		Time:    tick.Time,
		BestBid: tick.BestBid,
		BestAsk: tick.BestAsk,
		Mid:     tick.Mid,
		Trades:  tick.Trades,
		BuyVol:  tick.BuyVol,
		SellVol: tick.SellVol,
	}
	if err := rec.enc.Encode(&row); err != nil {
		return fmt.Errorf("tickrec encode %s: %w", posID, err)
	}
	rec.lastSec = sec
	return nil
}

// Stop closes the recording for posID. Idempotent: stopping an unknown posID
// is a no-op. After Stop the file remains on disk for replay.
func (r *Recorder) Stop(posID string) error {
	r.mu.Lock()
	rec, ok := r.open[posID]
	if ok {
		delete(r.open, posID)
	}
	r.mu.Unlock()
	if !ok {
		return nil
	}
	if err := rec.f.Close(); err != nil {
		return fmt.Errorf("tickrec close %s: %w", posID, err)
	}
	return nil
}

// Snapshot returns posID→assetID for all open recordings; the recorder ticker
// uses it to know which assets to poll without holding the lock.
func (r *Recorder) Snapshot() map[string]string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make(map[string]string, len(r.open))
	for k, v := range r.open {
		out[k] = v.assetID
	}
	return out
}

// Path returns the on-disk JSONL path for posID (whether or not currently
// open). Useful for tests and offline tooling.
func (r *Recorder) Path(posID string) string {
	return filepath.Join(r.dir, posID+".jsonl")
}
