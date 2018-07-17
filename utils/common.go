package utils

import (
	"fmt"
	"strings"
)

// UintToPaddedString converts an int to string of length 19 by padding with 0
func UintToPaddedString(num int64) string {
	return fmt.Sprintf("%019d", num)
}

// GetTickChannelID is used to get the channel id for OHLCV data streaming
// it takes pairname, duration and units of data streaming
func GetTickChannelID(pairName, unit string, duration int64) string {
	pairName = strings.ToLower(pairName)
	return fmt.Sprintf("%s::%d::%s", pairName, duration, unit)
}
