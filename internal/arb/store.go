package arb

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/15529214579/polymarket-go/internal/odds"
	_ "modernc.org/sqlite"
)

// Store persists odds snapshots to SQLite for iteration support.
type Store struct {
	db *sql.DB
}

// NewStore opens (or creates) the SQLite database at path.
func NewStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite %s: %w", path, err)
	}
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS odds_snapshot (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			market_id TEXT NOT NULL,
			sport TEXT,
			event_name TEXT,
			polymarket_price REAL,
			bookmaker TEXT,
			bookmaker_odds REAL,
			bookmaker_prob REAL,
			gap_pp REAL,
			snapshot_timestamp TEXT NOT NULL
		)
	`); err != nil {
		db.Close()
		return nil, fmt.Errorf("create table: %w", err)
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_snapshot_ts ON odds_snapshot(snapshot_timestamp)`); err != nil {
		db.Close()
		return nil, fmt.Errorf("create index: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

// Insert writes one snapshot row.
func (s *Store) Insert(snap odds.OddsSnapshot) {
	ts := snap.SnapshotTimestamp.UTC().Format(time.RFC3339)
	_, err := s.db.Exec(`
		INSERT INTO odds_snapshot (market_id, sport, event_name, polymarket_price,
			bookmaker, bookmaker_odds, bookmaker_prob, gap_pp, snapshot_timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		snap.MarketID, snap.Sport, snap.EventName, snap.PolymarketPrice,
		snap.Bookmaker, snap.BookmakerOddsVal, snap.BookmakerProb,
		snap.GapPP, ts,
	)
	if err != nil {
		slog.Warn("odds_snapshot_insert_fail", "err", err.Error())
	}
}

// Recent returns the last N snapshots ordered by timestamp desc.
func (s *Store) Recent(limit int) ([]odds.OddsSnapshot, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := s.db.Query(`
		SELECT id, market_id, sport, event_name, polymarket_price,
			bookmaker, bookmaker_odds, bookmaker_prob, gap_pp, snapshot_timestamp
		FROM odds_snapshot
		ORDER BY snapshot_timestamp DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []odds.OddsSnapshot
	for rows.Next() {
		var snap odds.OddsSnapshot
		var ts string
		err := rows.Scan(&snap.ID, &snap.MarketID, &snap.Sport, &snap.EventName,
			&snap.PolymarketPrice, &snap.Bookmaker, &snap.BookmakerOddsVal,
			&snap.BookmakerProb, &snap.GapPP, &ts)
		if err != nil {
			continue
		}
		snap.SnapshotTimestamp, _ = time.Parse(time.RFC3339, ts)
		result = append(result, snap)
	}
	return result, nil
}

// Count returns total rows in odds_snapshot.
func (s *Store) Count() int {
	var n int
	_ = s.db.QueryRow("SELECT COUNT(*) FROM odds_snapshot").Scan(&n)
	return n
}
