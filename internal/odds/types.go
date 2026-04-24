package odds

import "time"

// BookmakerOdds holds one outcome with juice-removed probability from a bookmaker.
type BookmakerOdds struct {
	Sport             string  `json:"sport"`
	EventID           string  `json:"event_id"`
	EventName         string  `json:"event_name"`
	TeamOrSide        string  `json:"team_or_side"`
	BookmakerProb     float64 `json:"bookmaker_prob"`
	Bookmaker         string  `json:"bookmaker"`
	MarketName        string  `json:"market_name"`
	EventCommenceTime string  `json:"event_commence_time,omitempty"`
}

// ArbOpportunity represents a price gap between a bookmaker and Polymarket.
type ArbOpportunity struct {
	TokenID         string  `json:"token_id"`
	PolymarketPrice float64 `json:"polymarket_price"`
	BookmakerProb   float64 `json:"bookmaker_prob"`
	GapPP           float64 `json:"gap_pp"`
	NetEvPP         float64 `json:"net_ev_pp"`
	Direction       string  `json:"direction"` // "BUY_YES" or "BUY_NO"
	MarketTitle     string  `json:"market_title"`
	Sport           string  `json:"sport"`
	EventName       string  `json:"event_name"`
	Bookmaker       string  `json:"bookmaker"`
	BetSizeUSDC     float64 `json:"bet_size_usdc"`
}

// OddsSnapshot is a row in the odds_snapshot DB table.
type OddsSnapshot struct {
	ID                int64     `json:"id"`
	MarketID          string    `json:"market_id"`
	Sport             string    `json:"sport"`
	EventName         string    `json:"event_name"`
	PolymarketPrice   float64   `json:"polymarket_price"`
	Bookmaker         string    `json:"bookmaker"`
	BookmakerOddsVal  float64   `json:"bookmaker_odds"`
	BookmakerProb     float64   `json:"bookmaker_prob"`
	GapPP             float64   `json:"gap_pp"`
	SnapshotTimestamp time.Time `json:"snapshot_timestamp"`
}
