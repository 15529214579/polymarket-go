package odds

import (
	"fmt"
	"log/slog"
	"math"
	"sort"
)

// ApplyConsensusFilter collapses per-bookmaker quotes into consensus rows.
// Groups by (event_id, team_or_side), drops outliers deviating >deviationPP
// from group median, requires minBooks surviving bookmakers.
// Returns synthetic rows with bookmaker="consensus(N)".
func ApplyConsensusFilter(odds []BookmakerOdds, minBooks int, deviationPP float64) []BookmakerOdds {
	if minBooks <= 0 {
		minBooks = 3
	}
	if deviationPP <= 0 {
		deviationPP = 0.05
	}

	type groupKey struct {
		eventID    string
		teamOrSide string
	}

	groups := map[groupKey][]BookmakerOdds{}
	for _, o := range odds {
		eid := o.EventID
		if eid == "" {
			eid = o.EventName
		}
		k := groupKey{eid, o.TeamOrSide}
		groups[k] = append(groups[k], o)
	}

	var result []BookmakerOdds
	droppedLow := 0
	droppedDeviants := 0

	for _, items := range groups {
		if len(items) == 0 {
			continue
		}

		probs := make([]float64, len(items))
		for i, o := range items {
			probs[i] = o.BookmakerProb
		}
		groupMedian := median(probs)

		var survivors []BookmakerOdds
		for _, o := range items {
			if math.Abs(o.BookmakerProb-groupMedian) < deviationPP {
				survivors = append(survivors, o)
			} else {
				droppedDeviants++
			}
		}

		if len(survivors) < minBooks {
			droppedLow++
			continue
		}

		survProbs := make([]float64, len(survivors))
		for i, o := range survivors {
			survProbs[i] = o.BookmakerProb
		}
		survMedian := median(survProbs)

		// Pick the bookmaker closest to the survivor median as anchor.
		anchor := survivors[0]
		bestDist := math.Abs(anchor.BookmakerProb - survMedian)
		for _, o := range survivors[1:] {
			d := math.Abs(o.BookmakerProb - survMedian)
			if d < bestDist {
				bestDist = d
				anchor = o
			}
		}

		result = append(result, BookmakerOdds{
			Sport:             anchor.Sport,
			EventID:           anchor.EventID,
			EventName:         anchor.EventName,
			TeamOrSide:        anchor.TeamOrSide,
			BookmakerProb:     math.Round(survMedian*10000) / 10000,
			Bookmaker:         fmt.Sprintf("consensus(%d)", len(survivors)),
			MarketName:        anchor.MarketName,
			EventCommenceTime: anchor.EventCommenceTime,
		})
	}

	if droppedLow > 0 || droppedDeviants > 0 {
		slog.Info("consensus_filter",
			"groups", len(groups),
			"consensus_rows", len(result),
			"dropped_low_count", droppedLow,
			"dropped_deviants", droppedDeviants,
			"min_books", minBooks,
		)
	}
	return result
}

func median(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sorted := make([]float64, len(vals))
	copy(sorted, vals)
	sort.Float64s(sorted)
	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}
