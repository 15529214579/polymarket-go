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
	"time"
)

type OrderbookDepth struct {
	TokenID    string
	Strike     float64
	BidDepth   float64 // total bid size in top 5 levels (USDC)
	AskDepth   float64 // total ask size in top 5 levels (USDC)
	Spread     float64 // best ask - best bid (0..1 scale)
	MidPrice   float64 // (best bid + best ask) / 2
	DepthScore float64 // 0-1: 0 = illiquid, 1 = deep
}

type MarketDepth struct {
	Depths     []OrderbookDepth
	AvgScore   float64
	MinScore   float64
}

type clobOrder struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}

type clobBook struct {
	Bids []clobOrder `json:"bids"`
	Asks []clobOrder `json:"asks"`
}

func FetchOrderbookDepth(ctx context.Context, markets []PMMarket) MarketDepth {
	client := &http.Client{Timeout: 10 * time.Second}
	var depths []OrderbookDepth
	totalScore := 0.0
	minScore := 1.0

	for _, m := range markets {
		if len(m.ClobTokenIDs) == 0 {
			continue
		}

		tokenID := m.ClobTokenIDs[0] // YES token
		book, err := fetchCLOBBook(ctx, client, tokenID)
		if err != nil {
			slog.Debug("orderbook.fetch_fail", "strike", m.Strike, "err", err.Error())
			continue
		}

		depth := computeDepth(tokenID, m.Strike, book)
		depths = append(depths, depth)
		totalScore += depth.DepthScore
		if depth.DepthScore < minScore {
			minScore = depth.DepthScore
		}
	}

	md := MarketDepth{Depths: depths}
	if len(depths) > 0 {
		md.AvgScore = totalScore / float64(len(depths))
		md.MinScore = minScore
	}

	slog.Info("orderbook.depth",
		"markets_checked", len(depths),
		"avg_score", fmt.Sprintf("%.2f", md.AvgScore),
		"min_score", fmt.Sprintf("%.2f", md.MinScore),
	)

	return md
}

func fetchCLOBBook(ctx context.Context, client *http.Client, tokenID string) (clobBook, error) {
	url := fmt.Sprintf("https://clob.polymarket.com/book?token_id=%s", tokenID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return clobBook{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return clobBook{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 65536))
	if err != nil {
		return clobBook{}, err
	}

	var book clobBook
	if err := json.Unmarshal(body, &book); err != nil {
		return clobBook{}, fmt.Errorf("decode book: %w", err)
	}
	return book, nil
}

func computeDepth(tokenID string, strike float64, book clobBook) OrderbookDepth {
	d := OrderbookDepth{TokenID: tokenID, Strike: strike}

	topN := 5
	for i, b := range book.Bids {
		if i >= topN {
			break
		}
		price, _ := strconv.ParseFloat(b.Price, 64)
		size, _ := strconv.ParseFloat(b.Size, 64)
		d.BidDepth += price * size
		if i == 0 {
			d.MidPrice = price
		}
	}

	for i, a := range book.Asks {
		if i >= topN {
			break
		}
		price, _ := strconv.ParseFloat(a.Price, 64)
		size, _ := strconv.ParseFloat(a.Size, 64)
		d.AskDepth += price * size
		if i == 0 {
			if d.MidPrice > 0 {
				d.Spread = price - d.MidPrice
				d.MidPrice = (d.MidPrice + price) / 2
			} else {
				d.MidPrice = price
			}
		}
	}

	// Depth score: based on total liquidity and spread
	totalDepth := d.BidDepth + d.AskDepth
	switch {
	case totalDepth >= 500:
		d.DepthScore = 1.0
	case totalDepth >= 200:
		d.DepthScore = 0.8
	case totalDepth >= 50:
		d.DepthScore = 0.6
	case totalDepth >= 10:
		d.DepthScore = 0.4
	case totalDepth > 0:
		d.DepthScore = 0.2
	default:
		d.DepthScore = 0.0
	}

	// Penalize wide spreads
	if d.Spread > 0.05 {
		d.DepthScore *= 0.7
	} else if d.Spread > 0.02 {
		d.DepthScore *= 0.85
	}

	return d
}

func SaveOrderbookDepth(ctx context.Context, db *sql.DB, md MarketDepth) error {
	const ddl = `CREATE TABLE IF NOT EXISTS btc_orderbook (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp   INTEGER NOT NULL,
		token_id    TEXT,
		strike      REAL,
		bid_depth   REAL,
		ask_depth   REAL,
		spread      REAL,
		mid_price   REAL,
		depth_score REAL
	);`
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	now := time.Now().Unix()
	for _, d := range md.Depths {
		if _, err := db.ExecContext(ctx,
			`INSERT INTO btc_orderbook(timestamp, token_id, strike, bid_depth, ask_depth, spread, mid_price, depth_score)
			 VALUES(?,?,?,?,?,?,?,?)`,
			now, d.TokenID, d.Strike, d.BidDepth, d.AskDepth, d.Spread, d.MidPrice, d.DepthScore); err != nil {
			return err
		}
	}
	return nil
}

func DepthModifier(md MarketDepth, strike float64) float64 {
	for _, d := range md.Depths {
		if d.Strike == strike {
			if d.DepthScore < 0.3 {
				return 0.5 // heavily penalize illiquid markets
			}
			if d.DepthScore < 0.5 {
				return 0.75
			}
			return 1.0
		}
	}
	return 0.8 // no data = conservative
}
