package walletdiscover

import (
	"regexp"
	"strings"
	"time"

	"github.com/15529214579/polymarket-go/internal/feed"
)

var addrRE = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
var nationalFootballWinRE = regexp.MustCompile(`\bwill ([a-z][a-z .'-]+) win on \d{4}-\d{2}-\d{2}\??\b`)
var soccerTokenRE = regexp.MustCompile(`\b(?:epl|ucl|uefa|fifa|fifwc|mls)\b`)

var footballNationNames = map[string]struct{}{
	"algeria": {}, "argentina": {}, "australia": {}, "austria": {},
	"belgium": {}, "brazil": {}, "canada": {}, "chile": {}, "china": {},
	"colombia": {}, "croatia": {}, "denmark": {}, "ecuador": {},
	"egypt": {}, "england": {}, "france": {}, "germany": {}, "ghana": {},
	"greece": {}, "hungary": {}, "iran": {}, "ireland": {}, "italy": {},
	"japan": {}, "mexico": {}, "morocco": {}, "netherlands": {},
	"new zealand": {}, "nigeria": {}, "norway": {}, "paraguay": {},
	"peru": {}, "poland": {}, "portugal": {}, "qatar": {}, "romania": {},
	"saudi arabia": {}, "scotland": {}, "senegal": {}, "serbia": {},
	"south africa": {}, "south korea": {}, "spain": {}, "sweden": {},
	"switzerland": {}, "tunisia": {}, "turkey": {}, "ukraine": {},
	"united states": {}, "uruguay": {}, "usa": {}, "wales": {},
}

func normalizeAddress(addr string) string {
	addr = strings.ToLower(strings.TrimSpace(addr))
	if !addrRE.MatchString(addr) {
		return ""
	}
	return addr
}

func IsNoisyMarket(m Market) bool {
	text := strings.ToLower(m.Question + " " + m.Slug + " " + m.Category)
	noisy := []string{
		"bitcoin up or down",
		"btc updown",
		"up or down -",
		"5m-",
		"15m-",
		"price of bitcoin",
		"solana up or down",
		"ethereum up or down",
		"chelsea clinton",
		"presidential",
		"presidential nomination",
		" election",
		" nomination",
	}
	for _, needle := range noisy {
		if strings.Contains(text, needle) {
			return true
		}
	}
	return false
}

func GoodDiscoveryMarket(m Market, now time.Time) bool {
	if m.ConditionID == "" || m.Closed || !m.Active {
		return false
	}
	if IsNoisyMarket(m) {
		return false
	}
	if m.LiquidityUSD() < 500 {
		return false
	}
	end := m.ParsedEndDate()
	if !end.IsZero() && end.Before(now.Add(-24*time.Hour)) {
		return false
	}
	return true
}

func MarketTargetCategory(m Market) string {
	return targetCategory(m.Question + " " + m.Slug + " " + m.Category)
}

func TradeTargetCategory(t Trade) string {
	return targetCategory(t.Title + " " + t.Slug + " " + t.EventSlug)
}

func TargetCategoryAllowed(category, rawAllowed string) bool {
	if rawAllowed == "" {
		return true
	}
	allowed := map[string]struct{}{}
	for _, part := range strings.Split(rawAllowed, ",") {
		part = strings.ToLower(strings.TrimSpace(part))
		switch part {
		case "football":
			part = "soccer"
		case "esport":
			part = "esports"
		}
		if part != "" {
			allowed[part] = struct{}{}
		}
	}
	_, ok := allowed[category]
	return ok
}

func targetCategory(text string) string {
	text = strings.ToLower(text)
	if isNonSportsTargetText(text) {
		return "other"
	}
	if isDerivativeTargetText(text) || feed.IsOutrightFollowMarketText(text) {
		return "other"
	}
	if strings.Contains(text, "tennis") ||
		strings.Contains(text, "wimbledon") ||
		strings.Contains(text, " atp") ||
		strings.Contains(text, " wta") ||
		strings.HasPrefix(text, "atp-") ||
		strings.HasPrefix(text, "wta-") {
		return "other"
	}
	if isNationalFootballWinText(text) {
		return "soccer"
	}
	switch {
	case containsAny(text, []string{
		"nba", "wnba", "basketball",
		"knicks", "celtics", "lakers", "warriors", "thunder", "pacers", "spurs", "mavericks",
		"timberwolves", "nuggets", "suns", "bucks", "sixers", "76ers", "clippers", "heat",
	}):
		return "basketball"
	case containsAny(text, []string{
		"premier league", "la liga", "bundesliga", "serie a", "ligue 1",
		"champions league", "fifa world cup",
		"copa ", "concacaf", "conmebol", "eredivisie", "liga mx",
		"soccer", "fútbol", "futbol", "club world cup",
		" arsenal", "chelsea", "man city", "manchester city", "manchester united", "liverpool",
		"barcelona", "real madrid", "psg", "inter milan", "bayern", "dortmund", " fc ",
	}) || soccerTokenRE.MatchString(text):
		return "soccer"
	case containsAny(text, []string{
		"lol", "league of legends", "lck", "lpl", "msi", "worlds",
		"dota", "dota2", "cs2", "csgo", "valorant", "esport",
	}):
		return "esports"
	default:
		return "other"
	}
}

func isNonSportsTargetText(text string) bool {
	return strings.Contains(text, "chelsea clinton") ||
		strings.Contains(text, "presidential") ||
		strings.Contains(text, " election") ||
		strings.Contains(text, " nomination")
}

func isDerivativeTargetText(text string) bool {
	if strings.Contains(text, " and ") || strings.Contains(text, " & ") {
		return true
	}
	for _, bad := range []string{
		"spread:",
		"game handicap",
		"games total",
		"total:",
		"over/under",
		"o/u",
		"correct score",
		"exact score",
		"first to score",
		"both teams to score",
		"to win 2-0",
		"match result",
		"goals during",
		" score ",
	} {
		if strings.Contains(text, bad) {
			return true
		}
	}
	return false
}

func isNationalFootballWinText(text string) bool {
	for _, match := range nationalFootballWinRE.FindAllStringSubmatch(text, -1) {
		if len(match) != 2 {
			continue
		}
		if _, ok := footballNationNames[strings.TrimSpace(match[1])]; ok {
			return true
		}
	}
	return false
}

func containsAny(text string, needles []string) bool {
	for _, n := range needles {
		if strings.Contains(text, n) {
			return true
		}
	}
	return false
}

func QualifyingTrade(t Trade, cfg Config, allowedMarkets map[string]struct{}) bool {
	if normalizeAddress(t.ProxyWallet) == "" {
		return false
	}
	if t.Type != "" && !strings.EqualFold(t.Type, "TRADE") {
		return false
	}
	side := strings.ToUpper(t.Side)
	if side != "BUY" && side != "SELL" {
		return false
	}
	if t.Price < cfg.MinPrice || t.Price > cfg.MaxPrice {
		return false
	}
	if t.NotionalUSD() < cfg.MinNotionalUSD {
		return false
	}
	if len(allowedMarkets) > 0 {
		if _, ok := allowedMarkets[t.ConditionID]; !ok {
			return false
		}
	}
	if !TargetCategoryAllowed(TradeTargetCategory(t), cfg.TargetCategories) {
		return false
	}
	if followTargetCategoriesOnly(cfg.TargetCategories) && !feed.IsFollowTargetMarket(feed.Market{Question: t.Title, Slug: firstNonEmpty(t.Slug, t.EventSlug)}) {
		return false
	}
	if strings.Contains(strings.ToLower(t.Title+" "+t.Slug), "bitcoin up or down") {
		return false
	}
	return true
}

func followTargetCategoriesOnly(rawAllowed string) bool {
	allowed := map[string]struct{}{}
	for _, part := range strings.Split(rawAllowed, ",") {
		part = strings.ToLower(strings.TrimSpace(part))
		if part == "" {
			continue
		}
		if part == "football" {
			part = "soccer"
		}
		if part == "esport" {
			part = "esports"
		}
		allowed[part] = struct{}{}
	}
	if len(allowed) == 0 {
		return false
	}
	for _, part := range []string{"basketball", "soccer", "esports"} {
		if _, ok := allowed[part]; ok {
			delete(allowed, part)
		}
	}
	return len(allowed) == 0
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
