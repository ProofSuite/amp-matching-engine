package e2e

// func testAccount(t *testing.T, tokens []types.Token) map[*ecdsa.PrivateKey]types.Account {
// 	fmt.Printf("\n=== Starting Account test ===\n")
// 	router := NewRouter()
// 	response := make(map[*ecdsa.PrivateKey]types.Account)

// 	pk1, _ := crypto.GenerateKey()
// 	pk2, _ := crypto.GenerateKey()

// 	account1 := types.Account{
// 		Address:       crypto.PubkeyToAddress(pk1.PublicKey),
// 		IsBlocked:     false,
// 		TokenBalances: make(map[common.Address]*types.TokenBalance),
// 	}

// 	account2 := types.Account{
// 		Address:       crypto.PubkeyToAddress(pk2.PublicKey),
// 		IsBlocked:     false,
// 		TokenBalances: make(map[common.Address]*types.TokenBalance),
// 	}

// 	initBalance := big.NewInt(10000000000000000)

// 	for _, token := range tokens {
// 		account1.TokenBalances[token.ContractAddress] = &types.TokenBalance{
// 			Address:       token.ContractAddress,
// 			Symbol:        token.Symbol,
// 			Balance:       initBalance,
// 			Allowance:     initBalance,
// 			LockedBalance: big.NewInt(0),
// 		}
// 		account2.TokenBalances[token.ContractAddress] = &types.TokenBalance{
// 			Address:       token.ContractAddress,
// 			Symbol:        token.Symbol,
// 			Balance:       initBalance,
// 			Allowance:     initBalance,
// 			LockedBalance: big.NewInt(0),
// 		}
// 	}
// 	// create account test
// 	res := testAPI(router, "POST", "/account", "{\"address\":\""+account1.Address.Hex()+"\"}")
// 	assert.Equal(t, http.StatusOK, res.Code, "t1 - create account")

// 	var resp types.Account
// 	if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
// 		fmt.Printf("%v", err)
// 	}
// 	if compareAccount(t, resp, account1) {
// 		fmt.Println("PASS  't1 - create account'")
// 		account1 = resp
// 		response[pk1] = account1
// 	} else {
// 		fmt.Println("FAIL  't1 - create account'")
// 	}
// 	// Duplicate account test
// 	res1 := testAPI(router, "POST", "/account", "{\"address\":\""+account1.Address.Hex()+"\"}")

// 	if err := json.Unmarshal(res1.Body.Bytes(), &resp); err != nil {
// 		fmt.Printf("%v", err)
// 	}

// 	if assert.Equal(t, account1.ID.Hex(), resp.ID.Hex(), "t2 - create duplicate account") {
// 		fmt.Println("PASS  't2 - create duplicate account'")
// 	} else {
// 		fmt.Println("FAIL  't2 - create duplicate account'")
// 	}

// 	// create 2nd account test
// 	res = testAPI(router, "POST", "/account", "{\"address\":\""+account2.Address.Hex()+"\"}")
// 	assert.Equal(t, http.StatusOK, res.Code, "t3 - create 2nd account")

// 	if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
// 		fmt.Printf("%v", err)
// 	}
// 	if compareAccount(t, resp, account2) {
// 		fmt.Println("PASS  't3 - create 2nd account'")
// 		response[pk2] = account2
// 	} else {
// 		fmt.Println("FAIL  't3 - create 2nd account'")
// 	}

// 	// Get Account
// 	res2 := testAPI(router, "GET", "/account/"+account1.Address.Hex(), "")

// 	assert.Equal(t, http.StatusOK, res2.Code, "t4 - fetch account")
// 	if err := json.Unmarshal(res2.Body.Bytes(), &resp); err != nil {
// 		fmt.Printf("%v", err)
// 	}
// 	if compareAccount(t, resp, account1) {
// 		fmt.Println("PASS  't4 - create account'")
// 	} else {
// 		fmt.Println("FAIL  't4 - create account'")
// 	}
// 	return response
// }

// func compareAccount(t *testing.T, actual, expected types.Account, msgs ...string) bool {
// 	for _, msg := range msgs {
// 		fmt.Println(msg)
// 	}
// 	response := true
// 	response = response && assert.Equalf(t, actual.Address, expected.Address, fmt.Sprintf("Address doesn't match. Expected: %v , Got: %v", expected.Address, actual.Address))
// 	response = response && assert.Equalf(t, actual.IsBlocked, expected.IsBlocked, fmt.Sprintf("Address IsBlocked doesn't match. Expected: %v , Got: %v", expected.IsBlocked, actual.IsBlocked))
// 	response = response && assert.Equalf(t, actual.TokenBalances, expected.TokenBalances, fmt.Sprintf("Balance Tokens doesn't match. Expected: %v , Got: %v", expected.TokenBalances, actual.TokenBalances))

// 	return response
// }
