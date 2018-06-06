package dex

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestComputeID(t *testing.T) {
	pair := config.TokenPairs["ZRXWETH"]
	ID := pair.ComputeID()
	expectedID := common.HexToHash("0x2a567f6079b13cc832a14db2ff45aa05827b5d1428ceb0982be7321ca78a739b")

	if ID != expectedID {
		t.Errorf("Expected computed ID to be equal to %v but got %v instead", expectedID.String(), ID.String())
	}
}
