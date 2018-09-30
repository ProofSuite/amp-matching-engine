package types

// type PairData struct {
// 	Pair PairAddresses `json:"pair" bson:"_id"`
// 	Open *big.Int `json:"open" bson:"open"`
// 	High *big.Int `json:"high" bson:"high"`
// 	Low *big.Int `json:"low" bson:"low"`
// 	Close *big.Int `json:"close" bson:"close"`
// }

// // MarshalJSON returns the json encoded byte array representing the trade struct
// func (p *PairData) MarshalJSON() ([]byte, error) {
// 	data := map[string]interface{}{
// 		"pair": map[string]interface{}{
// 			"name":       t.Pair.Name,
// 			"baseToken":  t.Pair.BaseToken.Hex(),
// 			"quoteToken": t.Pair.QuoteToken.Hex(),
// 		},
// 		"open":      t.Open.String(),
// 		"high":      t.High.String(),
// 		"low":       t.Low.String(),
// 		"close":     t.Close.String(),
// 		"volume":    t.Volume.String(),
// 		"count":     t.Count.String(),
// 	}

// 	bytes, err := json.Marshal(data)
// 	return bytes, err
// }

// // UnmarshalJSON creates a trade object from a json byte string
// func (p *PairData) UnmarshalJSON(b []byte) error {
// 	pairData := map[string]interface{}{}

// 	err := json.Unmarshal(b, &p)
// 	if err != nil {
// 		return err
// 	}

// 	pair := pairData["pair"].(map[string]interface{})
// 	if pairData["pair"] != nil {
// 		if pair["quoteToken"] == nil {
// 			return errors.New("Quote token is not set")
// 		}

// 		if pair["baseToken"] == nil {
// 			return errors.New("Base token is not set")
// 		}

// 		if pair["name"] == nil {
// 			return errors.New("Pair is not set")
// 		}
// 	}

// 	p.Pair = PairAddresses{}
// 	p.Pair.QuoteToken = common.HexToAddress(pair["quoteToken"].(string))
// 	p.Pair.BaseToken = common.HexToAddress(pair["baseToken"].(string))
// 	p.Pair.Name = common.HexToAddress(pair["name"].(string))

// 	if pairData["open"] != nil {
// 		p.Open = math.ToBigInt(pairData["open"].(string))
// 	}
// 	if pairData["high"] != nil {
// 		p.High = math.ToBigInt(pairData["high"].(string))
// 	}
// 	if pairData["low"] != nil {
// 		p.Low = math.ToBigInt(pairData["low"].(string))
// 	}
// 	if pairData["close"] != nil {
// 		p.Close = math.ToBigInt(pairData["close"].(string))
// 	}
// 	if pairData["volume"] != nil {
// 		p.Volume = math.ToBigInt(pairData["volume"].(string))
// 	}
// 	if pairData["count"] != nil {
// 		p.Volume = math.ToBigInt(pairData["count"].(string))
// 	}
// 	return nil
// }

// func (p *PairData) GetBSON() (interface{}, error) {
// 	count, err := bson.ParseDecimal128(t.Count.String())
// 	if err != nil {
// 		return nil, err
// 	}

// 	o, err := bson.ParseDecimal128(t.Open.String())
// 	if err != nil {
// 		return nil, err
// 	}

// 	h, err := bson.ParseDecimal128(t.High.String())
// 	if err != nil {
// 		return nil, err
// 	}

// 	l, err := bson.ParseDecimal128(t.Low.String())
// 	if err != nil {
// 		return nil, err
// 	}

// 	c, err := bson.ParseDecimal128(t.Close.String())
// 	if err != nil {
// 		return nil, err
// 	}

// 	v, err := bson.ParseDecimal128(t.Volume.String())
// 	if err != nil {
// 		return nil, err
// 	}

// 	return struct {
// 		Pair        PairAddresseRecord             `json:"pair,omitempty" bson:"_id"`
// 		Count     bson.Decimal128 `json:"count" bson:"count"`
// 		Open      bson.Decimal128 `json:"open" bson:"open"`
// 		High      bson.Decimal128 `json:"high" bson:"high"`
// 		Low       bson.Decimal128 `json:"low" bson:"low"`
// 		Close     bson.Decimal128 `json:"close" bson:"close"`
// 		Volume    bson.Decimal128 `json:"volume" bson:"volume"`
// 	}{
// 		ID: PairAddress{
// 			t.ID.Pair,
// 			t.ID.BaseToken.Hex(),
// 			t.ID.QuoteToken.Hex(),
// 		},

// 		Open:      o,
// 		High:      h,
// 		Low:       l,
// 		Close:     c,
// 		Volume:    v,
// 		Count:     count,
// 	}, nil
// }

// func (p *PairDao) SetBSON(raw bson.Raw) error {

// 	decoded := new(struct {
// 		Pair  PairAddressRecord                     `json:"pair,omitempty" bson:"_id"`
// 		Count     bson.Decimal128 `json:"count" bson:"count"`
// 		Open      bson.Decimal128 `json:"open" bson:"open"`
// 		High      bson.Decimal128 `json:"high" bson:"high"`
// 		Low       bson.Decimal128 `json:"low" bson:"low"`
// 		Close     bson.Decimal128 `json:"close" bson:"close"`
// 		Volume    bson.Decimal128 `json:"volume" bson:"volume"`
// 	})

// 	err := raw.Unmarshal(decoded)
// 	if err != nil {
// 		return err
// 	}

// 	t.ID = TickID{
// 		Pair:       decoded.ID.Pair,
// 		BaseToken:  common.HexToAddress(decoded.ID.BaseToken),
// 		QuoteToken: common.HexToAddress(decoded.ID.QuoteToken),
// 	}

// 	t.Count = new(big.Int)
// 	t.Close = new(big.Int)
// 	t.High = new(big.Int)
// 	t.Low = new(big.Int)
// 	t.Open = new(big.Int)
// 	t.Volume = new(big.Int)

// 	count := decoded.Count.String()
// 	o := decoded.Open.String()
// 	h := decoded.High.String()
// 	l := decoded.Low.String()
// 	c := decoded.Close.String()
// 	v := decoded.Volume.String()

// 	t.Count = math.ToBigInt(count)
// 	t.Close = math.ToBigInt(c)
// 	t.High = math.ToBigInt(h)
// 	t.Low = math.ToBigInt(l)
// 	t.Open = math.ToBigInt(o)
// 	t.Volume = math.ToBigInt(v)

// 	t.Timestamp = decoded.Timestamp
// 	return nil
// }
