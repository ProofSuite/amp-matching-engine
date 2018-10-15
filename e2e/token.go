package e2e

// func testToken(t *testing.T) []types.Token {
// 	fmt.Printf("\n=== Starting Token test ===\n")
// 	router := NewRouter()
// 	listTokens := make([]types.Token, 0)
// 	dbTokensList := make([]types.Token, 0)

// 	quoteToken := types.Token{
// 		Name:            "HotPotCoin",
// 		Symbol:          "HPC",
// 		Decimal:         18,
// 		ContractAddress: common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
// 		Active:          true,
// 		Quote:           true,
// 	}

// 	baseToken := types.Token{
// 		Name:            "Aura.Test",
// 		Symbol:          "AUT",
// 		Decimal:         18,
// 		ContractAddress: common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
// 		Active:          true,
// 		Quote:           false,
// 	}
// 	wethToken := types.Token{
// 		Name:            "Weth",
// 		Symbol:          "Weth",
// 		Decimal:         18,
// 		ContractAddress: common.HexToAddress("0x2EB24432177e82907dE24b7c5a6E0a5c03226135"),
// 		Active:          true,
// 		Quote:           false,
// 	}
// 	// create token test
// 	res := testAPI(router, "POST", "/tokens", `{  "name":"HotPotCoin", "symbol":"HPC", "decimal":18, "contractAddress":"0x1888a8db0b7db59413ce07150b3373972bf818d3","active":true,"quote":true}`)
// 	assert.Equal(t, http.StatusOK, res.Code, "t1 - create token")

// 	resp := types.Token{}
// 	if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
// 		fmt.Printf("%v", err)
// 	}

// 	if compareToken(t, resp, quoteToken) {
// 		fmt.Println("PASS  't1 - create token'")
// 	} else {
// 		fmt.Println("FAIL  't1 - create token'")
// 	}

// 	listTokens = append(listTokens, quoteToken)
// 	dbTokensList = append(dbTokensList, resp)

// 	// Duplicate token test
// 	res = testAPI(router, "POST", "/tokens", `{  "name":"HotPotCoin", "symbol":"HPC", "decimal":18, "contractAddress":"0x1888a8db0b7db59413ce07150b3373972bf818d3","active":true,"quote":true }`)

// 	if assert.Equal(t, 401, res.Code, "t2 - create duplicate token") {
// 		fmt.Println("PASS  't2 - create duplicate token'")
// 	} else {
// 		fmt.Println("FAIL  't2 - create duplicate token'")
// 	}

// 	// create second token test
// 	res = testAPI(router, "POST", "/tokens", `{  "name":"Aura.Test", "symbol":"AUT", "decimal":18, "contractAddress":"0x2034842261b82651885751fc293bba7ba5398156","active":true }`)
// 	assert.Equal(t, http.StatusOK, res.Code, "t3 - create second token")
// 	if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
// 		fmt.Printf("%v", err)
// 	}

// 	if compareToken(t, resp, baseToken) {
// 		fmt.Println("PASS  't3 - create second token'")
// 	} else {
// 		fmt.Println("FAIL  't3 - create second token'")
// 	}

// 	listTokens = append(listTokens, baseToken)
// 	dbTokensList = append(dbTokensList, resp)

// 	// fetch token detail test
// 	res = testAPI(router, "GET", "/tokens/0x1888a8db0b7db59413ce07150b3373972bf818d3", "")
// 	assert.Equal(t, http.StatusOK, res.Code, "t4 - fetch token")
// 	if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
// 		fmt.Printf("%v", err)
// 	}
// 	if compareToken(t, resp, quoteToken) {
// 		fmt.Println("PASS  't4 - fetch token'")
// 	} else {
// 		fmt.Println("FAIL  't4 - fetch token'")
// 	}

// 	// fetch tokens list
// 	res = testAPI(router, "GET", "/tokens", "")
// 	assert.Equal(t, http.StatusOK, res.Code, "t5 - fetch token list")
// 	var respArray []types.Token
// 	if err := json.Unmarshal(res.Body.Bytes(), &respArray); err != nil {
// 		fmt.Printf("%v", err)
// 	}

// 	if assert.Lenf(t, respArray, len(listTokens), fmt.Sprintf("Expected Length: %v Got length: %v", len(listTokens), len(respArray))) {
// 		rb := true
// 		for i, r := range respArray {
// 			if rb = rb && compareToken(t, r, listTokens[i]); !rb {
// 				fmt.Println("FAIL  't5 - fetch token list'")
// 				break
// 			}
// 		}
// 		if rb {
// 			fmt.Println("PASS  't5 - fetch token list'")
// 		}
// 	} else {
// 		fmt.Println("FAIL  't5 - fetch token list'")
// 	}

// 	// add weth token(for order fees validation)
// 	res = testAPI(router, "POST", "/tokens", `{  "name":"Weth", "symbol":"Weth", "decimal":18, "contractAddress":"0x2EB24432177e82907dE24b7c5a6E0a5c03226135","active":true }`)
// 	assert.Equal(t, http.StatusOK, res.Code, "t6 - create weth token")
// 	if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
// 		fmt.Printf("%v", err)
// 	}

// 	if compareToken(t, resp, wethToken) {
// 		fmt.Println("PASS  't6 - create weth token'")
// 	} else {
// 		fmt.Println("FAIL  't6 - create weth token'")
// 	}
// 	dbTokensList = append(dbTokensList, resp)

// 	// fetch quote tokens list
// 	res = testAPI(router, "GET", "/tokens/quote", "")
// 	assert.Equal(t, http.StatusOK, res.Code, "t7 - fetch quote token list")
// 	if err := json.Unmarshal(res.Body.Bytes(), &respArray); err != nil {
// 		fmt.Printf("%v", err)
// 	}
// 	quoteTokens := []types.Token{quoteToken}
// 	if assert.Lenf(t, respArray, len(quoteTokens), fmt.Sprintf("Expected Length: %v Got length: %v", len(quoteTokens), len(respArray))) {
// 		rb := true
// 		for i, r := range respArray {
// 			if rb = rb && compareToken(t, r, quoteTokens[i]); !rb {
// 				fmt.Println("FAIL  't7 - fetch quote token list'")
// 				break
// 			}
// 		}
// 		if rb {
// 			fmt.Println("PASS  't7 - fetch quote token list'")
// 		}
// 	} else {
// 		fmt.Println("FAIL  't7 - fetch quote token list'")
// 	}

// 	// fetch tokens list
// 	res = testAPI(router, "GET", "/tokens/base", "")
// 	assert.Equal(t, http.StatusOK, res.Code, "t8 - fetch base token list")
// 	if err := json.Unmarshal(res.Body.Bytes(), &respArray); err != nil {
// 		fmt.Printf("%v", err)
// 	}
// 	baseTokens := []types.Token{baseToken, wethToken}
// 	if assert.Lenf(t, respArray, len(baseTokens), fmt.Sprintf("Expected Length: %v Got length: %v", len(baseTokens), len(respArray))) {
// 		rb := true
// 		for i, r := range respArray {
// 			if rb = rb && compareToken(t, r, baseTokens[i]); !rb {
// 				fmt.Println("FAIL  't8 - fetch base token list'")
// 				break
// 			}
// 		}
// 		if rb {
// 			fmt.Println("PASS  't8 - fetch base token list'")
// 		}
// 	} else {
// 		fmt.Println("FAIL  't8 - fetch base token list'")
// 	}
// 	return dbTokensList
// }

// func compareToken(t *testing.T, actual, expected types.Token, msgs ...string) bool {
// 	for _, msg := range msgs {
// 		fmt.Println(msg)
// 	}

// 	response := true
// 	response = response && assert.Equalf(t, actual.Symbol, expected.Symbol, fmt.Sprintf("Token Symbol doesn't match. Expected: %v , Got: %v", expected.Symbol, actual.Symbol))
// 	response = response && assert.Equalf(t, actual.Name, expected.Name, fmt.Sprintf("Token Name doesn't match. Expected: %v , Got: %v", expected.Name, actual.Name))
// 	response = response && assert.Equalf(t, actual.Decimal, expected.Decimal, fmt.Sprintf("Token Decimal doesn't match. Expected: %v , Got: %v", expected.Decimal, actual.Decimal))
// 	response = response && assert.Equalf(t, actual.ContractAddress, expected.ContractAddress, fmt.Sprintf("Token ContractAddress doesn't match. Expected: %v , Got: %v", expected.ContractAddress, actual.ContractAddress))
// 	response = response && assert.Equalf(t, actual.Active, expected.Active, fmt.Sprintf("Token Active doesn't match. Expected: %v , Got: %v", expected.Active, actual.Active))
// 	response = response && assert.Equalf(t, actual.Quote, expected.Quote, fmt.Sprintf("Token Quote doesn't match. Expected: %v , Got: %v", expected.Quote, actual.Quote))

// 	return response
// }
