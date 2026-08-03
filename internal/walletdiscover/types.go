package walletdiscover

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	defaultGammaBase = "https://gamma-api.polymarket.com"
	defaultDataBase  = "https://data-api.polymarket.com"
	defaultLBBase    = "https://lb-api.polymarket.com"
)

type Config struct {
	GammaBase       string
	DataBase        string
	LeaderboardBase string
	HTTPMaxAttempts int
	HTTPRetryBase   time.Duration
	HTTPRetryMax    time.Duration
	HTTPTimeout     time.Duration

	MarketsLimit         int
	TradesPages          int
	TradesLimit          int
	HoldersLimit         int
	ActivityPages        int
	ActivityLimit        int
	ClosedLimit          int
	Concurrency          int
	MaxCandidates        int
	Days                 int
	MinNotionalUSD       float64
	MinHolderShares      float64
	MinPrice             float64
	MaxPrice             float64
	CopyStakeUSD         float64
	CopySlippageBP       float64
	CopyFeeBP            float64
	LeaderboardLimit     int
	LeaderboardWindows   string
	LeaderboardKinds     string
	TargetCategories     string
	ExistingWallets      string
	SportsTapeWallets    string
	RetainWallets        string
	OutputDir            string
	ReportPath           string
	GeneratedTier        string
	GeneratedWalletsPath string
	AutoWalletsPath      string
	PromptWalletsPath    string
	PositiveWalletsPath  string
}

func DefaultConfig() Config {
	return Config{
		GammaBase:            defaultGammaBase,
		DataBase:             defaultDataBase,
		LeaderboardBase:      defaultLBBase,
		HTTPMaxAttempts:      4,
		HTTPRetryBase:        250 * time.Millisecond,
		HTTPRetryMax:         3 * time.Second,
		HTTPTimeout:          30 * time.Second,
		MarketsLimit:         200,
		TradesPages:          4,
		TradesLimit:          500,
		HoldersLimit:         100,
		ActivityPages:        4,
		ActivityLimit:        500,
		ClosedLimit:          200,
		Concurrency:          5,
		MaxCandidates:        500,
		Days:                 90,
		MinNotionalUSD:       100,
		MinHolderShares:      100,
		MinPrice:             0.05,
		MaxPrice:             0.95,
		CopyStakeUSD:         10,
		CopySlippageBP:       50,
		CopyFeeBP:            0,
		LeaderboardLimit:     0,
		LeaderboardWindows:   "7d,30d,all",
		LeaderboardKinds:     "profit,volume",
		TargetCategories:     "",
		ExistingWallets:      "wallets.txt",
		SportsTapeWallets:    "",
		RetainWallets:        "",
		OutputDir:            "db",
		ReportPath:           "reports/wallet_discovery.md",
		GeneratedTier:        "C",
		GeneratedWalletsPath: "wallets.generated.txt",
		AutoWalletsPath:      "wallets.smartmoney-auto.txt",
		PromptWalletsPath:    "wallets.smartmoney-prompt.txt",
		PositiveWalletsPath:  "wallets.strategy-positive.txt",
	}
}

type FlexFloat float64

func (f *FlexFloat) UnmarshalJSON(b []byte) error {
	s := strings.TrimSpace(string(b))
	if s == "" || s == "null" {
		*f = 0
		return nil
	}
	if s[0] == '"' {
		var raw string
		if err := json.Unmarshal(b, &raw); err != nil {
			return err
		}
		if raw == "" {
			*f = 0
			return nil
		}
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return fmt.Errorf("parse float %q: %w", raw, err)
		}
		*f = FlexFloat(v)
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("parse float %q: %w", s, err)
	}
	*f = FlexFloat(v)
	return nil
}

type Market struct {
	ID              string    `json:"id"`
	Question        string    `json:"question"`
	ConditionID     string    `json:"conditionId"`
	Slug            string    `json:"slug"`
	Category        string    `json:"category"`
	Active          bool      `json:"active"`
	Closed          bool      `json:"closed"`
	NegRisk         bool      `json:"negRisk"`
	EndDate         string    `json:"endDate"`
	Liquidity       FlexFloat `json:"liquidity"`
	LiquidityClob   FlexFloat `json:"liquidityClob"`
	Volume          FlexFloat `json:"volume"`
	Volume24hr      FlexFloat `json:"volume24hr"`
	ClobTokenIDsRaw string    `json:"clobTokenIds"`
	OutcomesRaw     string    `json:"outcomes"`
}

func (m Market) ClobTokenIDs() []string { return parseStringList(m.ClobTokenIDsRaw) }
func (m Market) Outcomes() []string     { return parseStringList(m.OutcomesRaw) }

func parseStringList(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}

func (m Market) ParsedEndDate() time.Time {
	if m.EndDate == "" {
		return time.Time{}
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05Z"} {
		t, err := time.Parse(layout, m.EndDate)
		if err == nil {
			return t
		}
	}
	return time.Time{}
}

func (m Market) LiquidityUSD() float64 {
	if m.LiquidityClob > 0 {
		return float64(m.LiquidityClob)
	}
	return float64(m.Liquidity)
}

type Trade struct {
	ProxyWallet     string  `json:"proxyWallet"`
	Side            string  `json:"side"`
	Asset           string  `json:"asset"`
	ConditionID     string  `json:"conditionId"`
	Size            float64 `json:"size"`
	Price           float64 `json:"price"`
	Timestamp       int64   `json:"timestamp"`
	Title           string  `json:"title"`
	Slug            string  `json:"slug"`
	EventSlug       string  `json:"eventSlug"`
	Outcome         string  `json:"outcome"`
	OutcomeIndex    int     `json:"outcomeIndex"`
	TransactionHash string  `json:"transactionHash"`
	Type            string  `json:"type"`
}

func (t Trade) NotionalUSD() float64 {
	if t.Price <= 0 {
		return t.Size
	}
	return t.Size * t.Price
}

type HolderResponse struct {
	Token   string   `json:"token"`
	Holders []Holder `json:"holders"`
}

type Holder struct {
	ProxyWallet string  `json:"proxyWallet"`
	Asset       string  `json:"asset"`
	Amount      float64 `json:"amount"`
	Name        string  `json:"name"`
	Pseudonym   string  `json:"pseudonym"`
	Verified    bool    `json:"verified"`
}

type LeaderboardEntry struct {
	ProxyWallet           string  `json:"proxyWallet"`
	Amount                float64 `json:"amount"`
	Pseudonym             string  `json:"pseudonym"`
	Name                  string  `json:"name"`
	Bio                   string  `json:"bio"`
	ProfileImage          string  `json:"profileImage"`
	ProfileImageOptimized string  `json:"profileImageOptimized"`
}

type ClosedPosition struct {
	ProxyWallet string  `json:"proxyWallet"`
	Asset       string  `json:"asset"`
	ConditionID string  `json:"conditionId"`
	AvgPrice    float64 `json:"avgPrice"`
	TotalBought float64 `json:"totalBought"`
	RealizedPnL float64 `json:"realizedPnl"`
	CurPrice    float64 `json:"curPrice"`
	Title       string  `json:"title"`
	Slug        string  `json:"slug"`
	Outcome     string  `json:"outcome"`
}

type Candidate struct {
	Address          string             `json:"address"`
	Sources          map[string]int     `json:"sources"`
	Names            map[string]int     `json:"names,omitempty"`
	ObservedTrades   int                `json:"observed_trades"`
	ObservedNotional float64            `json:"observed_notional"`
	ObservedHolders  int                `json:"observed_holders"`
	MaxHolderShares  float64            `json:"max_holder_shares"`
	Markets          map[string]float64 `json:"markets,omitempty"`
}

type WalletStats struct {
	RawTrades                  int     `json:"raw_trades"`
	ValidTrades                int     `json:"valid_trades"`
	BuyTrades                  int     `json:"buy_trades"`
	SellTrades                 int     `json:"sell_trades"`
	LargeTrades                int     `json:"large_trades"`
	BuyRatio                   float64 `json:"buy_ratio"`
	ExtremePriceRatio          float64 `json:"extreme_price_ratio"`
	FixedAmountRatio           float64 `json:"fixed_amount_ratio"`
	FixedPriceRatio            float64 `json:"fixed_price_ratio"`
	BurstTrades                int     `json:"burst_trades"`
	MaxTradesPerMinute         int     `json:"max_trades_per_minute"`
	DistinctMarkets            int     `json:"distinct_markets"`
	DistinctCategories         int     `json:"distinct_categories"`
	TopCategory                string  `json:"top_category,omitempty"`
	TopCategoryRatio           float64 `json:"top_category_ratio"`
	TargetTrades               int     `json:"target_trades"`
	TargetLargeTrades          int     `json:"target_large_trades"`
	TargetTradeRatio           float64 `json:"target_trade_ratio"`
	TopMarketRatio             float64 `json:"top_market_ratio"`
	OppositeSideMarkets        int     `json:"opposite_side_markets"`
	ClosedPositions            int     `json:"closed_positions"`
	ClosedPnL                  float64 `json:"closed_pnl"`
	ClosedCapital              float64 `json:"closed_capital"`
	ClosedROI                  float64 `json:"closed_roi"`
	PositiveClosed             int     `json:"positive_closed"`
	ClosedWinRate              float64 `json:"closed_win_rate"`
	AvgTradeNotional           float64 `json:"avg_trade_notional"`
	CopyBuys                   int     `json:"copy_buys"`
	CopyClosedTrades           int     `json:"copy_closed_trades"`
	CopyWins                   int     `json:"copy_wins"`
	CopyPnL                    float64 `json:"copy_pnl"`
	CopyCapital                float64 `json:"copy_capital"`
	CopyROI                    float64 `json:"copy_roi"`
	CopyWinRate                float64 `json:"copy_win_rate"`
	CopyOpenPositions          int     `json:"copy_open_positions"`
	CopyOpenCost               float64 `json:"copy_open_cost"`
	TargetCopyBuys             int     `json:"target_copy_buys"`
	TargetCopyClosed           int     `json:"target_copy_closed_trades"`
	TargetCopyWins             int     `json:"target_copy_wins"`
	TargetCopyPnL              float64 `json:"target_copy_pnl"`
	TargetCopyCapital          float64 `json:"target_copy_capital"`
	TargetCopyROI              float64 `json:"target_copy_roi"`
	TargetCopyWinRate          float64 `json:"target_copy_win_rate"`
	TargetCopyOpen             int     `json:"target_copy_open_positions"`
	TargetCopyOpenCost         float64 `json:"target_copy_open_cost"`
	FootballScoreTrades        int     `json:"football_score_trades"`
	FootballScoreLargeTrades   int     `json:"football_score_large_trades"`
	FootballScoreClosed        int     `json:"football_score_closed_positions"`
	FootballScoreClosedPnL     float64 `json:"football_score_closed_pnl"`
	FootballScoreClosedCapital float64 `json:"football_score_closed_capital"`
	FootballScoreClosedROI     float64 `json:"football_score_closed_roi"`
	FootballScoreClosedWins    int     `json:"football_score_closed_wins"`
	FootballScoreCopyBuys      int     `json:"football_score_copy_buys"`
	FootballScoreCopyClosed    int     `json:"football_score_copy_closed_trades"`
	FootballScoreCopyWins      int     `json:"football_score_copy_wins"`
	FootballScoreCopyPnL       float64 `json:"football_score_copy_pnl"`
	FootballScoreCopyCapital   float64 `json:"football_score_copy_capital"`
	FootballScoreCopyROI       float64 `json:"football_score_copy_roi"`
	FootballScoreCopyWinRate   float64 `json:"football_score_copy_win_rate"`
	FootballScoreCopyOpen      int     `json:"football_score_copy_open_positions"`
	FootballScoreCopyOpenCost  float64 `json:"football_score_copy_open_cost"`
}

type WalletScore struct {
	Address         string         `json:"address"`
	Tier            string         `json:"tier"`
	Reason          string         `json:"reason"`
	FollowAction    string         `json:"follow_action"`
	BotScore        float64        `json:"bot_score"`
	Edge            float64        `json:"edge_score"`
	SmartMoneyScore float64        `json:"smart_money_score"`
	Score           float64        `json:"score"`
	RiskFlags       []string       `json:"risk_flags,omitempty"`
	Strengths       []string       `json:"strengths,omitempty"`
	Sources         map[string]int `json:"sources,omitempty"`
	Stats           WalletStats    `json:"stats"`
	DataStatus      string         `json:"data_status,omitempty"`
	DataIssues      []string       `json:"data_issues,omitempty"`
}

type HTTPStats struct {
	Retries    int64 `json:"retries"`
	RateLimits int64 `json:"rate_limits"`
	Failures   int64 `json:"failures"`
}

type DataQuality struct {
	Complete          int `json:"complete"`
	CachedActivity    int `json:"cached_activity"`
	PreservedPrevious int `json:"preserved_previous"`
	Incomplete        int `json:"incomplete"`
}

type Result struct {
	Markets     []Market      `json:"markets"`
	Candidates  []*Candidate  `json:"candidates"`
	Scores      []WalletScore `json:"scores"`
	HTTP        HTTPStats     `json:"http"`
	DataQuality DataQuality   `json:"data_quality"`
}
