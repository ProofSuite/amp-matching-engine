# matching-engine
Official proof decentralized exchange matching engine

# The Proof Decentralized Exchange

The proof decentralized exchange is a hybrid decentralized exchange that aims at bringing together the ease of use of centralized exchanges along with the security and privacy features of decentralized exchanges. Orders are matched through the an off-chain orderbook. After orders are matched and signed, the decentralized exchange operator has the sole ability to perform a transaction to the smart contract. This provides for the best UX as the exchange operator is the only party having to interact directly with the blockchain. Exchange users simply sign orders which are broadcasted then to the orderbook. This design enables users to queue and cancel their orders seamlessly.


# Getting Started

## Requirements

- **mongoDB** version 3.6 or newer
- **rabbitmq** version 3.7.7 or newer
- **redis** version 4.0 or newer

## Booting up the server



# Types

## Orders

Orders contain the information that is required to register an order in the orderbook as a "Maker".

- **id** is the primary ID of the order (possibly deprecated)
- **orderType** is either BUY or SELL. It is currently not parsed by the server and compute directly from tokenBuy, tokenSell, amountBuy, amountSell
- **exchangeAddress** is the exchange smart contract address
- **maker** is the maker (usually sender) ethereum account address
- **tokenBuy** is the BUY token ethereum address
- **tokenSell** is the SELL token ethereum address
- **amountBuy** is the BUY amount (in BUY_TOKEN units)
- **amountSell** is the SELL amount (in SELL_TOKEN units)
- **expires** is the order expiration timestamp
- **nonce** is the nonce that corresponds to
- **feeMake** is the maker fee (not implemented yet)
- **feeTake** is the taker fee (not implemented yet)
- **pairID** is a hash of the corresponding
- **hash** is a hash of the order details (see details below)
- **signature** is a signature of the order hash. The signer must equal to the maker address for the order to be valid.
- **price** corresponds to the pricepoint computed by the matching engine (not parsed)
- **amount** corresponds to the amount computed by the matching engine (not parsed)

**Order Price and Amount**

There are two ways to describe the amount of tokens being bought/sold. The smart-contract requires (tokenBuy, tokenSell, amountBuy, amountSell) while the
orderbook requires (pairID, amount, price).

The conversion between both systems can be found in the engine.ComputeOrderPrice
function


**Order Hash**

The order hash is a sha-256 hash of the following elements:
- Exchange address
- Token Buy address
- Amount Buy
- Token Sell Address
- Amount Sell
- Expires
- Nonce
- Maker Address


## Trades

When an order matches another order in the orderbook, the "taker" is required
to sign a trade object that matches an order.

- **orderHash** is the hash of the matching order
- **amount** is the amount of tokens that will be traded
- **trade nonce** is a unique integer to distinguish successive but identical orders (note: can probably be renamed to nonce)
- **taker** is the taker ethereum account address
- **pairID** is a hash identifying the token pair that will be traded
- **hash** is a unique identifier hash of the trade details (see details below)
- **signature** is a signature of the trade hash

Trade Hash:

The trade hash is a sha-256 hash of the following elements:
- Order Hash
- Amount
- Taker Address
- Trade Nonce


The (Order, Trade) tuple can then be used to perform an on-chain transaction for this trade.

## Quote Tokens and Token Pairs

In the same way as traditional exchanges function with the idea of base
currencies and quote currencies, the Proof decentralized exchange works with
base tokens and quote tokens under the following principles:

- Only the exchange operator can register a quote token
- Anybody can register a token pair (but the quote token needs to be registered)

Token pairs are identified by an ID (a hash of both token addresses)



# Websocket API

See [WEBSOCKET_API.md](WEBSOCKET_API.md)


# Contribution

Thank you for considering helping the Proof project !

To make the Proof project truely revolutionary, we need and accept contributions from anyone and are grateful even for the smallest fixes.

If you want to help Proof, please fork and setup the development environment of the appropriate repository. In the case you want to submit substantial changes, please get in touch with our development team on our slack channel (slack.proofsuite.com) to verify those modifications are in line with the general goal of the project and receive early feedback. Otherwise you are welcome to fix, commit and send a pull request for the maintainers to review and merge into the main code base.

Please make sure your contributions adhere to our coding guidelines:

Code must adhere as much as possible to standard conventions (DRY - Separation of concerns - Modular)
Pull requests need to be based and opened against the master branch
Commit messages should properly describe the code modified
Ensure all tests are passing before submitting a pull request

# License

The Proof CryptoFiat smart contract (i.e. all code inside of the contracts and test directories) is licensed under the MIT License, also included in our repository in the LICENSE file.