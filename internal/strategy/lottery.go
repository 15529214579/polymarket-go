package strategy

import (
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

// SportFamily is a coarse classification used by lottery-mode filtering.
// LoL outcomes are usually well-priced once the game is a few minutes in
// (static-advantage metagame), so we skip very low prices in LoL to avoid
// buying near-certain losses. NBA/EPL upsets happen far more often, so the
// global floor applies unchanged.
type SportFamily string

const (
	SportLoL        SportFamily = "lol"
	SportDota2      SportFamily = "dota2"
	SportBasketball SportFamily = "basketball"
	SportFootball   SportFamily = "football"
	SportUnknown    SportFamily = "unknown"
)

// ClassifySport maps a feed.Market to its SportFamily. Falls back to
// SportUnknown for anything outside FilterSports.
func ClassifySport(m feed.Market) SportFamily {
	switch {
	case feed.IsLoLMarket(m):
		return SportLoL
	case feed.IsDota2Market(m):
		return SportDota2
	case feed.IsBasketballMarket(m):
		return SportBasketball
	case feed.IsFootballMarket(m):
		return SportFootball
	default:
		return SportUnknown
	}
}

// LotteryConfig defines the lottery strategy band + per-sport floor overrides
// and scan cadence. See SPEC §2.5.
type LotteryConfig struct {
	MinPrice      float64       // global lottery floor (default 0.05)
	MaxPrice      float64       // global lottery ceiling (default 0.30)
	LoLMinPrice   float64       // LoL-specific floor; overrides MinPrice if higher (default 0.15)
	Dota2MinPrice float64       // Dota 2-specific floor; overrides MinPrice if higher (default 0.15)
	SizeUSD       float64       // per-entry paper size (default 1.0)
	ScanInterval  time.Duration // scanner cadence (default 5m)
}

// DefaultLotteryConfig returns the SPEC §2.5 defaults.
func DefaultLotteryConfig() LotteryConfig {
	return LotteryConfig{
		MinPrice:      0.05,
		MaxPrice:      0.30,
		LoLMinPrice:   0.15,
		Dota2MinPrice: 0.15,
		SizeUSD:       1.0,
		ScanInterval:  5 * time.Minute,
	}
}

// LotteryCandidate is one asset eligible to be paper-opened as a lottery entry.
type LotteryCandidate struct {
	AssetID string
	Market  string
	Mid     float64
	Sport   SportFamily
	Time    time.Time
}

// LotterySampler is the subset of *feed.Sampler that lottery needs.
// Kept narrow so tests don't need a full sampler.
type LotterySampler interface {
	TickTail(assetID string, n int) ([]feed.Tick, bool)
}

// EffectiveFloor returns the actual minimum price to allow for a given sport
// under cfg. LoL tightens to LoLMinPrice when it's higher than the global floor.
func EffectiveFloor(cfg LotteryConfig, sport SportFamily) float64 {
	floor := cfg.MinPrice
	switch sport {
	case SportLoL:
		if cfg.LoLMinPrice > floor {
			floor = cfg.LoLMinPrice
		}
	case SportDota2:
		if cfg.Dota2MinPrice > floor {
			floor = cfg.Dota2MinPrice
		}
	}
	return floor
}

// IsEligible reports whether (mid, sport) satisfies the lottery band rules.
// Callers are responsible for position de-duplication (per-asset / per-market),
// which is handled outside this package via PositionManager.
func IsEligible(cfg LotteryConfig, sport SportFamily, mid float64) bool {
	if mid <= 0 {
		return false
	}
	floor := EffectiveFloor(cfg, sport)
	return mid >= floor && mid <= cfg.MaxPrice
}

// ScanEligible walks assetSport, reads the latest mid per asset via sampler,
// and returns candidates that pass IsEligible. Order is undefined.
func ScanEligible(sampler LotterySampler, assetSport map[string]SportFamily, cfg LotteryConfig) []LotteryCandidate {
	if sampler == nil || len(assetSport) == 0 {
		return nil
	}
	out := make([]LotteryCandidate, 0, 8)
	for id, sport := range assetSport {
		ticks, ok := sampler.TickTail(id, 1)
		if !ok || len(ticks) == 0 {
			continue
		}
		t := ticks[len(ticks)-1]
		if !IsEligible(cfg, sport, t.Mid) {
			continue
		}
		out = append(out, LotteryCandidate{
			AssetID: id,
			Market:  t.Market,
			Mid:     t.Mid,
			Sport:   sport,
			Time:    t.Time,
		})
	}
	return out
}
