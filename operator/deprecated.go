package operator

//TODO what do we do if we don't pick up the event
// case event := <-tradeEvents:
// fmt.Println("TRADE_SUCCESS_EVENT")
// txh := event.Raw.TxHash

// go func() {
// 	_, err := op.EthereumProvider.WaitMined(txh)
// 	if err != nil {
// 		logger.Error(err)
// 	}

// 	takerOrderHash := event.TakerOrderHashes[0]
// 	makerOrderHashes := []common.Hash{}

// 	for _, h := range event.MakerOrderHashes {
// 		makerOrderHashes = append(makerOrderHashes, common.BytesToHash(h[:]))
// 	}

// 	trades, err := op.TradeService.GetByTakerOrderHash(takerOrderHash)
// 	if err != nil {
// 		logger.Error(err)
// 	}

// 	to, err := op.OrderService.GetByHash(takerOrderHash)
// 	if err != nil {
// 		logger.Error(err)
// 	}

// 	makerOrders, err := op.OrderService.GetByHashes(makerOrderHashes)
// 	if err != nil {
// 		logger.Error(err)
// 	}

// 	matches := types.Matches{
// 		MakerOrders: makerOrders,
// 		TakerOrder:  to,
// 		Trades:      trades,
// 	}

// 	err = op.Broker.PublishTradeSuccessMessage(&matches)
// 	if err != nil {
// 		logger.Error(err)
// 	}
// }()

// tradeEvents, err := op.Exchange.ListenToTrades()
// if err != nil {
// 	logger.Error(err)
// 	return err
// }

// tradeEvents, err := op.Exchange.ListenToBatchTrades()
// if err != nil {
// 	logger.Error(err)
// 	return err
// }
