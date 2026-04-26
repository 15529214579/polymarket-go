package main

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

func runBTCPnL(btcDBDir string) error {
	dbPath := filepath.Join(btcDBDir, "btc.db")
	uri := fmt.Sprintf("file:%s?mode=ro&_journal_mode=WAL", dbPath)
	db, err := sql.Open("sqlite", uri)
	if err != nil {
		return fmt.Errorf("open btc.db: %w", err)
	}
	defer db.Close()

	ctx := context.Background()

	fmt.Println("════════════════════════════════════════")
	fmt.Println(" BTC Up/Down — Live PnL Report")
	fmt.Printf(" db: %s\n", dbPath)
	fmt.Println("════════════════════════════════════════")

	if err := printBetSummary(ctx, db); err != nil {
		return err
	}
	if err := printPricingEfficiency(ctx, db); err != nil {
		return err
	}
	return nil
}

type liveBet struct {
	id        int64
	ts        int64
	slug      string
	predicted string
	upPrice   float64
	downPrice float64
	spot      float64
	sizeUSD   float64
	actual    sql.NullString
	pnl       sql.NullFloat64
	resolved  sql.NullInt64
}

func printBetSummary(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx,
		`SELECT id, timestamp, slug, predicted_direction,
		        pm_up_price, pm_down_price, btc_spot, size_usd,
		        actual_direction, pnl, resolved_at
		 FROM updown_bets ORDER BY timestamp`)
	if err != nil {
		return fmt.Errorf("query bets: %w", err)
	}
	defer rows.Close()

	var bets []liveBet
	for rows.Next() {
		var b liveBet
		if err := rows.Scan(&b.id, &b.ts, &b.slug, &b.predicted,
			&b.upPrice, &b.downPrice, &b.spot, &b.sizeUSD,
			&b.actual, &b.pnl, &b.resolved); err != nil {
			continue
		}
		bets = append(bets, b)
	}

	fmt.Printf("\n── Bet Log (%d total) ──\n", len(bets))
	fmt.Printf("  %-50s %-5s %-5s %8s %8s %10s\n",
		"slug", "pred", "act", "pm_px", "size", "pnl")

	var totalPnL, totalCapital float64
	var wins, losses, pending int
	for _, b := range bets {
		pmPrice := b.downPrice
		if b.predicted == "Up" {
			pmPrice = b.upPrice
		}

		actualStr := "..."
		pnlStr := "pending"
		if b.actual.Valid {
			actualStr = b.actual.String
			if b.pnl.Valid {
				pnlStr = fmt.Sprintf("$%+.2f", b.pnl.Float64)
				totalPnL += b.pnl.Float64
				if b.pnl.Float64 > 0 {
					wins++
				} else {
					losses++
				}
			}
		} else {
			pending++
		}
		totalCapital += b.sizeUSD

		shortSlug := b.slug
		if len(shortSlug) > 50 {
			shortSlug = shortSlug[:47] + "..."
		}
		fmt.Printf("  %-50s %-5s %-5s %8.3f %8.2f %10s\n",
			shortSlug, b.predicted, actualStr, pmPrice, b.sizeUSD, pnlStr)
	}

	total := wins + losses
	wr := 0.0
	if total > 0 {
		wr = float64(wins) / float64(total) * 100
	}

	fmt.Printf("\n── Summary ──\n")
	fmt.Printf("  Total bets:    %d (resolved: %d, pending: %d)\n", len(bets), total, pending)
	fmt.Printf("  Wins:          %d\n", wins)
	fmt.Printf("  Losses:        %d\n", losses)
	fmt.Printf("  Win rate:      %.1f%%\n", wr)
	fmt.Printf("  Total PnL:     $%+.2f\n", totalPnL)
	fmt.Printf("  Capital used:  $%.2f\n", totalCapital)
	if totalCapital > 0 {
		fmt.Printf("  ROI:           %+.1f%%\n", totalPnL/totalCapital*100)
	}

	if len(bets) > 0 {
		first := time.Unix(bets[0].ts, 0)
		last := time.Unix(bets[len(bets)-1].ts, 0)
		fmt.Printf("  Period:        %s to %s\n", first.Format("Jan 02 15:04"), last.Format("Jan 02 15:04"))
	}

	return nil
}

func printPricingEfficiency(ctx context.Context, db *sql.DB) error {
	var count int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM updown_prices").Scan(&count)
	if count == 0 {
		fmt.Println("\n── PM Pricing Efficiency ──")
		fmt.Println("  No price data yet (updown_prices table empty)")
		return nil
	}

	fmt.Printf("\n── PM Pricing Efficiency (%d snapshots) ──\n", count)

	rows, err := db.QueryContext(ctx,
		`SELECT
			COUNT(*) as n,
			AVG(up_price) as avg_up,
			AVG(down_price) as avg_down,
			AVG(spread) as avg_spread,
			AVG(deviation) as avg_dev,
			MIN(up_price) as min_up,
			MIN(down_price) as min_down,
			MAX(ABS(deviation)) as max_dev
		FROM updown_prices`)
	if err != nil {
		return fmt.Errorf("query prices: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var n int
		var avgUp, avgDown, avgSpread, avgDev, minUp, minDown, maxDev float64
		rows.Scan(&n, &avgUp, &avgDown, &avgSpread, &avgDev, &minUp, &minDown, &maxDev)

		fmt.Printf("  Avg Up price:      %.4f\n", avgUp)
		fmt.Printf("  Avg Down price:    %.4f\n", avgDown)
		fmt.Printf("  Avg spread:        %.4f (ideal=0, vig=positive)\n", avgSpread)
		fmt.Printf("  Avg deviation:     %.4f (from 0.50 fair)\n", avgDev)
		fmt.Printf("  Min Up price:      %.4f\n", minUp)
		fmt.Printf("  Min Down price:    %.4f\n", minDown)
		fmt.Printf("  Max deviation:     %.4f\n", maxDev)

		if maxDev > 0.01 {
			fmt.Printf("  → Pricing shows mispricing up to %.1fpp — edge exists\n", maxDev*100)
		} else {
			fmt.Printf("  → Pricing is efficient (deviation < 1pp) — limited edge\n")
		}
	}

	// Distribution of deviations
	distRows, err := db.QueryContext(ctx,
		`SELECT
			CASE
				WHEN ABS(deviation) < 0.005 THEN '<0.5pp'
				WHEN ABS(deviation) < 0.01 THEN '0.5-1pp'
				WHEN ABS(deviation) < 0.02 THEN '1-2pp'
				WHEN ABS(deviation) < 0.03 THEN '2-3pp'
				ELSE '>3pp'
			END as bucket,
			COUNT(*) as n
		FROM updown_prices
		GROUP BY bucket
		ORDER BY MIN(ABS(deviation))`)
	if err == nil {
		defer distRows.Close()
		fmt.Printf("\n  %-10s %6s\n", "deviation", "count")
		for distRows.Next() {
			var bucket string
			var n int
			distRows.Scan(&bucket, &n)
			pct := float64(n) / float64(count) * 100
			bar := ""
			barLen := int(math.Round(pct / 2))
			for j := 0; j < barLen; j++ {
				bar += "█"
			}
			fmt.Printf("  %-10s %6d %5.1f%% %s\n", bucket, n, pct, bar)
		}
	}

	return nil
}
