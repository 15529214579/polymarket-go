package btc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

const gammaMarketsURL = "https://gamma-api.polymarket.com/markets"

var etLocation = mustLoadLocation("America/New_York")

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic("load timezone " + name + ": " + err.Error())
	}
	return loc
}

type UpDownConfig struct {
	Enabled       bool
	ScanInterval  time.Duration
	MinConfidence float64
	SizeUSD       float64
	MaxDailyBets  int
	DBPath        string
}

func DefaultUpDownConfig() UpDownConfig {
	return UpDownConfig{
		ScanInterval:  15 * time.Minute,
		MinConfidence: 0.52,
		SizeUSD:       5.0,
		MaxDailyBets:  20,
		DBPath:        "db/btc.db",
	}
}

type UpDownMarket struct {
	Slug         string
	Question     string
	ConditionID  string
	UpTokenID    string
	DownTokenID  string
	UpPrice      float64
	DownPrice    float64
	EndDate      time.Time
	Active       bool
	AcceptOrders bool
}

type UpDownSignal struct {
	MarketSlug         string
	Question           string
	ConditionID        string
	TokenID            string
	PMPrice            float64
	PredictedDirection string
	Confidence         float64
	Spot               float64
	SizeUSD            float64
}

type UpDownSignalCallback func(sig UpDownSignal)

// --- Market Discovery ---

type marketCache struct {
	mu      sync.Mutex
	markets []UpDownMarket
	at      time.Time
	ttl     time.Duration
}

var udCache = &marketCache{ttl: 5 * time.Minute}

func DiscoverUpDownMarkets(ctx context.Context) ([]UpDownMarket, error) {
	udCache.mu.Lock()
	if time.Since(udCache.at) < udCache.ttl && len(udCache.markets) > 0 {
		cached := udCache.markets
		udCache.mu.Unlock()
		return cached, nil
	}
	udCache.mu.Unlock()

	now := time.Now().In(etLocation)
	slugs := generateSlugs(now, 48*time.Hour)

	var markets []UpDownMarket
	for _, slug := range slugs {
		m, err := fetchMarketBySlug(ctx, slug)
		if err != nil {
			continue
		}
		if m.EndDate.Before(time.Now()) {
			continue
		}
		if !m.Active || !m.AcceptOrders {
			continue
		}
		markets = append(markets, m)

		select {
		case <-ctx.Done():
			return markets, ctx.Err()
		case <-time.After(50 * time.Millisecond):
		}
	}

	udCache.mu.Lock()
	udCache.markets = markets
	udCache.at = time.Now()
	udCache.mu.Unlock()

	slog.Info("updown.discover", "found", len(markets), "slugs_tried", len(slugs))
	return markets, nil
}

func generateSlugs(now time.Time, window time.Duration) []string {
	months := []string{
		"january", "february", "march", "april", "may", "june",
		"july", "august", "september", "october", "november", "december",
	}

	var slugs []string
	end := now.Add(window)
	cur := now.Truncate(time.Hour)
	if cur.Before(now) {
		cur = cur.Add(time.Hour)
	}

	seen := make(map[string]bool)
	for cur.Before(end) {
		et := cur.In(etLocation)
		hour := et.Hour()
		if hour < 8 || hour > 23 {
			cur = cur.Add(time.Hour)
			continue
		}

		hourStr := formatHourET(hour)
		slug := fmt.Sprintf("bitcoin-up-or-down-%s-%d-%d-%s-et",
			months[et.Month()-1], et.Day(), et.Year(), hourStr)

		if !seen[slug] {
			slugs = append(slugs, slug)
			seen[slug] = true
		}
		cur = cur.Add(time.Hour)
	}
	return slugs
}

func formatHourET(hour int) string {
	if hour == 0 {
		return "12am"
	}
	if hour == 12 {
		return "12pm"
	}
	if hour < 12 {
		return fmt.Sprintf("%dam", hour)
	}
	return fmt.Sprintf("%dpm", hour-12)
}

type gammaMarketResp struct {
	ID              string `json:"id"`
	ConditionID     string `json:"conditionId"`
	Slug            string `json:"slug"`
	Question        string `json:"question"`
	Active          bool   `json:"active"`
	Closed          bool   `json:"closed"`
	AcceptingOrders bool   `json:"acceptingOrders"`
	EndDate         string `json:"endDate"`
	ClobTokenIDs    string `json:"clobTokenIds"`
	OutcomePrices   string `json:"outcomePrices"`
	Outcomes        string `json:"outcomes"`
}

func fetchMarketBySlug(ctx context.Context, slug string) (UpDownMarket, error) {
	url := fmt.Sprintf("%s?slug=%s", gammaMarketsURL, slug)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return UpDownMarket{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return UpDownMarket{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return UpDownMarket{}, fmt.Errorf("gamma HTTP %d: %s", resp.StatusCode, body)
	}

	var results []gammaMarketResp
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return UpDownMarket{}, err
	}
	if len(results) == 0 {
		return UpDownMarket{}, fmt.Errorf("no market for slug %s", slug)
	}

	r := results[0]
	endDate, _ := time.Parse(time.RFC3339, r.EndDate)

	outcomes := parseJSONStringArray(r.Outcomes)
	tokenIDs := parseJSONStringArray(r.ClobTokenIDs)
	prices := parseJSONStringArray(r.OutcomePrices)

	upIdx, downIdx := -1, -1
	for i, o := range outcomes {
		switch strings.ToLower(o) {
		case "up":
			upIdx = i
		case "down":
			downIdx = i
		}
	}
	if upIdx < 0 || downIdx < 0 {
		return UpDownMarket{}, fmt.Errorf("slug %s: missing Up/Down outcomes", slug)
	}

	upToken, downToken := safeIndex(tokenIDs, upIdx), safeIndex(tokenIDs, downIdx)
	upPrice, _ := strconv.ParseFloat(safeIndex(prices, upIdx), 64)
	downPrice, _ := strconv.ParseFloat(safeIndex(prices, downIdx), 64)

	return UpDownMarket{
		Slug:         slug,
		Question:     r.Question,
		ConditionID:  r.ConditionID,
		UpTokenID:    upToken,
		DownTokenID:  downToken,
		UpPrice:      upPrice,
		DownPrice:    downPrice,
		EndDate:      endDate,
		Active:       r.Active && !r.Closed,
		AcceptOrders: r.AcceptingOrders,
	}, nil
}

func parseJSONStringArray(s string) []string {
	if s == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return nil
	}
	return out
}

func safeIndex(arr []string, i int) string {
	if i < 0 || i >= len(arr) {
		return ""
	}
	return arr[i]
}

// --- Direction Prediction ---

type DirectionPrediction struct {
	Direction  string  // "Up" or "Down"
	Confidence float64 // 0-1
	Alignment  string
}

func PredictHourlyDirection(ctx context.Context) (DirectionPrediction, error) {
	mtf, err := PredictMultiTF(ctx)
	if err != nil {
		return DirectionPrediction{}, fmt.Errorf("multi_tf: %w", err)
	}

	var pred DirectionPrediction
	pred.Alignment = mtf.Alignment

	switch mtf.Alignment {
	case "ALIGNED_BULL":
		pred.Direction = "Up"
		pred.Confidence = mtf.CombinedBull
	case "ALIGNED_BEAR":
		pred.Direction = "Down"
		pred.Confidence = mtf.CombinedBear
	default:
		if mtf.CombinedBull >= mtf.CombinedBear {
			pred.Direction = "Up"
			pred.Confidence = mtf.Confidence
		} else {
			pred.Direction = "Down"
			pred.Confidence = mtf.Confidence
		}
	}

	slog.Info("updown.predict",
		"direction", pred.Direction,
		"confidence", fmt.Sprintf("%.3f", pred.Confidence),
		"alignment", pred.Alignment,
	)
	return pred, nil
}

// --- SQLite ---

func initUpDownDB(db *sql.DB) error {
	const ddl = `CREATE TABLE IF NOT EXISTS updown_bets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp INTEGER NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    question TEXT,
    condition_id TEXT,
    token_id TEXT,
    predicted_direction TEXT,
    confidence REAL,
    pm_up_price REAL,
    pm_down_price REAL,
    btc_spot REAL,
    size_usd REAL,
    actual_direction TEXT,
    pnl REAL,
    resolved_at INTEGER
);
CREATE TABLE IF NOT EXISTS updown_prices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp INTEGER NOT NULL,
    slug TEXT NOT NULL,
    up_price REAL,
    down_price REAL,
    spread REAL,
    deviation REAL,
    UNIQUE(timestamp, slug)
);`
	_, err := db.Exec(ddl)
	return err
}

func logUpDownPrices(ctx context.Context, db *sql.DB, markets []UpDownMarket) {
	now := time.Now().Unix()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer tx.Rollback() //nolint:errcheck

	stmt, err := tx.PrepareContext(ctx,
		`INSERT OR IGNORE INTO updown_prices(timestamp, slug, up_price, down_price, spread, deviation)
		 VALUES(?,?,?,?,?,?)`)
	if err != nil {
		return
	}
	defer stmt.Close()

	for _, m := range markets {
		if m.UpPrice <= 0 || m.DownPrice <= 0 {
			continue
		}
		spread := m.UpPrice + m.DownPrice - 1.0
		cheaperSide := m.UpPrice
		if m.DownPrice < cheaperSide {
			cheaperSide = m.DownPrice
		}
		deviation := 0.50 - cheaperSide
		stmt.ExecContext(ctx, now, m.Slug, m.UpPrice, m.DownPrice, spread, deviation) //nolint:errcheck
	}
	tx.Commit() //nolint:errcheck
}

func slugAlreadyBet(ctx context.Context, db *sql.DB, slug string) bool {
	var n int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM updown_bets WHERE slug = ?", slug).Scan(&n)
	return err == nil && n > 0
}

func countBetsToday(ctx context.Context, db *sql.DB) int {
	startOfDay := time.Now().Truncate(24 * time.Hour).Unix()
	var n int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM updown_bets WHERE timestamp >= ?", startOfDay).Scan(&n) //nolint:errcheck
	return n
}

func recordBet(ctx context.Context, db *sql.DB, m UpDownMarket, pred DirectionPrediction, spot, sizeUSD float64, tokenID string) {
	_, err := db.ExecContext(ctx, `INSERT OR IGNORE INTO updown_bets
		(timestamp, slug, question, condition_id, token_id, predicted_direction, confidence, pm_up_price, pm_down_price, btc_spot, size_usd)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		time.Now().Unix(), m.Slug, m.Question, m.ConditionID, tokenID,
		pred.Direction, pred.Confidence, m.UpPrice, m.DownPrice, spot, sizeUSD,
	)
	if err != nil {
		slog.Warn("updown.record_bet_fail", "slug", m.Slug, "err", err.Error())
	}
}

func logPnLSummary(ctx context.Context, db *sql.DB) {
	var totalBets, resolved, wins int
	var totalPnL float64
	var pending int

	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM updown_bets").Scan(&totalBets)
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM updown_bets WHERE actual_direction IS NOT NULL").Scan(&resolved)
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM updown_bets WHERE actual_direction IS NULL").Scan(&pending)
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM updown_bets WHERE pnl > 0").Scan(&wins)
	db.QueryRowContext(ctx, "SELECT COALESCE(SUM(pnl), 0) FROM updown_bets WHERE pnl IS NOT NULL").Scan(&totalPnL)

	wr := 0.0
	if resolved > 0 {
		wr = float64(wins) / float64(resolved) * 100
	}

	var priceSnapshots int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM updown_prices").Scan(&priceSnapshots)

	slog.Info("updown.pnl_summary",
		"total_bets", totalBets,
		"resolved", resolved,
		"pending", pending,
		"wins", wins,
		"win_rate", fmt.Sprintf("%.1f%%", wr),
		"total_pnl", fmt.Sprintf("$%+.2f", totalPnL),
		"price_snapshots", priceSnapshots,
	)
}

// --- Resolution Checker ---

func resolveSettledBets(ctx context.Context, db *sql.DB) {
	rows, err := db.QueryContext(ctx,
		"SELECT id, slug, predicted_direction, pm_up_price, pm_down_price, size_usd, timestamp FROM updown_bets WHERE actual_direction IS NULL")
	if err != nil {
		slog.Warn("updown.resolve_query_fail", "err", err.Error())
		return
	}
	defer rows.Close()

	type pendingBet struct {
		id        int64
		slug      string
		predicted string
		upPrice   float64
		downPrice float64
		sizeUSD   float64
		ts        int64
	}

	var pending []pendingBet
	for rows.Next() {
		var b pendingBet
		if err := rows.Scan(&b.id, &b.slug, &b.predicted, &b.upPrice, &b.downPrice, &b.sizeUSD, &b.ts); err != nil {
			continue
		}
		pending = append(pending, b)
	}
	if err := rows.Err(); err != nil {
		slog.Warn("updown.resolve_rows_err", "err", err.Error())
	}

	for _, b := range pending {
		candleEnd := parseCandleEndFromSlug(b.slug)
		if candleEnd.IsZero() || time.Now().Before(candleEnd.Add(2*time.Minute)) {
			continue
		}

		candleStart := candleEnd.Add(-1 * time.Hour)
		candles, err := FetchCandlesRange(ctx, "BTCUSDT", Interval1h, candleStart, candleEnd)
		if err != nil || len(candles) == 0 {
			continue
		}

		c := candles[len(candles)-1]
		actual := "Down"
		if c.Close >= c.Open {
			actual = "Up"
		}

		pmPrice := b.downPrice
		if b.predicted == "Up" {
			pmPrice = b.upPrice
		}

		var pnl float64
		if b.predicted == actual {
			if pmPrice > 0 {
				pnl = (1.0/pmPrice - 1.0) * b.sizeUSD
			}
		} else {
			pnl = -b.sizeUSD
		}

		_, err = db.ExecContext(ctx,
			"UPDATE updown_bets SET actual_direction = ?, pnl = ?, resolved_at = ? WHERE id = ?",
			actual, pnl, time.Now().Unix(), b.id)
		if err != nil {
			slog.Warn("updown.resolve_update_fail", "id", b.id, "err", err.Error())
			continue
		}

		candleRange := (c.High - c.Low) / c.Open * 100
		candleBody := (c.Close - c.Open) / c.Open * 100
		ev := ExpectedValue(pmPrice, 0.50)
		kellyF := KellyFraction(pmPrice, 0.50)

		slog.Info("updown.resolved",
			"slug", b.slug,
			"predicted", b.predicted,
			"actual", actual,
			"pnl", fmt.Sprintf("%.2f", pnl),
			"size_usd", fmt.Sprintf("%.2f", b.sizeUSD),
			"pm_price", fmt.Sprintf("%.3f", pmPrice),
			"ev_at_entry", fmt.Sprintf("%.3f", ev),
			"kelly_at_entry", fmt.Sprintf("%.3f", kellyF),
			"candle_open", fmt.Sprintf("%.2f", c.Open),
			"candle_close", fmt.Sprintf("%.2f", c.Close),
			"candle_range_pct", fmt.Sprintf("%.2f%%", candleRange),
			"candle_body_pct", fmt.Sprintf("%+.2f%%", candleBody),
		)
	}
}

// parseCandleEndFromSlug extracts the candle end time from an up/down slug.
// Slug: bitcoin-up-or-down-april-26-2026-10am-et → 2026-04-26 11:00 ET (end of the 10am candle).
func parseCandleEndFromSlug(slug string) time.Time {
	slug = strings.TrimPrefix(slug, "bitcoin-up-or-down-")
	slug = strings.TrimSuffix(slug, "-et")

	parts := strings.Split(slug, "-")
	if len(parts) < 4 {
		return time.Time{}
	}

	monthStr := parts[0]
	dayStr := parts[1]
	yearStr := parts[2]
	hourStr := parts[3]

	months := map[string]time.Month{
		"january": time.January, "february": time.February, "march": time.March,
		"april": time.April, "may": time.May, "june": time.June,
		"july": time.July, "august": time.August, "september": time.September,
		"october": time.October, "november": time.November, "december": time.December,
	}

	month, ok := months[monthStr]
	if !ok {
		return time.Time{}
	}
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return time.Time{}
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return time.Time{}
	}

	hour := parseHourStr(hourStr)
	if hour < 0 {
		return time.Time{}
	}

	// End of the candle = start hour + 1
	start := time.Date(year, month, day, hour, 0, 0, 0, etLocation)
	return start.Add(time.Hour)
}

func parseHourStr(s string) int {
	s = strings.ToLower(s)
	isPM := strings.HasSuffix(s, "pm")
	isAM := strings.HasSuffix(s, "am")
	if !isPM && !isAM {
		return -1
	}
	numStr := strings.TrimSuffix(strings.TrimSuffix(s, "pm"), "am")
	n, err := strconv.Atoi(numStr)
	if err != nil || n < 1 || n > 12 {
		return -1
	}
	if isAM {
		if n == 12 {
			return 0
		}
		return n
	}
	if n == 12 {
		return 12
	}
	return n + 12
}

// --- Auto-Trading Loop ---

func RunUpDownStrategy(ctx context.Context, cfg UpDownConfig, cb UpDownSignalCallback) error {
	db, err := sql.Open("sqlite", cfg.DBPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return fmt.Errorf("updown db open: %w", err)
	}
	defer db.Close()

	if err := initUpDownDB(db); err != nil {
		return fmt.Errorf("updown db init: %w", err)
	}

	slog.Info("updown_strategy.ready",
		"interval", cfg.ScanInterval.String(),
		"min_confidence", cfg.MinConfidence,
		"size_usd", cfg.SizeUSD,
		"max_daily", cfg.MaxDailyBets,
	)

	scanCount := 0
	scan := func() {
		resolveSettledBets(ctx, db)

		if err := updownScanOnce(ctx, db, cfg, cb); err != nil {
			slog.Warn("updown_strategy.scan_fail", "err", err.Error())
		}

		scanCount++
		if scanCount%6 == 0 {
			logPnLSummary(ctx, db)
		}
	}

	scan()

	tk := time.NewTicker(cfg.ScanInterval)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tk.C:
			scan()
		}
	}
}

func updownScanOnce(ctx context.Context, db *sql.DB, cfg UpDownConfig, cb UpDownSignalCallback) error {
	markets, err := DiscoverUpDownMarkets(ctx)
	if err != nil {
		return fmt.Errorf("discover: %w", err)
	}

	now := time.Now()
	minEnd := now.Add(1 * time.Hour)
	maxEnd := now.Add(4 * time.Hour)

	logUpDownPrices(ctx, db, markets)

	var candidates []UpDownMarket
	for _, m := range markets {
		if m.EndDate.Before(minEnd) || m.EndDate.After(maxEnd) {
			continue
		}
		if slugAlreadyBet(ctx, db, m.Slug) {
			continue
		}
		candidates = append(candidates, m)
	}

	if len(candidates) == 0 {
		slog.Info("updown_strategy.no_candidates",
			"total_markets", len(markets),
			"window", fmt.Sprintf("%s-%s", minEnd.Format("15:04"), maxEnd.Format("15:04")),
		)
		return nil
	}

	dailyBets := countBetsToday(ctx, db)
	if dailyBets >= cfg.MaxDailyBets {
		slog.Info("updown_strategy.daily_limit", "bets_today", dailyBets, "max", cfg.MaxDailyBets)
		return nil
	}

	candles, err := FetchCandles(ctx, "BTCUSDT", Interval1h, 200)
	if err != nil {
		return fmt.Errorf("fetch spot: %w", err)
	}
	spot := 0.0
	if len(candles) > 0 {
		spot = candles[len(candles)-1].Close
	}

	regime, rName, rConf := DetectCurrentRegime(candles)
	slog.Info("updown_strategy.regime",
		"regime", rName,
		"regime_conf", fmt.Sprintf("%.2f", rConf),
	)

	pred, err := PredictHourlyDirection(ctx)
	if err != nil {
		slog.Warn("updown_strategy.predict_fail", "err", err.Error())
	}

	for _, m := range candidates {
		if dailyBets >= cfg.MaxDailyBets {
			break
		}

		// Value betting: buy whichever side PM underprices.
		// BTC 1h up/down is ~50/50, so fair price for each side is ~0.50.
		// Buy the side that's cheaper (further below 0.50).
		direction := "Up"
		tokenID := m.UpTokenID
		pmPrice := m.UpPrice
		edge := 0.50 - m.UpPrice

		if m.DownPrice < m.UpPrice {
			direction = "Down"
			tokenID = m.DownTokenID
			pmPrice = m.DownPrice
			edge = 0.50 - m.DownPrice
		}

		// Markov tiebreaker: if model has directional conviction, follow it
		if pred.Confidence >= cfg.MinConfidence {
			if pred.Direction == "Up" && m.UpPrice <= 0.52 {
				direction = "Up"
				tokenID = m.UpTokenID
				pmPrice = m.UpPrice
				edge = 0.50 - m.UpPrice
			} else if pred.Direction == "Down" && m.DownPrice <= 0.52 {
				direction = "Down"
				tokenID = m.DownTokenID
				pmPrice = m.DownPrice
				edge = 0.50 - m.DownPrice
			}
		}

		// Skip if no value edge (both sides priced at or above fair value)
		if pmPrice > 0.49 {
			slog.Info("updown_strategy.skip_no_edge",
				"slug", m.Slug,
				"up_price", fmt.Sprintf("%.3f", m.UpPrice),
				"down_price", fmt.Sprintf("%.3f", m.DownPrice),
			)
			continue
		}

		// Skip VOLATILE regime — consistently bad for directional bets
		if regime == RegimeVolat && rConf >= 0.6 {
			slog.Info("updown_strategy.skip_volatile",
				"slug", m.Slug,
				"regime_conf", fmt.Sprintf("%.2f", rConf),
			)
			continue
		}

		conf := edge * 2 // normalize edge to 0-1 scale
		if conf < 0 {
			conf = 0
		}
		if conf > 1 {
			conf = 1
		}

		// Kelly sizing: fair prob ≈ 0.50, PM price = pmPrice
		betSize := KellySizeUSD(90.0, pmPrice, 0.50, cfg.SizeUSD*3)
		if betSize < 1.0 {
			betSize = cfg.SizeUSD // fallback to default
		}

		ev := ExpectedValue(pmPrice, 0.50)

		usePred := DirectionPrediction{
			Direction:  direction,
			Confidence: conf,
			Alignment:  "VALUE",
		}

		sig := UpDownSignal{
			MarketSlug:         m.Slug,
			Question:           m.Question,
			ConditionID:        m.ConditionID,
			TokenID:            tokenID,
			PMPrice:            pmPrice,
			PredictedDirection: direction,
			Confidence:         conf,
			Spot:               spot,
			SizeUSD:            betSize,
		}

		cb(sig)
		recordBet(ctx, db, m, usePred, spot, betSize, tokenID)
		dailyBets++

		slog.Info("updown_strategy.signal",
			"slug", m.Slug,
			"direction", direction,
			"pm_price", fmt.Sprintf("%.3f", pmPrice),
			"edge_pp", fmt.Sprintf("%.1f", edge*100),
			"ev_per_dollar", fmt.Sprintf("%.3f", ev),
			"kelly_size", fmt.Sprintf("%.2f", betSize),
			"regime", rName,
			"spot", fmt.Sprintf("%.0f", spot),
		)
	}

	return nil
}
