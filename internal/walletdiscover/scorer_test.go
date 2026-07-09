package walletdiscover

import (
	"fmt"
	"testing"
)

func TestScoreWallet_BotLikeExtremeFixedFlow(t *testing.T) {
	cfg := DefaultConfig()
	var trades []Trade
	for i := 0; i < 40; i++ {
		side := "BUY"
		if i%2 == 0 {
			side = "SELL"
		}
		trades = append(trades, Trade{
			ProxyWallet: "0x2222222222222222222222222222222222222222",
			Side:        side,
			Type:        "TRADE",
			Asset:       "asset",
			ConditionID: "0xmarket",
			Size:        1000,
			Price:       0.999,
			Timestamp:   1000 + int64(i%8),
			Title:       "Will Trump say something?",
		})
	}
	score := ScoreWallet("0x2222222222222222222222222222222222222222", nil, trades, nil, cfg)
	if score.Tier != "BOT" {
		t.Fatalf("tier=%s, want BOT (bot_score=%.2f)", score.Tier, score.BotScore)
	}
	if score.BotScore < 60 {
		t.Fatalf("bot score %.2f too low", score.BotScore)
	}
}

func TestScoreWallet_StrongLowBotWallet(t *testing.T) {
	cfg := DefaultConfig()
	cand := &Candidate{Sources: map[string]int{"holder": 3, "recent_trade": 5}}
	var trades []Trade
	for i := 0; i < 30; i++ {
		price := 0.35 + float64(i%5)*0.03
		trades = append(trades, Trade{
			ProxyWallet: "0x3333333333333333333333333333333333333333",
			Side:        "BUY",
			Type:        "TRADE",
			Asset:       "asset",
			ConditionID: "0xmarket",
			Size:        300 + float64(i),
			Price:       price,
			Timestamp:   1000 + int64(i*3600),
			Title:       "Will Arsenal FC win on 2026-05-18?",
		})
	}
	var closed []ClosedPosition
	for i := 0; i < 15; i++ {
		closed = append(closed, ClosedPosition{
			TotalBought: 500,
			RealizedPnL: 80,
		})
	}
	score := ScoreWallet("0x3333333333333333333333333333333333333333", cand, trades, closed, cfg)
	if score.Tier != "A" && score.Tier != "B" {
		t.Fatalf("tier=%s, want A/B (bot=%.2f edge=%.2f)", score.Tier, score.BotScore, score.Edge)
	}
	if score.BotScore >= 35 {
		t.Fatalf("bot score %.2f too high", score.BotScore)
	}
}

func TestScoreWallet_BotRiskFlagsRejectAutoFollow(t *testing.T) {
	cfg := DefaultConfig()
	var trades []Trade
	for i := 0; i < 50; i++ {
		trades = append(trades, Trade{
			ProxyWallet: "0x4444444444444444444444444444444444444444",
			Side:        "BUY",
			Type:        "TRADE",
			ConditionID: "0xmarket",
			Size:        100,
			Price:       0.99,
			Timestamp:   1000 + int64(i%5),
			Title:       "Bitcoin up or down - 15m",
		})
	}
	score := ScoreWallet("0x4444444444444444444444444444444444444444", nil, trades, nil, cfg)
	if score.Tier != "BOT" {
		t.Fatalf("tier=%s, want BOT", score.Tier)
	}
	if score.FollowAction != "reject-bot" {
		t.Fatalf("follow_action=%s, want reject-bot", score.FollowAction)
	}
	if len(score.RiskFlags) == 0 {
		t.Fatalf("expected risk flags")
	}
}

func TestScoreWallet_SmartMoneyActionForFocusedWinner(t *testing.T) {
	cfg := DefaultConfig()
	cand := &Candidate{Sources: map[string]int{"holder": 5, "recent_trade": 8, "existing": 1}}
	var trades []Trade
	for i := 0; i < 45; i++ {
		asset := fmt.Sprintf("asset-%02d", i)
		price := 0.25 + float64(i%6)*0.04
		trades = append(trades, Trade{
			ProxyWallet: "0x5555555555555555555555555555555555555555",
			Side:        "BUY",
			Type:        "TRADE",
			Asset:       asset,
			ConditionID: "0xmarket",
			Size:        250 + float64(i%7)*15,
			Price:       price,
			Timestamp:   1000 + int64(i*7200),
			Title:       "Will Arsenal FC beat Chelsea?",
			Slug:        "arsenal-chelsea",
		})
		if i < 18 {
			trades = append(trades, Trade{
				ProxyWallet: "0x5555555555555555555555555555555555555555",
				Side:        "SELL",
				Type:        "TRADE",
				Asset:       asset,
				ConditionID: "0xmarket",
				Size:        250 + float64(i%7)*15,
				Price:       price + 0.18,
				Timestamp:   1000 + int64(i*7200) + 3600,
				Title:       "Will Arsenal FC beat Chelsea?",
				Slug:        "arsenal-chelsea",
			})
		}
	}
	var closed []ClosedPosition
	for i := 0; i < 24; i++ {
		pnl := 90.0
		if i%5 == 0 {
			pnl = -40
		}
		closed = append(closed, ClosedPosition{
			TotalBought: 400,
			RealizedPnL: pnl,
		})
	}
	score := ScoreWallet("0x5555555555555555555555555555555555555555", cand, trades, closed, cfg)
	if score.Tier != "A" {
		t.Fatalf("tier=%s, want A (smart=%.2f bot=%.2f edge=%.2f)", score.Tier, score.SmartMoneyScore, score.BotScore, score.Edge)
	}
	if score.FollowAction != "auto-small" && score.FollowAction != "prompt" {
		t.Fatalf("follow_action=%s, want auto-small or prompt", score.FollowAction)
	}
	if score.SmartMoneyScore < 70 {
		t.Fatalf("smart score %.2f too low", score.SmartMoneyScore)
	}
}
