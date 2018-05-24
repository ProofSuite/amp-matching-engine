package dex

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestComputeID(t *testing.T) {
	pair := config.TokenPairs["ZRXWETH"]
	ID := pair.ComputeID()
	expectedID := common.HexToHash("0xba6d84bb63fd28b7ca98ab0edf677fadd1354e572c830880c5ed0a2a0ff0cadb")

	if ID != expectedID {
		t.Errorf("Expected computed ID to be equal to %v but got %v instead", expectedID.String(), ID.String())
	}
}
