package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type whaleTrade struct {
	TS      string  `json:"ts"`
	Wallet  string  `json:"wallet"`
	Label   string  `json:"label"`
	Side    string  `json:"side"`
	Market  string  `json:"market"`
	Outcome string  `json:"outcome"`
	Price   float64 `json:"price"`
	Size    float64 `json:"size"`
	AssetID string  `json:"asset_id"`
	TradeID string  `json:"trade_id"`
	Action  string  `json:"action"`
	List    string  `json:"list"`
	Tier    string  `json:"tier"`
	Smart   float64 `json:"smart"`
	Bot     float64 `json:"bot"`

	Time time.Time `json:"-"`
}

type tapeTrade struct {
	Time     time.Time `json:"time"`
	Wallet   string    `json:"wallet"`
	Side     string    `json:"side"`
	Notional float64   `json:"notional"`
	Price    float64   `json:"price"`
	Outcome  string    `json:"outcome"`
	Market   string    `json:"market"`
	Asset    string    `json:"asset"`
	Tx       string    `json:"transaction"`
	List     string    `json:"known_list,omitempty"`
	Tier     string    `json:"tier,omitempty"`
	Smart    float64   `json:"smart,omitempty"`
	Bot      float64   `json:"bot,omitempty"`
}

type edgeSnapshot struct {
	SignalID    string  `json:"signal_id"`
	Wallet      string  `json:"wallet"`
	Label       string  `json:"label,omitempty"`
	List        string  `json:"list,omitempty"`
	Tier        string  `json:"tier,omitempty"`
	Market      string  `json:"market"`
	Outcome     string  `json:"outcome"`
	AssetID     string  `json:"asset_id"`
	TradeID     string  `json:"trade_id,omitempty"`
	Action      string  `json:"action"`
	TradeTime   string  `json:"trade_time"`
	EntryPrice  float64 `json:"entry_price"`
	NotionalUSD float64 `json:"notional_usd"`
	HorizonSec  int64   `json:"horizon_sec"`
	SampleTime  string  `json:"sample_time"`
	Mid         float64 `json:"mid"`
	DeltaPP     float64 `json:"delta_pp"`
	DelaySec    int64   `json:"delay_sec"`
}

type walletEdgeStats struct {
	Wallet      string
	Label       string
	List        string
	Tier        string
	Samples     int
	Wins        int
	DeltaPPSum  float64
	NotionalSum float64
	LastTime    time.Time
}

type walletMeta struct {
	List string
	Tier string
}

func main() {
	logPath := flag.String("log", "db/journal/whale_trades.jsonl", "whale trades JSONL")
	tapeLogPath := flag.String("tape_log", "", "optional sports-tape JSONL to edge-track observation-only large orders")
	walletsPath := flag.String("wallets", "wallets.strategy-push.txt", "optional comma-separated wallet allowlist/meta files")
	snapshotPath := flag.String("snapshots", "db/strategy_iteration/whale_edge_snapshots.jsonl", "edge snapshot JSONL")
	reportPath := flag.String("report", "reports/whale_edge.md", "edge report path")
	actionsRaw := flag.String("actions", "alert,followed", "BUY actions to track")
	horizonsRaw := flag.String("horizons", "0m,5m,15m,30m,60m", "comma-separated post-trade horizons")
	tolerance := flag.Duration("tolerance", 2*time.Minute, "sample only while horizon is within this lateness window")
	minNotional := flag.Float64("min_notional", 500, "minimum BUY notional")
	timeout := flag.Duration("timeout", 20*time.Second, "overall midpoint fetch timeout")
	flag.Parse()

	horizons, err := parseDurations(*horizonsRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "whale-edge: parse horizons: %v\n", err)
		os.Exit(1)
	}
	wallets, err := loadWalletMetas(*walletsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "whale-edge: load wallets: %v\n", err)
		os.Exit(1)
	}
	actions := parseSet(*actionsRaw)
	trades, err := loadTrades(*logPath, wallets, actions, *minNotional)
	if err != nil {
		fmt.Fprintf(os.Stderr, "whale-edge: load trades: %v\n", err)
		os.Exit(1)
	}
	tapeTrades, err := loadTapeTrades(*tapeLogPath, wallets, *minNotional)
	if err != nil {
		fmt.Fprintf(os.Stderr, "whale-edge: load tape trades: %v\n", err)
		os.Exit(1)
	}
	trades = append(trades, tapeTrades...)
	sort.Slice(trades, func(i, j int) bool { return trades[i].Time.Before(trades[j].Time) })
	existing, err := loadSnapshotKeys(*snapshotPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "whale-edge: load snapshots: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	newSnaps, err := collectDueSnapshots(ctx, trades, existing, horizons, *tolerance, time.Now(), fetchMidpoint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "whale-edge: collect snapshots: %v\n", err)
		os.Exit(1)
	}
	if err := appendSnapshots(*snapshotPath, newSnaps); err != nil {
		fmt.Fprintf(os.Stderr, "whale-edge: write snapshots: %v\n", err)
		os.Exit(1)
	}
	allSnaps, err := loadSnapshots(*snapshotPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "whale-edge: reload snapshots: %v\n", err)
		os.Exit(1)
	}
	applySnapshotMetas(allSnaps, wallets)
	if err := writeReport(*reportPath, allSnaps, horizons, len(trades), len(tapeTrades), len(newSnaps)); err != nil {
		fmt.Fprintf(os.Stderr, "whale-edge: write report: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("whale-edge done: trades=%d tape_trades=%d new_snapshots=%d total_snapshots=%d report=%s\n", len(trades), len(tapeTrades), len(newSnaps), len(allSnaps), *reportPath)
}

func parseDurations(raw string) ([]time.Duration, error) {
	var out []time.Duration
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		d, err := time.ParseDuration(part)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out, nil
}

func parseSet(raw string) map[string]struct{} {
	out := map[string]struct{}{}
	for _, part := range strings.Split(raw, ",") {
		part = strings.ToLower(strings.TrimSpace(part))
		if part != "" {
			out[part] = struct{}{}
		}
	}
	return out
}

func loadWalletMetas(path string) (map[string]walletMeta, error) {
	if path == "" {
		return nil, nil
	}
	if strings.Contains(path, ",") {
		out := map[string]walletMeta{}
		for _, part := range strings.Split(path, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			metas, err := loadWalletMetas(part)
			if err != nil {
				return nil, err
			}
			for addr, meta := range metas {
				out[addr] = meta
			}
		}
		return out, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	out := map[string]walletMeta{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "#", 2)
		fields := strings.Fields(parts[0])
		if len(fields) == 0 {
			continue
		}
		addr := strings.ToLower(fields[0])
		if strings.HasPrefix(addr, "0x") && len(addr) == 42 {
			meta := walletMeta{}
			if len(parts) > 1 {
				for _, field := range strings.Fields(parts[1]) {
					k, v, ok := strings.Cut(field, "=")
					if !ok {
						continue
					}
					switch k {
					case "list":
						meta.List = v
					case "tier":
						meta.Tier = v
					}
				}
			}
			out[addr] = meta
		}
	}
	return out, sc.Err()
}

func loadTapeTrades(path string, wallets map[string]walletMeta, minNotional float64) ([]whaleTrade, error) {
	if path == "" {
		return nil, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var out []whaleTrade
	seen := map[string]struct{}{}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1024*1024), 8*1024*1024)
	for sc.Scan() {
		var tape tapeTrade
		if err := json.Unmarshal(sc.Bytes(), &tape); err != nil {
			continue
		}
		wallet := strings.ToLower(strings.TrimSpace(tape.Wallet))
		side := strings.ToUpper(strings.TrimSpace(tape.Side))
		if wallet == "" || side != "BUY" || tape.Asset == "" || tape.Price <= 0 || tape.Notional < minNotional || tape.Time.IsZero() {
			continue
		}
		meta := walletMeta{}
		if len(wallets) > 0 {
			var ok bool
			meta, ok = wallets[wallet]
			if !ok {
				continue
			}
		}
		tr := whaleTrade{
			TS:      tape.Time.Format(time.RFC3339),
			Wallet:  wallet,
			Label:   "sports_tape",
			Side:    side,
			Market:  tape.Market,
			Outcome: tape.Outcome,
			Price:   tape.Price,
			Size:    tape.Notional,
			AssetID: tape.Asset,
			TradeID: tape.Tx,
			Action:  "tape",
			List:    firstNonEmpty(meta.List, tape.List, "tape_observe"),
			Tier:    firstNonEmpty(meta.Tier, tape.Tier),
			Smart:   tape.Smart,
			Bot:     tape.Bot,
			Time:    tape.Time,
		}
		id := signalID(tr)
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, tr)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Time.Before(out[j].Time) })
	return out, nil
}

func loadTrades(path string, wallets map[string]walletMeta, actions map[string]struct{}, minNotional float64) ([]whaleTrade, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var out []whaleTrade
	seen := map[string]struct{}{}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1024*1024), 8*1024*1024)
	for sc.Scan() {
		var tr whaleTrade
		if err := json.Unmarshal(sc.Bytes(), &tr); err != nil {
			continue
		}
		ts, err := time.Parse(time.RFC3339, tr.TS)
		if err != nil {
			continue
		}
		tr.Time = ts
		tr.Wallet = strings.ToLower(strings.TrimSpace(tr.Wallet))
		tr.Side = strings.ToUpper(strings.TrimSpace(tr.Side))
		tr.Action = strings.ToLower(strings.TrimSpace(tr.Action))
		if tr.Wallet == "" || tr.Side != "BUY" || tr.AssetID == "" || tr.Price <= 0 || tr.Size < minNotional {
			continue
		}
		if len(wallets) > 0 {
			meta, ok := wallets[tr.Wallet]
			if !ok {
				continue
			}
			if meta.List != "" {
				tr.List = meta.List
			}
			if meta.Tier != "" {
				tr.Tier = meta.Tier
			}
		}
		if _, ok := actions[tr.Action]; !ok {
			continue
		}
		id := signalID(tr)
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, tr)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Time.Before(out[j].Time) })
	return out, nil
}

func applySnapshotMetas(snaps []edgeSnapshot, metas map[string]walletMeta) {
	if len(metas) == 0 {
		return
	}
	for i := range snaps {
		meta, ok := metas[strings.ToLower(strings.TrimSpace(snaps[i].Wallet))]
		if !ok {
			continue
		}
		if meta.List != "" {
			snaps[i].List = meta.List
		}
		if meta.Tier != "" {
			snaps[i].Tier = meta.Tier
		}
	}
}

func signalID(tr whaleTrade) string {
	if tr.TradeID != "" {
		return tr.Wallet + "|" + tr.AssetID + "|" + tr.TradeID
	}
	return fmt.Sprintf("%s|%s|%d|%.8f|%.4f", tr.Wallet, tr.AssetID, tr.Time.Unix(), tr.Price, tr.Size)
}

func snapshotKey(signalID string, horizon time.Duration) string {
	return fmt.Sprintf("%s|%d", signalID, int64(horizon.Seconds()))
}

func loadSnapshotKeys(path string) (map[string]struct{}, error) {
	snaps, err := loadSnapshots(path)
	if err != nil {
		return nil, err
	}
	out := map[string]struct{}{}
	for _, s := range snaps {
		out[snapshotKey(s.SignalID, time.Duration(s.HorizonSec)*time.Second)] = struct{}{}
	}
	return out, nil
}

func loadSnapshots(path string) ([]edgeSnapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var out []edgeSnapshot
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1024*1024), 8*1024*1024)
	for sc.Scan() {
		var s edgeSnapshot
		if err := json.Unmarshal(sc.Bytes(), &s); err == nil && s.SignalID != "" {
			out = append(out, s)
		}
	}
	return out, sc.Err()
}

func collectDueSnapshots(ctx context.Context, trades []whaleTrade, existing map[string]struct{}, horizons []time.Duration, tolerance time.Duration, now time.Time, fetchMid func(context.Context, string) (float64, error)) ([]edgeSnapshot, error) {
	var out []edgeSnapshot
	midCache := map[string]float64{}
	for _, tr := range trades {
		id := signalID(tr)
		for _, h := range horizons {
			key := snapshotKey(id, h)
			if _, ok := existing[key]; ok {
				continue
			}
			age := now.Sub(tr.Time)
			if age < h || age > h+tolerance {
				continue
			}
			mid, ok := midCache[tr.AssetID]
			if !ok {
				var err error
				mid, err = fetchMid(ctx, tr.AssetID)
				if err != nil {
					return out, err
				}
				midCache[tr.AssetID] = mid
			}
			out = append(out, edgeSnapshot{
				SignalID:    id,
				Wallet:      tr.Wallet,
				Label:       tr.Label,
				List:        tr.List,
				Tier:        tr.Tier,
				Market:      tr.Market,
				Outcome:     tr.Outcome,
				AssetID:     tr.AssetID,
				TradeID:     tr.TradeID,
				Action:      tr.Action,
				TradeTime:   tr.Time.Format(time.RFC3339),
				EntryPrice:  tr.Price,
				NotionalUSD: tr.Size,
				HorizonSec:  int64(h.Seconds()),
				SampleTime:  now.Format(time.RFC3339),
				Mid:         mid,
				DeltaPP:     (mid - tr.Price) * 100,
				DelaySec:    int64(age.Seconds()) - int64(h.Seconds()),
			})
		}
	}
	return out, nil
}

func fetchMidpoint(ctx context.Context, assetID string) (float64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://clob.polymarket.com/midpoint?token_id="+url.QueryEscape(assetID), nil)
	if err != nil {
		return 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var body struct {
		Mid string `json:"mid"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0, err
	}
	var mid float64
	if _, err := fmt.Sscanf(body.Mid, "%f", &mid); err != nil {
		return 0, err
	}
	return mid, nil
}

func appendSnapshots(path string, snaps []edgeSnapshot) error {
	if len(snaps) == 0 {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, s := range snaps {
		if err := enc.Encode(s); err != nil {
			return err
		}
	}
	return nil
}

func writeReport(path string, snaps []edgeSnapshot, horizons []time.Duration, trackedTrades, trackedTapeTrades, newSnapshots int) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# Whale Edge Report\n\n")
	fmt.Fprintf(&b, "**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04 MST"))
	fmt.Fprintf(&b, "- Snapshots: %d\n", len(snaps))
	fmt.Fprintf(&b, "- New snapshots: %d\n", newSnapshots)
	fmt.Fprintf(&b, "- Tracked BUY trades: %d\n", trackedTrades)
	fmt.Fprintf(&b, "- Tracked sports-tape trades: %d\n", trackedTapeTrades)
	fmt.Fprintf(&b, "- Horizons: %s\n\n", formatHorizons(horizons))
	for _, h := range horizons {
		writeHorizonSection(&b, snaps, h)
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func writeHorizonSection(b *strings.Builder, snaps []edgeSnapshot, horizon time.Duration) {
	stats := map[string]*walletEdgeStats{}
	for _, s := range snaps {
		if s.HorizonSec != int64(horizon.Seconds()) {
			continue
		}
		st := stats[s.Wallet]
		if st == nil {
			st = &walletEdgeStats{Wallet: s.Wallet, Label: s.Label, List: s.List, Tier: s.Tier}
			stats[s.Wallet] = st
		}
		st.Samples++
		if s.DeltaPP > 0 {
			st.Wins++
		}
		st.DeltaPPSum += s.DeltaPP
		st.NotionalSum += s.NotionalUSD
		if t, err := time.Parse(time.RFC3339, s.SampleTime); err == nil && t.After(st.LastTime) {
			st.LastTime = t
		}
	}
	rows := make([]*walletEdgeStats, 0, len(stats))
	for _, st := range stats {
		rows = append(rows, st)
	}
	sort.Slice(rows, func(i, j int) bool {
		ai := rows[i].DeltaPPSum / float64(maxInt(rows[i].Samples, 1))
		aj := rows[j].DeltaPPSum / float64(maxInt(rows[j].Samples, 1))
		if ai != aj {
			return ai > aj
		}
		return rows[i].Samples > rows[j].Samples
	})
	fmt.Fprintf(b, "## %s Edge\n\n", horizon)
	fmt.Fprintf(b, "| Wallet | List | Tier | Samples | Win | AvgDeltaPP | Notional | Last |\n")
	fmt.Fprintf(b, "|---|---|---|---:|---:|---:|---:|---|\n")
	if len(rows) == 0 {
		fmt.Fprintf(b, "| n/a |  |  | 0 | 0.0%% | 0.00 | $0 |  |\n\n")
		return
	}
	limit := minInt(len(rows), 20)
	for i := 0; i < limit; i++ {
		st := rows[i]
		avg := st.DeltaPPSum / float64(st.Samples)
		fmt.Fprintf(b, "| `%s` | %s | %s | %d | %.1f%% | %+.2f | $%.0f | %s |\n",
			shortAddr(st.Wallet), st.List, st.Tier, st.Samples, pct(float64(st.Wins), float64(st.Samples)), avg, st.NotionalSum, formatShortTime(st.LastTime))
	}
	fmt.Fprintf(b, "\n")
}

func formatHorizons(horizons []time.Duration) string {
	parts := make([]string, 0, len(horizons))
	for _, h := range horizons {
		parts = append(parts, h.String())
	}
	return strings.Join(parts, ",")
}

func pct(n, d float64) float64 {
	if d <= 0 {
		return 0
	}
	return n / d * 100
}

func shortAddr(addr string) string {
	if len(addr) <= 12 {
		return addr
	}
	return addr[:6] + "..." + addr[len(addr)-4:]
}

func formatShortTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("01-02 15:04")
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		v = strings.TrimSpace(v)
		if v != "" {
			return v
		}
	}
	return ""
}
