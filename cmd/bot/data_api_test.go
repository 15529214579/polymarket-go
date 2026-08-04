package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestFetchDataAPIPositionsPaginatesAndIncludesRedeemable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		redeemable := r.URL.Query().Get("redeemable") == "true"
		var page []dataAPIPosition
		switch {
		case !redeemable && offset == 0:
			page = make([]dataAPIPosition, 500)
			for i := range page {
				page[i] = dataAPIPosition{Asset: fmt.Sprintf("asset-%d", i), ConditionID: "condition"}
			}
		case !redeemable && offset == 500:
			page = []dataAPIPosition{{Asset: "asset-500", ConditionID: "condition"}}
		case redeemable && offset == 0:
			page = []dataAPIPosition{
				{Asset: "asset-1", ConditionID: "condition", Redeemable: true},
				{Asset: "redeemable", ConditionID: "condition", Redeemable: true},
			}
		}
		if err := json.NewEncoder(w).Encode(page); err != nil {
			t.Error(err)
		}
	}))
	t.Cleanup(server.Close)

	positions, err := fetchDataAPIPositionsFrom(context.Background(), server.Client(), server.URL, "0xabc")
	if err != nil {
		t.Fatal(err)
	}
	if len(positions) != 502 {
		t.Fatalf("positions=%d, want 502", len(positions))
	}
}
