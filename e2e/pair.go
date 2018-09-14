package e2e

// func testPair(t *testing.T, tokens []types.Token) []types.Pair {
// 	fmt.Printf("\n=== Starting Pair test ===\n")
// 	router := NewRouter()
// 	listPairs := make([]types.Pair, 0)

// 	neededPair := types.Pair{
// 		Name:              strings.ToUpper(tokens[0].Symbol + "/" + tokens[1].Symbol),
// 		BaseTokenAddress:  tokens[1].ContractAddress,
// 		BaseTokenSymbol:   tokens[1].Symbol,
// 		QuoteTokenAddress: tokens[0].ContractAddress,
// 		QuoteTokenSymbol:  tokens[0].Symbol,
// 		Active:            true,
// 		// MakeFee:           big.NewInt(0),
// 		// TakeFee:           big.NewInt(0),
// 	}

// 	// create pair test
// 	res := testAPI(router, "POST", "/pairs", `{"quoteTokenAddress":"`+tokens[0].ContractAddress.Hex()+`", "baseTokenAddress":"`+tokens[1].ContractAddress.Hex()+`", "active":true}`)

// 	assert.Equal(t, http.StatusOK, res.Code, "t1 - create pair")
// 	var resp types.Pair
// 	if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
// 		fmt.Printf("%v", err)
// 	}

// 	if comparePair(t, resp, neededPair) {
// 		fmt.Println("PASS  't1 - create pair'")
// 	} else {
// 		fmt.Println("FAIL  't1 - create pair'")
// 	}

// 	listPairs = append(listPairs, neededPair)

// 	// Duplicate pair test
// 	res = testAPI(router, "POST", "/pairs", `{"quoteTokenAddress":"`+tokens[0].ContractAddress.Hex()+`", "baseTokenAddress":"`+tokens[1].ContractAddress.Hex()+`"}`)
// 	if assert.Equal(t, 401, res.Code, "t2 - create duplicate pair") {
// 		fmt.Println("PASS  't2 - create duplicate pair'")
// 	} else {
// 		fmt.Println("FAIL  't2 - create duplicate pair'")
// 	}

// 	// fetch pair detail test
// 	res = testAPI(router, "GET", "/pairs/"+tokens[1].ContractAddress.Hex()+"/"+tokens[0].ContractAddress.Hex(), "")
// 	assert.Equal(t, http.StatusOK, res.Code, "t2 - fetch pair")
// 	if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
// 		fmt.Printf("%v", err)
// 	}
// 	if comparePair(t, resp, neededPair) {
// 		fmt.Println("PASS  't2 - fetch pair'")
// 	} else {
// 		fmt.Println("FAIL  't2 - fetch pair'")
// 	}

// 	// fetch pairs list
// 	res = testAPI(router, "GET", "/pairs", "")
// 	assert.Equal(t, http.StatusOK, res.Code, "t3 - fetch pair list")
// 	var respArray []types.Pair
// 	if err := json.Unmarshal(res.Body.Bytes(), &respArray); err != nil {
// 		fmt.Printf("%v", err)
// 	}

// 	if assert.Lenf(t, respArray, len(listPairs), fmt.Sprintf("Expected Length: %v Got length: %v", len(listPairs), len(respArray))) {
// 		rb := true
// 		for i, r := range respArray {
// 			if rb = rb && comparePair(t, r, listPairs[i]); !rb {
// 				fmt.Println("FAIL  't3 - fetch pair list'")
// 				break
// 			}
// 		}
// 		if rb {
// 			fmt.Println("PASS  't3 - fetch pair list'")
// 		}
// 	} else {
// 		fmt.Println("FAIL  't3 - fetch pair list'")
// 	}
// 	return listPairs
// }

// func comparePair(t *testing.T, actual, expected types.Pair, msgs ...string) bool {
// 	for _, msg := range msgs {
// 		fmt.Println(msg)
// 	}

// 	response := true
// 	response = response && assert.Equalf(t, actual.Name, expected.Name, fmt.Sprintf("Pair Name doesn't match. Expected: %v , Got: %v", expected.Name, actual.Name))
// 	response = response && assert.Equalf(t, actual.BaseTokenAddress, expected.BaseTokenAddress, fmt.Sprintf("Pair BaseToken Address doesn't match. Expected: %v , Got: %v", expected.BaseTokenAddress, actual.BaseTokenAddress))
// 	response = response && assert.Equalf(t, actual.QuoteTokenAddress, expected.QuoteTokenAddress, fmt.Sprintf("Pair QuoteToken Address doesn't match. Expected: %v , Got: %v", expected.QuoteTokenAddress, actual.QuoteTokenAddress))
// 	response = response && assert.Equalf(t, actual.BaseTokenSymbol, expected.BaseTokenSymbol, fmt.Sprintf("Pair BaseTokenSymbol doesn't match. Expected: %v , Got: %v", expected.BaseTokenSymbol, actual.BaseTokenSymbol))
// 	response = response && assert.Equalf(t, actual.QuoteTokenSymbol, expected.QuoteTokenSymbol, fmt.Sprintf("Pair QuoteTokenSymbol doesn't match. Expected: %v , Got: %v", expected.QuoteTokenSymbol, actual.QuoteTokenSymbol))
// 	response = response && assert.Equalf(t, actual.Active, expected.Active, fmt.Sprintf("Pair Active doesn't match. Expected: %v , Got: %v", expected.Active, actual.Active))
// 	response = response && assert.Equalf(t, actual.MakeFee, expected.MakeFee, fmt.Sprintf("Pair MakerFee doesn't match. Expected: %v , Got: %v", expected.MakeFee, actual.MakeFee))
// 	response = response && assert.Equalf(t, actual.TakeFee, expected.TakeFee, fmt.Sprintf("Pair TakerFee doesn't match. Expected: %v , Got: %v", expected.TakeFee, actual.TakeFee))

// 	return response
// }
