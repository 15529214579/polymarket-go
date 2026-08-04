package order

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestReadOnlyOnChainCannotSignTransaction(t *testing.T) {
	client := &OnChain{address: common.HexToAddress("0x1234")}
	err := client.sendTx(context.Background(), common.HexToAddress(PUSDAddress), nil, 21_000)
	if err == nil || !strings.Contains(err.Error(), "read-only") {
		t.Fatalf("expected read-only signing rejection, got %v", err)
	}
}

func TestRedemptionIndexSetRejectsNonBinaryOutcome(t *testing.T) {
	if _, err := redemptionIndexSet(-1); err == nil {
		t.Fatal("negative outcome index accepted")
	}
	if _, err := redemptionIndexSet(2); err == nil {
		t.Fatal("third outcome accepted for binary redemption")
	}
	for index, want := range []int64{1, 2} {
		got, err := redemptionIndexSet(index)
		if err != nil || got.Int64() != want {
			t.Fatalf("index=%d got=%v err=%v want=%d", index, got, err, want)
		}
	}
}

func TestNegRiskFallbackRequiresConfirmedRevert(t *testing.T) {
	if negRiskFallbackAllowed(errors.New("receipt timeout")) {
		t.Fatal("receipt timeout must not trigger a duplicate redemption")
	}
	reverted := &txRevertedError{hash: common.HexToHash("0x1")}
	if !negRiskFallbackAllowed(fmt.Errorf("adapter failed: %w", reverted)) {
		t.Fatal("confirmed revert should allow direct CTF fallback")
	}
}
