package order

import (
	"context"
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
