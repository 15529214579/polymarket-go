package walletdiscover

import "fmt"

// FootballScoreRating is independent of the generic wallet tier. Correct
// score traders commonly place baskets in bursts, so promotion is based on
// executable Yes-side history instead of relaxing the generic bot score.
type FootballScoreRating struct {
	Tier   string
	Reason string
}

func RateFootballScoreWallet(score WalletScore) FootballScoreRating {
	st := score.Stats
	if st.FootballScoreTrades < 8 || st.FootballScoreLargeTrades < 2 {
		return FootballScoreRating{Tier: "D", Reason: "insufficient score sample"}
	}
	if st.FootballScoreCopyClosed >= 3 && st.FootballScoreCopyROI <= 0 {
		return FootballScoreRating{Tier: "D", Reason: "negative score copy ROI"}
	}
	if st.FootballScoreCopyClosed < 3 && st.FootballScoreClosed >= 3 && st.FootballScoreClosedROI <= 0 {
		return FootballScoreRating{Tier: "D", Reason: "negative score closed ROI"}
	}
	copyA := st.FootballScoreCopyClosed >= 8 && st.FootballScoreCopyROI >= 10 && st.FootballScoreCopyWinRate >= 45
	closedA := st.FootballScoreClosed >= 10 && st.FootballScoreClosedROI >= 10
	if st.FootballScoreTrades >= 20 && st.FootballScoreLargeTrades >= 5 && (copyA || closedA) {
		return FootballScoreRating{Tier: "A", Reason: fmt.Sprintf("proven score edge: copy %.1f%% closed %.1f%%", st.FootballScoreCopyROI, st.FootballScoreClosedROI)}
	}
	copyB := st.FootballScoreCopyClosed >= 3 && st.FootballScoreCopyROI >= 3
	closedB := st.FootballScoreClosed >= 3 && st.FootballScoreClosedROI >= 5
	if copyB || closedB {
		return FootballScoreRating{Tier: "B", Reason: fmt.Sprintf("positive score evidence: copy %.1f%% closed %.1f%%", st.FootballScoreCopyROI, st.FootballScoreClosedROI)}
	}
	return FootballScoreRating{Tier: "C", Reason: "score holder under observation"}
}
