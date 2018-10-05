# Websocket API

**Websocket Endpoint**: `/socket`

### PLACE_ORDER (client -> engine)

The PLACE_ORDER message payload consists in an order in the  format. This
message is send in order to provide

**Payload**
```
{
	"channel": "order_channel",
	"message":
	{
		"msgType": "NEW_ORDER",
		"data": {
			"userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
			"amount": 50,
			"price": 2.2,
			"type": 2,
			"buyToken": "0x2034842261b82651885751fc293bba7ba5398156",
			"sellToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3",
      "signature":[string]
		}
	}
}
```
**Response** [engine->client]
If order added to orderbook:

```
{
  "msgType": "ORDER_ADDED",
  "orderId": "5b50d0ac7b44578e7e436815",
  "data": {
    "Order": {
      "id": "5b50d0ac7b44578e7e436815",
      "buyToken": "AUT",
      "sellToken": "HPC",
      "buyTokenAddress": "0x2034842261b82651885751fc293bba7ba5398156",
      "sellTokenAddress": "0x1888a8db0b7db59413ce07150b3373972bf818d3",
      "filledAmount": 0,
      "amount": 5000000000,
      "price": 220000000,
      "fee": 0,
      "type": "SELL",
      "amountBuy": 5000000000,
      "amountSell": 11000000000,
      "exchangeAddress": "",
      "status": "OPEN",
      "pairID": "5b3e82a07b44576ba8000003",
      "pairName": "HPC-AUT",
      "hash": "0xa9a89346cc62330626c5853b74493a1f8e933db582c444bf2288bd6a211586ee",
      "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
      "orderBook": null,
      "createdAt": "2018-07-19T23:25:56.28077166+05:30",
      "updatedAt": "2018-07-19T23:25:56.280771695+05:30"
    },
    "Trades": [],
    "RemainingOrder": {
      "id": "5b50d0ac7b44578e7e436815",
      "buyToken": "AUT",
      "sellToken": "HPC",
      "buyTokenAddress": "0x2034842261b82651885751fc293bba7ba5398156",
      "sellTokenAddress": "0x1888a8db0b7db59413ce07150b3373972bf818d3",
      "filledAmount": 0,
      "amount": 5000000000,
      "price": 220000000,
      "fee": 0,
      "type": "SELL",
      "amountBuy": 5000000000,
      "amountSell": 11000000000,
      "exchangeAddress": "",
      "status": "OPEN",
      "pairID": "5b3e82a07b44576ba8000003",
      "pairName": "HPC-AUT",
      "hash": "0xa9a89346cc62330626c5853b74493a1f8e933db582c444bf2288bd6a211586ee",
      "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
      "orderBook": null,
      "createdAt": "2018-07-19T23:25:56.28077166+05:30",
      "updatedAt": "2018-07-19T23:25:56.280771695+05:30"
    },
    "FillStatus": 1,
    "MatchingOrders": []
  }
}
```

If order filled/partially filled, Then we need to send a new message with payload as signed trades and signed remaining order within a given amount of time, otherwise the match will be reverted and order will be marked as error

```
{
  "msgType": "REQUEST_SIGNATURE",
  "orderId": "5b50d0ff7b44578e7e436816",
  "data": {
    "Order": {
      "id": "5b50d0ff7b44578e7e436816",
      "buyToken": "AUT",
      "sellToken": "HPC",
      "buyTokenAddress": "0x2034842261b82651885751fc293bba7ba5398156",
      "sellTokenAddress": "0x1888a8db0b7db59413ce07150b3373972bf818d3",
      "filledAmount": 6000000000,
      "amount": 6000000000,
      "price": 220000000,
      "fee": 0,
      "type": "BUY",
      "amountBuy": 6000000000,
      "amountSell": 13200000000,
      "exchangeAddress": "",
      "status": "FILLED",
      "pairID": "5b3e82a07b44576ba8000003",
      "pairName": "HPC-AUT",
      "hash": "0xf9f05bf18c6e43c1622b9613707fe931546bf7207341b34cab33a14c0d5a8a10",
      "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
      "orderBook": null,
      "createdAt": "2018-07-19T23:27:19.878444154+05:30",
      "updatedAt": "2018-07-19T23:27:19.878444175+05:30"
    },
    "Trades": [{
      "orderHash": "0x23e38e470bd683414f2fad7916811c35050e43ff3d71b0c053ef5ae22e41708d",
      "amount": 1000000000,
      "price": 220000000,
      "type": "BUY",
      "tradeNonce": 0,
      "taker": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
      "maker": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
      "takerOrderId": "5b50d0ff7b44578e7e436816",
      "makerOrderId": "5b46ea917b445743bb0ffaff",
      "signature": null,
      "hash": "0xf43c261f49a2fe2398830024f32ac86b1ee5abfd9b7299ac1a64df19ec961bd7",
      "pairName": "HPC-AUT",
      "createdAt": "0001-01-01T00:00:00Z",
      "updatedAt": "0001-01-01T00:00:00Z"
    }, {
      "orderHash": "0xa9a89346cc62330626c5853b74493a1f8e933db582c444bf2288bd6a211586ee",
      "amount": 5000000000,
      "price": 220000000,
      "type": "BUY",
      "tradeNonce": 0,
      "taker": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
      "maker": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
      "takerOrderId": "5b50d0ff7b44578e7e436816",
      "makerOrderId": "5b50d0ac7b44578e7e436815",
      "signature": null,
      "hash": "0x0902fa3e6c4e6bd3045fe4081b75cc17c12c2fa6135a30e45d203f2e731133b3",
      "pairName": "HPC-AUT",
      "createdAt": "0001-01-01T00:00:00Z",
      "updatedAt": "0001-01-01T00:00:00Z"
    }],
    "RemainingOrder": null,
    "FillStatus": 3,
    "MatchingOrders": [{
      "Amount": 1000000000,
      "Order": {
        "id": "5b46ea917b445743bb0ffaff",
        "buyToken": "AUT",
        "sellToken": "HPC",
        "buyTokenAddress": "0x2034842261b82651885751fc293bba7ba5398156",
        "sellTokenAddress": "0x1888a8db0b7db59413ce07150b3373972bf818d3",
        "filledAmount": 2000000000,
        "amount": 2000000000,
        "price": 220000000,
        "fee": 0,
        "type": "SELL",
        "amountBuy": 2000000000,
        "amountSell": 4400000000,
        "exchangeAddress": "",
        "status": "FILLED",
        "pairID": "5b3e82a07b44576ba8000003",
        "pairName": "HPC-AUT",
        "hash": "0x23e38e470bd683414f2fad7916811c35050e43ff3d71b0c053ef5ae22e41708d",
        "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
        "orderBook": null,
        "createdAt": "2018-07-12T11:13:45.031963564+05:30",
        "updatedAt": "2018-07-12T11:13:45.031963664+05:30"
      }
    }, {
      "Amount": 5000000000,
      "Order": {
        "id": "5b50d0ac7b44578e7e436815",
        "buyToken": "AUT",
        "sellToken": "HPC",
        "buyTokenAddress": "0x2034842261b82651885751fc293bba7ba5398156",
        "sellTokenAddress": "0x1888a8db0b7db59413ce07150b3373972bf818d3",
        "filledAmount": 5000000000,
        "amount": 5000000000,
        "price": 220000000,
        "fee": 0,
        "type": "SELL",
        "amountBuy": 5000000000,
        "amountSell": 11000000000,
        "exchangeAddress": "",
        "status": "FILLED",
        "pairID": "5b3e82a07b44576ba8000003",
        "pairName": "HPC-AUT",
        "hash": "0xa9a89346cc62330626c5853b74493a1f8e933db582c444bf2288bd6a211586ee",
        "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
        "orderBook": null,
        "createdAt": "2018-07-19T23:25:56.28077166+05:30",
        "updatedAt": "2018-07-19T23:25:56.280771695+05:30"
      }
    }]
  }
}
```

SIGNED_DATA (client -> engine)

To execute an order/trade on the connected ethereum chain, the matching engine requires a cryptographic signature that will be used to verify the users will to trade tokens.

Payload:
```
{
	"channel": "order_channel",
	"message":
	{
  "msgType": "REQUEST_SIGNATURE",
  "orderId": "5b50d0ff7b44578e7e436816",
  "data": {
    "Order": {
      "id": "5b50d0ff7b44578e7e436816",
      "buyToken": "AUT",
      "sellToken": "HPC",
      "buyTokenAddress": "0x2034842261b82651885751fc293bba7ba5398156",
      "sellTokenAddress": "0x1888a8db0b7db59413ce07150b3373972bf818d3",
      "filledAmount": 6000000000,
      "amount": 6000000000,
      "price": 220000000,
      "fee": 0,
      "type": "BUY",
      "amountBuy": 6000000000,
      "amountSell": 13200000000,
      "exchangeAddress": "",
      "status": "FILLED",
      "pairID": "5b3e82a07b44576ba8000003",
      "pairName": "HPC-AUT",
      "hash": "0xf9f05bf18c6e43c1622b9613707fe931546bf7207341b34cab33a14c0d5a8a10",
      "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
      "signature":""
    },
    "Trades": [{
      "orderHash": "0x23e38e470bd683414f2fad7916811c35050e43ff3d71b0c053ef5ae22e41708d",
      "amount": 1000000000,
      "price": 220000000,
      "type": "BUY",
      "tradeNonce": 0,
      "taker": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
      "maker": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
      "takerOrderId": "5b50d0ff7b44578e7e436816",
      "makerOrderId": "5b46ea917b445743bb0ffaff",
      "signature": "",
      "hash": "0xf43c261f49a2fe2398830024f32ac86b1ee5abfd9b7299ac1a64df19ec961bd7",
      "pairName": "HPC-AUT"
    }, {
      "orderHash": "0xa9a89346cc62330626c5853b74493a1f8e933db582c444bf2288bd6a211586ee",
      "amount": 5000000000,
      "price": 220000000,
      "type": "BUY",
      "tradeNonce": 0,
      "taker": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
      "maker": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
      "takerOrderId": "5b50d0ff7b44578e7e436816",
      "makerOrderId": "5b50d0ac7b44578e7e436815",
      "signature": "",
      "hash": "0x0902fa3e6c4e6bd3045fe4081b75cc17c12c2fa6135a30e45d203f2e731133b3",
      "pairName": "HPC-AUT"
    }],
    "RemainingOrder": null,
    "FillStatus": 3,
    "MatchingOrders": [{
      "Amount": 1000000000,
      "Order": {
        "id": "5b46ea917b445743bb0ffaff",
        "buyToken": "AUT",
        "sellToken": "HPC",
        "buyTokenAddress": "0x2034842261b82651885751fc293bba7ba5398156",
        "sellTokenAddress": "0x1888a8db0b7db59413ce07150b3373972bf818d3",
        "filledAmount": 2000000000,
        "amount": 2000000000,
        "price": 220000000,
        "fee": 0,
        "type": "SELL",
        "amountBuy": 2000000000,
        "amountSell": 4400000000,
        "exchangeAddress": "",
        "status": "FILLED",
        "pairID": "5b3e82a07b44576ba8000003",
        "pairName": "HPC-AUT",
        "hash": "0x23e38e470bd683414f2fad7916811c35050e43ff3d71b0c053ef5ae22e41708d",
        "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
        "signature":""
      }
    }, {
      "Amount": 5000000000,
      "Order": {
        "id": "5b50d0ac7b44578e7e436815",
        "buyToken": "AUT",
        "sellToken": "HPC",
        "buyTokenAddress": "0x2034842261b82651885751fc293bba7ba5398156",
        "sellTokenAddress": "0x1888a8db0b7db59413ce07150b3373972bf818d3",
        "filledAmount": 5000000000,
        "amount": 5000000000,
        "price": 220000000,
        "fee": 0,
        "type": "SELL",
        "amountBuy": 5000000000,
        "amountSell": 11000000000,
        "exchangeAddress": "",
        "status": "FILLED",
        "pairName": "HPC-AUT",
        "hash": "0xa9a89346cc62330626c5853b74493a1f8e933db582c444bf2288bd6a211586ee",
        "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
        "signature":""
      }
    }]
  }
}
}
```
CANCEL_ORDER (client -> engine)

To cancel an order (off-chain), the client sends a CANCEL_ORDER message.

Payload:
```
{
	"channel": "order_channel",
	"message":
	{
		"msgType": "CANCEL_ORDER",
		"data": {
		  "id":"5b46eb9f7b445747d1673b21",
      "hash":"",
			"userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63"
      "signature":""
		}
	}
}
```
ORDER_BOOK_SUBSCRIBE (client->engine)

To subscribe to orderbook channel for any given pair. client needs to send message with payload:
**Payload**:
```
{
	"channel": "order_book",
	"message": {
		"event":"subscribe",
		"key":"hpc-aut"
	}
}
```
**Response**
```
{
  "buy": null,
  "sell": [{
    "price": 2.2,
    "volume": 60
  }],
  "trades": [{
    "id": "5b46ead87b445743bb0ffb03",
    "orderHash": "0x23e38e470bd683414f2fad7916811c35050e43ff3d71b0c053ef5ae22e41708d",
    "amount": 2000000000,
    "price": 220000000,
    "type": "BUY",
    "tradeNonce": 0,
    "taker": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
    "maker": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
    "takerOrderId": "5b46ead37b445743bb0ffb02",
    "makerOrderId": "5b46ea8c7b445743bb0ffafd",
    "signature": null,
    "hash": "0x293b6d2aa83841af6e56c1ae86b8fbb953c1f8b19f482fc7b2df64109c320920",
    "pairName": "HPC-AUT",
    "createdAt": "2018-07-12T11:14:56.443+05:30",
    "updatedAt": "2018-07-12T11:14:56.443+05:30"
  }, {
    "id": "5b46ead87b445743bb0ffb04",
    "orderHash": "0x23e38e470bd683414f2fad7916811c35050e43ff3d71b0c053ef5ae22e41708d",
    "amount": 2000000000,
    "price": 220000000,
    "type": "BUY",
    "tradeNonce": 0,
    "taker": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
    "maker": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
    "takerOrderId": "5b46ead37b445743bb0ffb02",
    "makerOrderId": "5b46ea8e7b445743bb0ffafe",
    "signature": null,
    "hash": "0x293b6d2aa83841af6e56c1ae86b8fbb953c1f8b19f482fc7b2df64109c320920",
    "pairName": "HPC-AUT",
    "createdAt": "2018-07-12T11:14:56.443+05:30",
    "updatedAt": "2018-07-12T11:14:56.443+05:30"
  }, {
    "id": "5b46ead87b445743bb0ffb05",
    "orderHash": "0x23e38e470bd683414f2fad7916811c35050e43ff3d71b0c053ef5ae22e41708d",
    "amount": 1000000000,
    "price": 220000000,
    "type": "BUY",
    "tradeNonce": 0,
    "taker": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
    "maker": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
    "takerOrderId": "5b46ead37b445743bb0ffb02",
    "makerOrderId": "5b46ea917b445743bb0ffaff",
    "signature": null,
    "hash": "0xf43c261f49a2fe2398830024f32ac86b1ee5abfd9b7299ac1a64df19ec961bd7",
    "pairName": "HPC-AUT",
    "createdAt": "2018-07-12T11:14:56.443+05:30",
    "updatedAt": "2018-07-12T11:14:56.443+05:30"
  }]
}
```

ORDER_BOOK_UNSUBSCRIBE (client->engine)
To unsubscribe from orderbook channel for any given pair. client needs to send message with payload:
**Payload**
```
{
	"channel": "order_book",
	"message": {
		"event":"unsubscribe",
		"key":"hpc-aut"
	}
}
```

TRADES_SUBSCRIBE (client->engine)
**Payload**
```
{
	"channel": "trades",
	"message": {
		"event":"subscribe",
		"key":"hpc-aut",
		"params":{
		  "from":1531440000,
		  "to":1532075163,
		  "duration":30,
		  "units":"min"
		}
	}
}
```
TRADES_UNSUBSCRIBE (client->engine)
**Payload**
```
{
	"channel": "trades",
	"message": {
		"event":"unsubscribe",
		"key":"hpc-aut",
		"params":{
		  "from":1531440000,
		  "to":1532075163,
		  "duration":30,
		  "units":"min"
		}
	}
}
```
ORDER_PLACED (engine -> client)

Payload:
```
{
  "msgType": "ORDER_ADDED",
  "orderId": "",
  "data": {
    "id": "5b523800dbd89c85aee3dcaf",
    "buyToken": "AUT",
    "sellToken": "HPC",
    "buyTokenAddress": "0x2034842261b82651885751fc293bba7ba5398156",
    "sellTokenAddress": "0x1888a8db0b7db59413ce07150b3373972bf818d3",
    "filledAmount": 0,
    "amount": 6000000000,
    "price": 220000000,
    "fee": 0,
    "type": "SELL",
    "amountBuy": 6000000000,
    "amountSell": 13200000000,
    "exchangeAddress": "",
    "status": "OPEN",
    "pairID": "5b3e82a07b44576ba8000003",
    "pairName": "HPC-AUT",
    "hash": "0xfc4c476cb2166d520a28ec1d4d4a42f1ef1ce5d7ac10433b03a39c99320a69d0",
    "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
    "orderBook": null,
    "createdAt": "2018-07-21T00:59:04.290048711+05:30",
    "updatedAt": "2018-07-21T00:59:04.290048765+05:30"
  }
}
```

ORDER_CANCELED (engine -> client)

To report an order has been successfully cancelled, the matching engine sends
back an ORDER_CANCELED message.

Payload:
```
{
  "order": {
    "id": [Number],
    "exchangeAddress": [String],
    "maker": [String],
    "tokenBuy": [String],
    "tokenSell": [String],
    "amountBuy": [String],
    "amountSell": [String],
    "expires": [String],
    "nonce": [String],
    "feeMake": [String],
    "feeTake": [String],
    "signature": {
      "R": [Number],
      "S": [String],
      "V": [String]
    },
    "pairID": [String],
    "hash": [String]
  }
}
```

ORDER_FILLED (engine -> client)

To report an order has been successfully filled, the matching engine sends
an ORDER_FILLED message to the client. The client needs to follow-up by
signing the given "makerOrder" and return it in a SIGNED_DATA message

Payload:
```
{
  "makerOrder": {
    "id": [Number],
    "exchangeAddress": [String],
    "maker": [String],
    "tokenBuy": [String],
    "tokenSell": [String],
    "amountBuy": [String],
    "amountSell": [String],
    "expires": [String],
    "nonce": [String],
    "feeMake": [String],
    "feeTake": [String],
    "signature": {
      "R": [Number],
      "S": [String],
      "V": [String]
    },
    "pairID": [String],
    "hash": [String]
  },
  "takerOrder": {
    "id": [Number],
    "exchangeAddress": [String],
    "maker": [String],
    "tokenBuy": [String],
    "tokenSell": [String],
    "amountBuy": [String],
    "amountSell": [String],
    "expires": [String],
    "nonce": [String],
    "feeMake": [String],
    "feeTake": [String],
    "signature": {
      "R": [Number],
      "S": [String],
      "V": [String]
    },
    "pairID": [String],
    "hash": [String]
  }
}
```

ORDER_PARTIALLY_FILLED (engine -> client)

To report an order has been partially filled, the matching engine sends an
ORDER_PARTIALLY_FILLED message to the client. The client needs to follow-up
by signing the given "takerOrder" and return it in a SIGNED_DATA message

Payload:
```
{
  "makerOrder": {
    "id": [Number],
    "exchangeAddress": [String],
    "maker": [String],
    "tokenBuy": [String],
    "tokenSell": [String],
    "amountBuy": [String],
    "amountSell": [String],
    "expires": [String],
    "nonce": [String],
    "feeMake": [String],
    "feeTake": [String],
    "signature": {
      "R": [Number],
      "S": [String],
      "V": [String]
    },
    "pairID": [String],
    "hash": [String]
  },
  "takerOrder": {
    "id": [Number],
    "exchangeAddress": [String],
    "maker": [String],
    "tokenBuy": [String],
    "tokenSell": [String],
    "amountBuy": [String],
    "amountSell": [String],
    "expires": [String],
    "nonce": [String],
    "feeMake": [String],
    "feeTake": [String],
    "signature": {
      "R": [Number],
      "S": [String],
      "V": [String]
    },
    "pairID": [String],
    "hash": [String]
  }
}
```

Note: This is currently not completely implemeted I believe.

ORDER_EXECUTED (engine -> client)

The order executed message is sent by the server to the maker client
when a (order/trade) pair is sent the exchange smart contract. This does not mean that the transaction is sucessful but merely that a transaction has been sent.

Payload:
```
{
  "order": {
    "id": [Number],
    "exchangeAddress": [String],
    "maker": [String],
    "tokenBuy": [String],
    "tokenSell": [String],
    "amountBuy": [String],
    "amountSell": [String],
    "expires": [String],
    "nonce": [String],
    "feeMake": [String],
    "feeTake": [String],
    "signature": {
      "R": [Number],
      "S": [String],
      "V": [String]
    },
    "pairID": [String],
    "hash": [String]
  },
  "tx": [String]
}
```
TRADE_EXECUTED: (engine -> client)

The TRADE_EXECUTED message is sent by the server to the taker client
when a (order/trade) pair is sent the exchange smart contract. This does not mean that the transaction is sucessful but merely that a transaction has been sent.

Payload:
```
{
  "order": {
    "id": [Number],
    "exchangeAddress": [String],
    "maker": [String],
    "tokenBuy": [String],
    "tokenSell": [String],
    "amountBuy": [String],
    "amountSell": [String],
    "expires": [String],
    "nonce": [String],
    "feeMake": [String],
    "feeTake": [String],
    "signature": {
      "R": [Number],
      "S": [String],
      "V": [String]
    },
    "pairID": [String],
    "hash": [String],
  },
  "trade": {
    "orderHash": [String],
    "amount": [Number],
    "tradeNonce": [String],
    "taker": [String],
    "signature": {
      "R": [Number],
      "S": [String],
      "V": [String]
    },
    "hash": [String],
    "pairID": [String]
  },
  "tx":
}
```
ORDER_TX_SUCCESS / TRADE_TX_SUCCESS (engine -> client)

The ORDER_TX_SUCCESS and TRADE_TX_SUCCESS messages are sent by the server respectively to the maker client and to the taker client when a trade has been successfully executed on the exchange smart-contract.

Payload:
```
{
  "order": {
    "id": [Number],
    "exchangeAddress": [String],
    "maker": [String],
    "tokenBuy": [String],
    "tokenSell": [String],
    "amountBuy": [String],
    "amountSell": [String],
    "expires": [String],
    "nonce": [String],
    "feeMake": [String],
    "feeTake": [String],
    "signature": {
      "R": [Number],
      "S": [String],
      "V": [String]
    },
    "pairID": [String],
    "hash": [String],
  },
  "trade": {
    "orderHash": [String],
    "amount": [Number],
    "tradeNonce": [String],
    "taker": [String],
    "signature": {
      "R": [Number],
      "S": [String],
      "V": [String]
    },
    "hash": [String],
    "pairID": [String]
  },
  "tx":
}
```

ORDER_TX_ERROR / TRADE_TX_ERROR (engine -> client)

The ORDER_TX_ERROR and TRADE_TX_ERROR messages are sent by the server respectively to the maker client and to the taker client when trade blockchain transaction fails after being sent to the exchange smart contract.

Payload:
```
{
  "order": {
    "id": [Number],
    "exchangeAddress": [String],
    "maker": [String],
    "tokenBuy": [String],
    "tokenSell": [String],
    "amountBuy": [String],
    "amountSell": [String],
    "expires": [String],
    "nonce": [String],
    "feeMake": [String],
    "feeTake": [String],
    "signature": {
      "R": [Number],
      "S": [String],
      "V": [String]
    },
    "pairID": [String],
    "hash": [String],
  },
  "trade": {
    "orderHash": [String],
    "amount": [Number],
    "tradeNonce": [String],
    "taker": [String],
    "signature": {
      "R": [Number],
      "S": [String],
      "V": [String]
    },
    "hash": [String],
    "pairID": [String]
  },
  "tx":
}
```
Note: errorId will be renamed to errorID