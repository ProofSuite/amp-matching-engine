package testutils

import "github.com/ethereum/go-ethereum/common"

func GetTestAddress1() common.Address {
	return common.HexToAddress("0x1")
}

func GetTestAddress2() common.Address {
	return common.HexToAddress("0x2")
}

func GetTestAddress3() common.Address {
	return common.HexToAddress("0x3")
}
