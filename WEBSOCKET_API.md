# Websocket API

PLACE_ORDER (client -> engine)

The PLACE_ORDER message payload consists in an order in the json format. This
message is send in order to provide

```json
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
      "V": [String],
    },
    "pairID": [String],
    "hash": [String]
  }
}
```


SIGNED_DATA (client -> engine)

To execute an order/trade on the connected ethereum chain, the matching engine requires a cryptographic signature that will be used to verify the users will to trade tokens.

Payload:
```json
{
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
  }
}
```
CANCEL_ORDER (client -> engine)

To cancel an order (off-chain), the client sends a CANCEL_ORDER message.

Payload:
```json
{
  "orderCancel": {
    "orderId": [Number],
    "pair": [String],
    "orderHash": [String],
    "hash": [String],
    "signature": {
      "R": [Number],
      "S": [String],
      "V": [String]
    },
  }
}
```

ORDER_PLACED (engine -> client)

Payload:
```json
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

ORDER_CANCELED (engine -> client)

To report an order has been successfully canceled, the matching engine sends
back an ORDER_CANCELED message.

Payload:
```json
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
```json
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
```json
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
```json
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
```json
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
```json
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
```json
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