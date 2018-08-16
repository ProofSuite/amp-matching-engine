package utils

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// UintToPaddedString converts an int to string of length 19 by padding with 0
func UintToPaddedString(num int64) string {
	return fmt.Sprintf("%019d", num)
}

// GetTickChannelID is used to get the channel id for OHLCV data streaming
// it takes pairname, duration and units of data streaming
func GetTickChannelID(bt, qt common.Address, unit string, duration int64) string {
	pair := GetPairKey(bt, qt)
	return fmt.Sprintf("%s::%d::%s", pair, duration, unit)
}

// GetPairKey return the pair key identifier corresponding to two
func GetPairKey(bt, qt common.Address) string {
	return strings.ToLower(fmt.Sprintf("%s::%s", bt.Hex(), qt.Hex()))
}
