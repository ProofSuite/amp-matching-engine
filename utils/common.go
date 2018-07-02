package utils

import (
	"fmt"
)

func UintToPaddedString(num int64) string {
	return fmt.Sprintf("%019d", num)
}
