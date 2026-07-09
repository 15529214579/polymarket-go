package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/15529214579/polymarket-go/internal/walletdiscover"
)

func TestDecide_ReviewsNoisyNonCoreWalletWithoutQuarantine(t *testing.T) {
	perfs := []*walletPerf{{
		Meta:    walletMeta{Address: "0x1111111111111111111111111111111111111111", List: "target", Tier: "C"},
		Pending: 2,
		EventCD: 1,
	}}

	decide(perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if perfs[0].Decision != "review_noise" {
		t.Fatalf("Decision=%q, want review_noise", perfs[0].Decision)
	}
}

func TestDecide_PositiveEdgeWalletNotNoiseReviewed(t *testing.T) {
	perfs := []*walletPerf{{
		Meta:        walletMeta{Address: "0xe8725f4bd4de73c888dda35c4850f702c322819a", List: "watch", Tier: "B", Bot: 33.4},
		Pending:     2,
		AssetCD:     1,
		EventCD:     1,
		EdgeSamples: 5,
		EdgeWins:    5,
		EdgeAvgPP:   17.09,
		Score: &walletdiscover.WalletScore{
			BotScore: 33.4,
			Stats:    walletdiscover.WalletStats{CopyROI: 92.3, CopyClosedTrades: 29},
		},
	}}

	decide(perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if perfs[0].Decision == "review_noise" {
		t.Fatalf("Decision=%q, want positive-edge wallet not noise reviewed", perfs[0].Decision)
	}
	if perfs[0].Decision != "learning" {
		t.Fatalf("Decision=%q, want learning until closed live signals accrue", perfs[0].Decision)
	}
}

func TestDecide_EventQuarantineTakesPrecedenceOverNoiseReview(t *testing.T) {
	perfs := []*walletPerf{{
		Meta:               walletMeta{Address: "0x2222222222222222222222222222222222222222", List: "watch", Tier: "B"},
		EventSignals:       1,
		EventROI:           -35,
		ProvenEventSignals: 1,
		ProvenEventROI:     -35,
		Pending:            3,
	}}

	decide(perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if perfs[0].Decision != "quarantine" {
		t.Fatalf("Decision=%q, want quarantine", perfs[0].Decision)
	}
}

func TestDecide_DoesNotQuarantineOnMarkedOnlyEventLoss(t *testing.T) {
	perfs := []*walletPerf{{
		Meta:         walletMeta{Address: "0x4444444444444444444444444444444444444444", List: "watch", Tier: "B"},
		EventSignals: 1,
		EventMarked:  1,
		EventROI:     -35,
	}}

	decide(perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if perfs[0].Decision != "learning" {
		t.Fatalf("Decision=%q, want learning", perfs[0].Decision)
	}
}

func TestDecide_CoreWalletNotNoiseReviewed(t *testing.T) {
	perfs := []*walletPerf{{
		Meta:    walletMeta{Address: "0x3333333333333333333333333333333333333333", List: "core", Tier: "A"},
		Pending: 5,
	}}

	decide(perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if perfs[0].Decision != "learning" {
		t.Fatalf("Decision=%q, want learning", perfs[0].Decision)
	}
}

func TestApplyPerformanceSnapshot_LoadsProvenByWallet(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "perf.json")
	raw := `{
		"event_capped_by_wallet": [
			{"wallet":"0x5555555555555555555555555555555555555555","signals":2,"marked":2,"pnl_usd":-7,"stake_usd":20,"return_pct":-35}
		],
		"event_capped_proven_by_wallet": [
			{"wallet":"0x5555555555555555555555555555555555555555","signals":1,"closed":1,"wins":1,"pnl_usd":4,"stake_usd":10,"return_pct":40}
		]
	}`
	if err := os.WriteFile(path, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	perfs := []*walletPerf{{
		Meta: walletMeta{Address: "0x5555555555555555555555555555555555555555", List: "watch"},
	}}

	if err := applyPerformanceSnapshot(perfs, path); err != nil {
		t.Fatal(err)
	}
	if perfs[0].EventSignals != 2 || perfs[0].EventROI != -35 {
		t.Fatalf("event stats=%d/%.1f, want 2/-35", perfs[0].EventSignals, perfs[0].EventROI)
	}
	if perfs[0].ProvenEventSignals != 1 || perfs[0].ProvenEventROI != 40 {
		t.Fatalf("proven event stats=%d/%.1f, want 1/40", perfs[0].ProvenEventSignals, perfs[0].ProvenEventROI)
	}
}

func TestDecide_PromotesTapeObservationOnStrongEdge(t *testing.T) {
	perfs := []*walletPerf{{
		Meta:           walletMeta{Address: "0x7777777777777777777777777777777777777777", List: "tape_observe", Tier: "TAPE"},
		EdgeSamples:    3,
		EdgeWins:       2,
		EdgeDeltaPPSum: 4.5,
		EdgeAvgPP:      1.5,
	}}

	decide(perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if perfs[0].Decision != "promote_tape_candidate" {
		t.Fatalf("Decision=%q, want promote_tape_candidate", perfs[0].Decision)
	}
}

func TestDecide_KeepsExistingTapeCandidateOnStrongEdge(t *testing.T) {
	perfs := []*walletPerf{{
		Meta:        walletMeta{Address: "0x7777777777777777777777777777777777777777", List: "tape_candidate", Tier: "D"},
		EdgeSamples: 3,
		EdgeWins:    2,
		EdgeAvgPP:   1.5,
	}}

	decide(perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if perfs[0].Decision != "promote_tape_candidate" {
		t.Fatalf("Decision=%q, want promote_tape_candidate", perfs[0].Decision)
	}
}

func TestDecide_DoesNotPromoteTapeCandidateWithHighBot(t *testing.T) {
	perfs := []*walletPerf{{
		Meta:        walletMeta{Address: "0xc0f5345be7d24194892f10627123ffe91c633d04", List: "tape_observe", Tier: "D", Bot: 51.7},
		EdgeSamples: 3,
		EdgeWins:    3,
		EdgeAvgPP:   24.77,
	}}

	decide(perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if perfs[0].Decision == "promote_tape_candidate" {
		t.Fatalf("Decision=%q, want high bot wallet blocked", perfs[0].Decision)
	}
	if !strings.Contains(perfs[0].Reason, "bot 51.7") {
		t.Fatalf("Reason=%q, want bot rejection", perfs[0].Reason)
	}
}

func TestDecide_PromotesTapeFollowOnlyOnStrictEdge(t *testing.T) {
	perfs := []*walletPerf{{
		Meta:             walletMeta{Address: "0x7777777777777777777777777777777777777777", List: "tape_candidate", Tier: "B", Bot: 32},
		EdgeSamples:      6,
		EdgeWins:         5,
		EdgeAvgPP:        1.8,
		Edge5mSamples:    3,
		Edge5mAvgPP:      0.8,
		Edge15mSamples:   2,
		Edge15mAvgPP:     0.2,
		EdgeDeltaPPSum:   10.8,
		Edge5mDeltaPPSum: 2.4,
	}}

	decide(perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if perfs[0].Decision != "promote_tape_follow" {
		t.Fatalf("Decision=%q, want promote_tape_follow", perfs[0].Decision)
	}
}

func TestDecide_DoesNotFollowTapeCandidateWithNegative15m(t *testing.T) {
	perfs := []*walletPerf{{
		Meta:           walletMeta{Address: "0x7777777777777777777777777777777777777777", List: "tape_candidate", Tier: "B", Bot: 32},
		EdgeSamples:    6,
		EdgeWins:       5,
		EdgeAvgPP:      1.8,
		Edge5mSamples:  3,
		Edge5mAvgPP:    0.8,
		Edge15mSamples: 2,
		Edge15mAvgPP:   -1.5,
	}}

	decide(perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if perfs[0].Decision != "review_tape_reversal" {
		t.Fatalf("Decision=%q, want review_tape_reversal", perfs[0].Decision)
	}
}

func TestDecide_DoesNotReverseTapeCandidateOnSingle15mLoss(t *testing.T) {
	perfs := []*walletPerf{{
		Meta:           walletMeta{Address: "0x7777777777777777777777777777777777777777", List: "tape_candidate", Tier: "B", Bot: 32},
		EdgeSamples:    4,
		EdgeWins:       3,
		EdgeAvgPP:      1.8,
		Edge15mSamples: 1,
		Edge15mAvgPP:   -2.0,
	}}

	decide(perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if perfs[0].Decision != "promote_tape_candidate" {
		t.Fatalf("Decision=%q, want promote_tape_candidate", perfs[0].Decision)
	}
}

func TestDecide_ReversesTapeObserveOnSevereSingleEdgeDrawdown(t *testing.T) {
	perfs := []*walletPerf{{
		Meta:        walletMeta{Address: "0xf3ce7f04dde4f8c5aed19a20e5ecd8520e5ca57a", List: "tape_observe", Tier: "B", Bot: 24.5},
		EdgeSamples: 1,
		EdgeWins:    0,
		EdgeAvgPP:   -47.28,
	}}

	decide(perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if perfs[0].Decision != "review_tape_reversal" {
		t.Fatalf("Decision=%q, want review_tape_reversal", perfs[0].Decision)
	}
	if !strings.Contains(perfs[0].Reason, "severe edge drawdown") {
		t.Fatalf("Reason=%q, want severe edge drawdown", perfs[0].Reason)
	}
}

func TestDecide_ReversesTapeObserveOnPersistentNegativeEdge(t *testing.T) {
	perfs := []*walletPerf{{
		Meta:        walletMeta{Address: "0x84cf65356d3ac7d87e26a94119e1e2c8dd802f63", List: "tape_observe", Tier: "BOT", Bot: 62},
		EdgeSamples: 8,
		EdgeWins:    1,
		EdgeAvgPP:   -0.42,
	}}

	decide(perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if perfs[0].Decision != "review_tape_reversal" {
		t.Fatalf("Decision=%q, want review_tape_reversal", perfs[0].Decision)
	}
	if !strings.Contains(perfs[0].Reason, "persistent negative edge") {
		t.Fatalf("Reason=%q, want persistent negative edge", perfs[0].Reason)
	}
}

func TestDecide_DoesNotReverseTapeObserveOnSmallNegativeSample(t *testing.T) {
	perfs := []*walletPerf{{
		Meta:        walletMeta{Address: "0x9999999999999999999999999999999999999999", List: "tape_observe", Tier: "D", Bot: 30},
		EdgeSamples: 3,
		EdgeWins:    0,
		EdgeAvgPP:   -0.50,
	}}

	decide(perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if perfs[0].Decision != "learning" {
		t.Fatalf("Decision=%q, want learning until persistent negative sample threshold", perfs[0].Decision)
	}
}

func TestDecide_RetainsTapeReversalLayer(t *testing.T) {
	perfs := []*walletPerf{{
		Meta:        walletMeta{Address: "0xf3ce7f04dde4f8c5aed19a20e5ecd8520e5ca57a", List: "tape_reversal", Tier: "B", Bot: 24.5},
		EdgeSamples: 1,
		EdgeWins:    0,
		EdgeAvgPP:   -47.28,
	}}

	decide(perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if perfs[0].Decision != "review_tape_reversal" {
		t.Fatalf("Decision=%q, want review_tape_reversal", perfs[0].Decision)
	}
	if !strings.Contains(perfs[0].Reason, "retained reversal") {
		t.Fatalf("Reason=%q, want retained reversal", perfs[0].Reason)
	}
}

func TestLoadWalletMetas_CommaSeparatedFiles(t *testing.T) {
	dir := t.TempDir()
	push := filepath.Join(dir, "push.txt")
	observe := filepath.Join(dir, "observe.txt")
	if err := os.WriteFile(push, []byte("0x1111111111111111111111111111111111111111 # list=core tier=A\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(observe, []byte("0x2222222222222222222222222222222222222222 # list=tape_observe tier=TAPE\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := loadWalletMetas(push + "," + observe)
	if err != nil {
		t.Fatal(err)
	}
	if got["0x1111111111111111111111111111111111111111"].List != "core" {
		t.Fatalf("push meta=%q, want core", got["0x1111111111111111111111111111111111111111"].List)
	}
	if got["0x2222222222222222222222222222222222222222"].List != "tape_observe" {
		t.Fatalf("observe meta=%q, want tape_observe", got["0x2222222222222222222222222222222222222222"].List)
	}
}

func TestWriteOutputs_WritesTapeCandidatesSeparately(t *testing.T) {
	dir := t.TempDir()
	report := filepath.Join(dir, "report.md")
	quarantine := filepath.Join(dir, "quarantine.txt")
	reviewNoise := filepath.Join(dir, "review-noise.txt")
	candidates := filepath.Join(dir, "tape-candidates.txt")
	perfs := []*walletPerf{{
		Meta:           walletMeta{Address: "0x7777777777777777777777777777777777777777", List: "tape_observe", Tier: "D"},
		Decision:       "promote_tape_candidate",
		Reason:         "edge 66.7% win and avg +3.49pp over 3 samples",
		EdgeSamples:    3,
		EdgeWins:       2,
		EdgeAvgPP:      3.49,
		Edge5mAvgPP:    0.26,
		Edge15mAvgPP:   0,
		SignalRank:     0,
		EdgeDeltaPPSum: 10.47,
	}}

	q, n, c, f, r, err := writeOutputs(report, quarantine, reviewNoise, candidates, filepath.Join(dir, "follow.txt"), filepath.Join(dir, "reversal.txt"), perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if err != nil {
		t.Fatal(err)
	}
	if q != 0 || n != 0 || c != 1 || f != 0 || r != 0 {
		t.Fatalf("quarantine/noise/candidates/follow/reversal=%d/%d/%d/%d/%d, want 0/0/1/0/0", q, n, c, f, r)
	}
	qb, err := os.ReadFile(quarantine)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(qb)) != "" {
		t.Fatalf("quarantine content=%q, want empty", string(qb))
	}
	cb, err := os.ReadFile(candidates)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(cb), "0x7777777777777777777777777777777777777777 # list=tape_candidate") {
		t.Fatalf("candidate file missing wallet:\n%s", string(cb))
	}
	if !strings.Contains(string(cb), "edgeAvgPP=+3.49") {
		t.Fatalf("candidate file missing edge stats:\n%s", string(cb))
	}
}

func TestWriteOutputs_WritesTapeFollowSeparately(t *testing.T) {
	dir := t.TempDir()
	report := filepath.Join(dir, "report.md")
	quarantine := filepath.Join(dir, "quarantine.txt")
	reviewNoise := filepath.Join(dir, "review-noise.txt")
	candidates := filepath.Join(dir, "tape-candidates.txt")
	follow := filepath.Join(dir, "tape-follow.txt")
	perfs := []*walletPerf{{
		Meta:         walletMeta{Address: "0x8888888888888888888888888888888888888888", List: "tape_candidate", Tier: "B", Bot: 32},
		Decision:     "promote_tape_follow",
		Reason:       "follow-ready",
		EdgeSamples:  6,
		EdgeWins:     5,
		EdgeAvgPP:    1.8,
		Edge5mAvgPP:  0.8,
		Edge15mAvgPP: 0.2,
	}}

	q, n, c, f, r, err := writeOutputs(report, quarantine, reviewNoise, candidates, follow, filepath.Join(dir, "reversal.txt"), perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if err != nil {
		t.Fatal(err)
	}
	if q != 0 || n != 0 || c != 1 || f != 1 || r != 0 {
		t.Fatalf("quarantine/noise/candidates/follow/reversal=%d/%d/%d/%d/%d, want 0/0/1/1/0", q, n, c, f, r)
	}
	fb, err := os.ReadFile(follow)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(fb), "0x8888888888888888888888888888888888888888 # list=tape_follow") {
		t.Fatalf("follow file missing wallet:\n%s", string(fb))
	}
}

func TestWriteOutputs_WritesTapeReversalSeparately(t *testing.T) {
	dir := t.TempDir()
	report := filepath.Join(dir, "report.md")
	quarantine := filepath.Join(dir, "quarantine.txt")
	reviewNoise := filepath.Join(dir, "review-noise.txt")
	candidates := filepath.Join(dir, "tape-candidates.txt")
	follow := filepath.Join(dir, "tape-follow.txt")
	reversal := filepath.Join(dir, "tape-reversal.txt")
	perfs := []*walletPerf{{
		Meta:           walletMeta{Address: "0x9999999999999999999999999999999999999999", List: "tape_candidate", Tier: "B", Bot: 28},
		Decision:       "review_tape_reversal",
		Reason:         "15m edge reversed -2.49pp over 2 samples",
		EdgeSamples:    5,
		EdgeWins:       3,
		EdgeAvgPP:      1.2,
		Edge5mAvgPP:    0.3,
		Edge15mAvgPP:   -2.49,
		SignalRank:     0,
		EdgeDeltaPPSum: 6,
	}}

	q, n, c, f, r, err := writeOutputs(report, quarantine, reviewNoise, candidates, follow, reversal, perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if err != nil {
		t.Fatal(err)
	}
	if q != 0 || n != 0 || c != 0 || f != 0 || r != 1 {
		t.Fatalf("quarantine/noise/candidates/follow/reversal=%d/%d/%d/%d/%d, want 0/0/0/0/1", q, n, c, f, r)
	}
	rb, err := os.ReadFile(reversal)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(rb), "0x9999999999999999999999999999999999999999 # list=tape_reversal") {
		t.Fatalf("reversal file missing wallet:\n%s", string(rb))
	}
	cb, err := os.ReadFile(candidates)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(cb)) != "" {
		t.Fatalf("candidate file should be empty for reversal:\n%s", string(cb))
	}
}

func TestWriteOutputs_WritesReviewNoiseSeparately(t *testing.T) {
	dir := t.TempDir()
	report := filepath.Join(dir, "report.md")
	quarantine := filepath.Join(dir, "quarantine.txt")
	reviewNoise := filepath.Join(dir, "review-noise.txt")
	candidates := filepath.Join(dir, "tape-candidates.txt")
	follow := filepath.Join(dir, "tape-follow.txt")
	reversal := filepath.Join(dir, "tape-reversal.txt")
	perfs := []*walletPerf{{
		Meta:     walletMeta{Address: "0x620d7e06ce27d16532c061eba9b46c7e1833c67f", List: "scout", Tier: "B", Bot: 23.7},
		Decision: "review_noise",
		Reason:   "suppressed BUYs 8 >= 3 with no evaluated entries",
		Pending:  1,
		AssetCD:  6,
		EventCD:  1,
		Score: &walletdiscover.WalletScore{
			Stats:    walletdiscover.WalletStats{CopyROI: -60.6},
			BotScore: 23.7,
		},
	}}

	q, n, c, f, r, err := writeOutputs(report, quarantine, reviewNoise, candidates, follow, reversal, perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if err != nil {
		t.Fatal(err)
	}
	if q != 0 || n != 1 || c != 0 || f != 0 || r != 0 {
		t.Fatalf("quarantine/noise/candidates/follow/reversal=%d/%d/%d/%d/%d, want 0/1/0/0/0", q, n, c, f, r)
	}
	qb, err := os.ReadFile(quarantine)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(qb)) != "" {
		t.Fatalf("quarantine content=%q, want empty", string(qb))
	}
	nb, err := os.ReadFile(reviewNoise)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(nb), "0x620d7e06ce27d16532c061eba9b46c7e1833c67f # list=scout") {
		t.Fatalf("review-noise file missing wallet:\n%s", string(nb))
	}
	if !strings.Contains(string(nb), "suppressed=8") {
		t.Fatalf("review-noise file missing suppressed count:\n%s", string(nb))
	}
}

func TestWriteOutputs_RetainsExistingReviewNoiseOutsideCurrentPush(t *testing.T) {
	dir := t.TempDir()
	report := filepath.Join(dir, "report.md")
	quarantine := filepath.Join(dir, "quarantine.txt")
	reviewNoise := filepath.Join(dir, "review-noise.txt")
	candidates := filepath.Join(dir, "tape-candidates.txt")
	follow := filepath.Join(dir, "tape-follow.txt")
	reversal := filepath.Join(dir, "tape-reversal.txt")
	existing := "0x620d7e06ce27d16532c061eba9b46c7e1833c67f # list=scout tier=B suppressed=8 reason=suppressed BUYs 8 >= 3 with no evaluated entries\n"
	if err := os.WriteFile(reviewNoise, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}
	perfs := []*walletPerf{{
		Meta:     walletMeta{Address: "0x1111111111111111111111111111111111111111", List: "watch", Tier: "B"},
		Decision: "learning",
		Reason:   "closed signals 0 < 5",
	}}

	q, n, c, f, r, err := writeOutputs(report, quarantine, reviewNoise, candidates, follow, reversal, perfs, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if err != nil {
		t.Fatal(err)
	}
	if q != 0 || n != 1 || c != 0 || f != 0 || r != 0 {
		t.Fatalf("quarantine/noise/candidates/follow/reversal=%d/%d/%d/%d/%d, want 0/1/0/0/0", q, n, c, f, r)
	}
	nb, err := os.ReadFile(reviewNoise)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(nb)) != strings.TrimSpace(existing) {
		t.Fatalf("review-noise content=%q, want retained existing %q", string(nb), existing)
	}
	rb, err := os.ReadFile(report)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(rb), "- retained_existing_review_noise: 1") {
		t.Fatalf("report missing retained review-noise count:\n%s", string(rb))
	}
}

func TestWriteOutputs_DropsRetainedReviewNoiseWithStrongPositiveEdge(t *testing.T) {
	dir := t.TempDir()
	report := filepath.Join(dir, "report.md")
	quarantine := filepath.Join(dir, "quarantine.txt")
	reviewNoise := filepath.Join(dir, "review-noise.txt")
	candidates := filepath.Join(dir, "tape-candidates.txt")
	follow := filepath.Join(dir, "tape-follow.txt")
	reversal := filepath.Join(dir, "tape-reversal.txt")
	existing := "0xe8725f4bd4de73c888dda35c4850f702c322819a # list=watch tier=B suppressed=5 edgeN=5 edgeAvgPP=+17.09 bot=33.4 copyROI=92.3% reason=suppressed BUYs 5 >= 3 with no evaluated entries\n"
	if err := os.WriteFile(reviewNoise, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	q, n, c, f, r, err := writeOutputs(report, quarantine, reviewNoise, candidates, follow, reversal, nil, 5, 0, 5, defaultEdgePromote(), defaultTapeFollow(), 1, -30, 3)
	if err != nil {
		t.Fatal(err)
	}
	if q != 0 || n != 0 || c != 0 || f != 0 || r != 0 {
		t.Fatalf("quarantine/noise/candidates/follow/reversal=%d/%d/%d/%d/%d, want 0/0/0/0/0", q, n, c, f, r)
	}
	nb, err := os.ReadFile(reviewNoise)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(nb)) != "" {
		t.Fatalf("strong positive-edge review-noise should be dropped, got:\n%s", string(nb))
	}
}

func TestApplyEdgeSnapshots_AggregatesWalletEdge(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "edge.jsonl")
	raw := strings.Join([]string{
		`{"wallet":"0x5555555555555555555555555555555555555555","horizon_sec":300,"delta_pp":4.5}`,
		`{"wallet":"0x5555555555555555555555555555555555555555","horizon_sec":900,"delta_pp":-1.5}`,
		`{"wallet":"0x6666666666666666666666666666666666666666","horizon_sec":300,"delta_pp":9.0}`,
	}, "\n") + "\n"
	if err := os.WriteFile(path, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	perfs := []*walletPerf{{
		Meta: walletMeta{Address: "0x5555555555555555555555555555555555555555", List: "watch"},
	}}

	if err := applyEdgeSnapshots(perfs, path); err != nil {
		t.Fatal(err)
	}
	if perfs[0].EdgeSamples != 2 || perfs[0].EdgeWins != 1 {
		t.Fatalf("edge samples/wins=%d/%d, want 2/1", perfs[0].EdgeSamples, perfs[0].EdgeWins)
	}
	if perfs[0].EdgeAvgPP != 1.5 {
		t.Fatalf("EdgeAvgPP=%.2f, want 1.5", perfs[0].EdgeAvgPP)
	}
	if perfs[0].Edge5mAvgPP != 4.5 || perfs[0].Edge15mAvgPP != -1.5 {
		t.Fatalf("Edge5m/15m=%.2f/%.2f, want 4.5/-1.5", perfs[0].Edge5mAvgPP, perfs[0].Edge15mAvgPP)
	}
}

func TestEvaluate_ClassifiesFilteredSkipReasons(t *testing.T) {
	addr := "0x5555555555555555555555555555555555555555"
	base := whaleTrade{
		TS:     "2026-07-09T01:00:00Z",
		Wallet: addr,
		Side:   "BUY",
		Price:  0.5,
		Size:   1000,
		Time:   mustParseTime(t, "2026-07-09T01:00:00Z"),
	}
	trades := []whaleTrade{
		base,
		base,
		base,
	}
	trades[0].Action = "skip"
	trades[0].Reason = "derivative_filtered"
	trades[1].Action = "skip"
	trades[1].Reason = "category_filtered"
	trades[2].Action = "skip"
	trades[2].Reason = "other"

	perfs := evaluate(map[string]walletMeta{addr: {Address: addr, List: "watch"}}, nil, trades, 10)
	if len(perfs) != 1 {
		t.Fatalf("perfs len=%d, want 1", len(perfs))
	}
	if perfs[0].DerivativeFiltered != 1 || perfs[0].CategoryFiltered != 1 || perfs[0].OtherNoise != 1 {
		t.Fatalf("filtered counts deriv/category/other=%d/%d/%d, want 1/1/1", perfs[0].DerivativeFiltered, perfs[0].CategoryFiltered, perfs[0].OtherNoise)
	}
}

func mustParseTime(t *testing.T, raw string) time.Time {
	t.Helper()
	ts, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		t.Fatal(err)
	}
	return ts
}

func defaultEdgePromote() edgePromoteParams {
	return edgePromoteParams{MinSamples: 3, MinAvgPP: 1, MinWin: 60, MaxBot: 45, ReversalMin15mSamples: 2, ReversalMax15mAvgPP: -1, SevereMinSamples: 1, SevereMaxAvgPP: -20, NegativeMinSamples: 5, NegativeMaxAvgPP: -0.25, NegativeMaxWin: 20}
}

func defaultTapeFollow() tapeFollowParams {
	return tapeFollowParams{MinSamples: 6, MinAvgPP: 1.5, MinWin: 65, Min5mAvgPP: 0.5, Min15mAvgPP: 0, MaxBot: 45}
}
