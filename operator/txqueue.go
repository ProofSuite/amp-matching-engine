package operator

type TxQueue {
	Name string
	Wallet *types.Wallet
	Exchange *contracts.Exchange
}

func NewTxQueue(n string, w *types.Wallet, ex *contracts.Exchange) {
	txq := &TxQueue{
		Name: n,
		Wallet: w,
		Exchange: ex
	}

	tradeEvents, err := ex.ListenToTrades()
	if err != nil {
		return nil, err
	}

	errorEvents, err := ex.ListenToErrors()
	if err != nil {
		return nil, err
	}

	go txq.HandleEvents()
	return txq, nil
}


// Bug: In certain cases, the trade channel seems to be receiving additional unexpected trades.
	// In the case TestSocketExecuteOrder (in file socket_test.go) is run on its own, everything is working correctly.
	// However, in the case TestSocketExecuteOrder is run among other tests, some tradeLogs do not correspond to an
	// order hash in the ordertrade mapping. I suspect this is because the event listener catches events from previous
	// tests. It might be helpful to see how to listen to events from up to a certain block.
func (txq *TxQueue) HandleEvents() {
	for {
		select {
		case event := <- errorEvents:
			tradeHash := event.tradeHash
			errID := int(event.ErrorId)

			tr, err := txq.TradeService.GetByHash(tradeHash)
			if err != nil {
				log.Print(err)
			}

			err = op.PublishTxErrorMessage(tr, errID)
			if err != nil {
				log.Print(err)
			}

			err = op.PublishTradeCancelMessage(tr)
			if err != nil {
				log.Print(err)
			}

		case event := <- tradeEvents:
			tr, err := tradeService.GetByHash(event.TradeHash)
			if err != nil {
				log.Print(err)
			}

			go func() {
				_, err := op.EthereumService.WaitMined(tr.Tx)
				if err != nil {
					log.Print(err)
				}

				err = op.PublishTradeSuccessMessage(tr)
				if err != nil {
					log.Print(err)
				}

				qname := "TX_QUEUES:" + txq.Name
				ch := getChannel(qname)
				q := getQueue(ch, qname)

				len := q.Messages
				if len > 0 {
					msg, _, _ := ch.Get(
						q.Name,
						true
					)

					pding := &PendingTradeMessage{}
					err = json.Unmarshal(msg.Body, &pendingTrade)
					if err != nil {
						log.Print(err)
					}

					_, err = txq.ExecuteTrade(pding.Order, pding.Trade)
					if err != nil {
						log.Print(err)
					}
				}
			}()
		}
	}
}

func (txq *TxQueue) Length() int {
	name := "TX_QUEUES:" + txq.Name
	ch := getChannel(name)
	q := getQueue(ch, name)
	return q.Messages
}


// AddTradeToExecutionList adds a new trade to the execution list. If the execution list is empty (= contains 1 element
// after adding the transaction hash), the given order/trade pair gets executed. If the tranasction queue is full,
// we return an error. Ultimately we want to account send the transaction to another queue that is handled by another ethereum account
// func (op *Operator) QueueTrade(o *types.Order, t *types.Trade) error {
func (txq *TxQueue) QueueTrade(o *types.Order, tr *types.Trade) error {
	name := "TX_QUEUES:" + txq.Name
	ch := getChannel(name)
	q := getQueue(ch, name)

	bytes, err := json.Marshal(t)
	if err != nil {
		return errors.New("Failed to marshal trade object")
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		ampq.Publishing{
			ContentType: "text/json",
			Body: bytes,
		}
	)

	length := q.Messages
	if length == 1 {
		txq.ExecuteTrade(o, t)
	}

	return nil
}

func (txq *TxQueue) ExecuteTrade(o *types.Order, tr *types.Trade) (*eth.Transaction, error) {
	tx, err := op.Exchange.Trade(o, tr)
	if err != nil {
		return nil, err
	}

	err = op.TradeService.UpdateTradeTx(tr, tx)
	if err != nil {
		return nil, errors.New("Could not update trade tx attribute")
	}

	err = op.
}