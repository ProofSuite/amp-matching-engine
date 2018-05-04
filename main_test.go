package main

import (
	"reflect"
	"sync"
	"testing"
)

func TestSingleOrderbook(t *testing.T) {

	actions := make(chan *Action)
	done := make(chan bool)
	engine := NewTradingEngine()
	engine.CreateNewOrderBook("ETHEOS", actions)
	log := make([]*Action, 0)

	go func() {
		for {
			action := <-actions
			log = append(log, action)
			if action.actionType == AT_DONE {
				done <- true
				return
			}
		}
	}()

	engine.AddOrder(&Order{orderType: SELL, symbol: "ETHEOS", id: 1, price: 50, amount: 50})
	engine.AddOrder(&Order{orderType: SELL, symbol: "ETHEOS", id: 2, price: 45, amount: 25})
	engine.AddOrder(&Order{orderType: SELL, symbol: "ETHEOS", id: 3, price: 45, amount: 25})
	engine.AddOrder(&Order{orderType: BUY, symbol: "ETHEOS", id: 4, price: 55, amount: 75})
	engine.CancelOrder(1, "ETHEOS")
	engine.AddOrder(&Order{orderType: BUY, symbol: "ETHEOS", id: 5, price: 55, amount: 20})
	engine.AddOrder(&Order{orderType: BUY, symbol: "ETHEOS", id: 6, price: 50, amount: 15})
	engine.AddOrder(&Order{orderType: SELL, symbol: "ETHEOS", id: 7, price: 45, amount: 25})
	engine.CloseOrderBookChannel("ETHEOS")

	<-done

	expected := []*Action{
		&Action{AT_SELL, "ETHEOS", 1, 0, 50, 50},
		&Action{AT_SELL, "ETHEOS", 2, 0, 25, 45},
		&Action{AT_SELL, "ETHEOS", 3, 0, 25, 45},
		&Action{AT_BUY, "ETHEOS", 4, 0, 75, 55},
		&Action{AT_PARTIAL_FILLED, "ETHEOS", 4, 2, 25, 45},
		&Action{AT_PARTIAL_FILLED, "ETHEOS", 4, 3, 25, 45},
		&Action{AT_FILLED, "ETHEOS", 4, 1, 25, 50},
		&Action{AT_CANCEL, "ETHEOS", 1, 0, 0, 0},
		&Action{AT_CANCELLED, "ETHEOS", 1, 0, 0, 0},
		&Action{AT_BUY, "ETHEOS", 5, 0, 20, 55},
		&Action{AT_BUY, "ETHEOS", 6, 0, 15, 50},
		&Action{AT_SELL, "ETHEOS", 7, 0, 25, 45},
		&Action{AT_PARTIAL_FILLED, "ETHEOS", 7, 5, 20, 55},
		&Action{AT_FILLED, "ETHEOS", 7, 6, 5, 50},
		&Action{AT_DONE, "", 0, 0, 0, 0},
	}

	if !reflect.DeepEqual(log, expected) {
		t.Error("\n\nExpected:\n\n", expected, "\n\nGot:\n\n", log, "\n\n")
	}

}

func TestMultipleOrderbook(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	actions1 := make(chan *Action)
	actions2 := make(chan *Action)
	engine := NewTradingEngine()
	engine.CreateNewOrderBook("ETHEOS", actions1)
	engine.CreateNewOrderBook("ZRXEOS", actions2)
	log1 := make([]*Action, 0)
	log2 := make([]*Action, 0)

	go func() {
		for {
			select {
			case action := <-actions1:
				log1 = append(log1, action)
				if action.actionType == AT_DONE {
					wg.Done()
				}
			case action := <-actions2:
				log2 = append(log2, action)
				if action.actionType == AT_DONE {
					wg.Done()
				}
			}
		}
	}()

	engine.AddOrder(&Order{orderType: SELL, symbol: "ETHEOS", id: 1, price: 50, amount: 50})
	engine.AddOrder(&Order{orderType: SELL, symbol: "ETHEOS", id: 2, price: 45, amount: 25})
	engine.AddOrder(&Order{orderType: SELL, symbol: "ZRXEOS", id: 1, price: 50, amount: 50})
	engine.AddOrder(&Order{orderType: SELL, symbol: "ZRXEOS", id: 2, price: 45, amount: 25})
	engine.AddOrder(&Order{orderType: SELL, symbol: "ETHEOS", id: 3, price: 45, amount: 25})
	engine.AddOrder(&Order{orderType: BUY, symbol: "ETHEOS", id: 4, price: 55, amount: 75})
	engine.AddOrder(&Order{orderType: SELL, symbol: "ZRXEOS", id: 3, price: 45, amount: 25})
	engine.AddOrder(&Order{orderType: BUY, symbol: "ZRXEOS", id: 4, price: 55, amount: 75})
	engine.CancelOrder(1, "ETHEOS")
	engine.AddOrder(&Order{orderType: BUY, symbol: "ETHEOS", id: 5, price: 55, amount: 20})
	engine.CancelOrder(1, "ZRXEOS")
	engine.AddOrder(&Order{orderType: BUY, symbol: "ZRXEOS", id: 5, price: 55, amount: 20})
	engine.AddOrder(&Order{orderType: BUY, symbol: "ZRXEOS", id: 6, price: 50, amount: 15})
	engine.AddOrder(&Order{orderType: SELL, symbol: "ZRXEOS", id: 7, price: 45, amount: 25})
	engine.AddOrder(&Order{orderType: BUY, symbol: "ETHEOS", id: 6, price: 50, amount: 15})
	engine.AddOrder(&Order{orderType: SELL, symbol: "ETHEOS", id: 7, price: 45, amount: 25})
	engine.CloseOrderBookChannel("ETHEOS")
	engine.CloseOrderBookChannel("ZRXEOS")

	wg.Wait()

	expected1 := []*Action{
		&Action{AT_SELL, "ETHEOS", 1, 0, 50, 50},
		&Action{AT_SELL, "ETHEOS", 2, 0, 25, 45},
		&Action{AT_SELL, "ETHEOS", 3, 0, 25, 45},
		&Action{AT_BUY, "ETHEOS", 4, 0, 75, 55},
		&Action{AT_PARTIAL_FILLED, "ETHEOS", 4, 2, 25, 45},
		&Action{AT_PARTIAL_FILLED, "ETHEOS", 4, 3, 25, 45},
		&Action{AT_FILLED, "ETHEOS", 4, 1, 25, 50},
		&Action{AT_CANCEL, "ETHEOS", 1, 0, 0, 0},
		&Action{AT_CANCELLED, "ETHEOS", 1, 0, 0, 0},
		&Action{AT_BUY, "ETHEOS", 5, 0, 20, 55},
		&Action{AT_BUY, "ETHEOS", 6, 0, 15, 50},
		&Action{AT_SELL, "ETHEOS", 7, 0, 25, 45},
		&Action{AT_PARTIAL_FILLED, "ETHEOS", 7, 5, 20, 55},
		&Action{AT_FILLED, "ETHEOS", 7, 6, 5, 50},
		&Action{AT_DONE, "", 0, 0, 0, 0},
	}

	expected2 := []*Action{
		&Action{AT_SELL, "ZRXEOS", 1, 0, 50, 50},
		&Action{AT_SELL, "ZRXEOS", 2, 0, 25, 45},
		&Action{AT_SELL, "ZRXEOS", 3, 0, 25, 45},
		&Action{AT_BUY, "ZRXEOS", 4, 0, 75, 55},
		&Action{AT_PARTIAL_FILLED, "ZRXEOS", 4, 2, 25, 45},
		&Action{AT_PARTIAL_FILLED, "ZRXEOS", 4, 3, 25, 45},
		&Action{AT_FILLED, "ZRXEOS", 4, 1, 25, 50},
		&Action{AT_CANCEL, "ZRXEOS", 1, 0, 0, 0},
		&Action{AT_CANCELLED, "ZRXEOS", 1, 0, 0, 0},
		&Action{AT_BUY, "ZRXEOS", 5, 0, 20, 55},
		&Action{AT_BUY, "ZRXEOS", 6, 0, 15, 50},
		&Action{AT_SELL, "ZRXEOS", 7, 0, 25, 45},
		&Action{AT_PARTIAL_FILLED, "ZRXEOS", 7, 5, 20, 55},
		&Action{AT_FILLED, "ZRXEOS", 7, 6, 5, 50},
		&Action{AT_DONE, "", 0, 0, 0, 0},
	}

	if !reflect.DeepEqual(log1, expected1) {
		t.Error("\n\nExpected:\n\n", expected1, "\n\nGot:\n\n", log1, "\n\n")
	}

	if !reflect.DeepEqual(log2, expected2) {
		t.Error("\n\nExpected:\n\n", expected2, "\n\nGot:\n\n", log2, "\n\n")
	}
}
