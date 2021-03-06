<p align="center">
  <a href="https://www.rosetta-api.org">
    <img width="90%" alt="Rosetta" src="https://www.rosetta-api.org/img/rosetta_header.png">
  </a>
</p>

<h3 align="center">
   Rosetta Solana
</h3>

<p align="center"><b>
ROSETTA-SOLANA IS CONSIDERED <a href="https://en.wikipedia.org/wiki/Software_release_life_cycle#Alpha">ALPHA SOFTWARE</a>.
USE AT YOUR OWN RISK!
</b></p>


## Overview
`rosetta-solana-go` provides a reference implementation of the Rosetta API for
Solana in Go. If you haven't heard of the Rosetta API, you can find more
information [here](https://rosetta-api.org).

## Features
* Rosetta API implementation (both Data API and Construction API)
* Stateless, offline, curve-based transaction construction

## Usage
As specified in the [Rosetta API Principles](https://www.rosetta-api.org/docs/automated_deployment.html),
all Rosetta implementations must be deployable via Docker and support running via either an
[`online` or `offline` mode](https://www.rosetta-api.org/docs/node_deployment.html#multiple-modes).


### Direct Install
After cloning this repository, run:
```text
go build -o rosettasolanago
```
### docker
```
docker build -t rosetta-sol-go .
RPC_URL=https://api.mainnet-beta.solana.com PORT=8080 NETWORK=MAINNET MODE=ONLINE docker-compose up
```
## Testing with rosetta-cli
To validate `rosetta-solana`, [install `rosetta-cli`](https://github.com/coinbase/rosetta-cli#install)
and run one of the following commands:
* `rosetta-cli check:data --configuration-file rosetta-cli-conf/testnet/devnet.json`
* `rosetta-cli check:construction --configuration-file rosetta-cli-conf/testnet/devnet.json`

## Development
* `RPC_URL=https://devnet.solana.com PORT=8080 NETWORK=TESTNET MODE=ONLINE go run *.go run` to run server

## Details

### Endpoints Implemented

```
    /network/list (network_list)
    /network/options (network_options)
    /network/status (network_status)
    /account/balance (account_balance)
    /block (get_block)
    /block/transaction (block_transaction)
    /construction/combine (construction_combine)
    /construction/derive (construction_derive)
    /construction/hash (construction_hash)
    /construction/metadata (construction_metadata)
    /construction/payloads (construction_payloads)
    /construction/preprocess (construction_preprocess)
    /construction/submit (construction_submit)
    /construction/parse (construction_parse)
    /call (call)
        
```
#### Environment variables
```
RPC_URL = "https://api.mainnet-beta.solana.com" (optional)
NETWORK = "MAINNET" //MAINNET/TESTNET/DEVNET (required)
PORT = "8080" (optional)
MODE = "ONLINE" //ONLINE/OFFLINE (required)
```

#### Operations supported
See `types::OperationType` to see full list of current operations supported . This list might not be up to date.

```
 
		System__Transfer,
		System__CreateAccount,
		System__Assign,
		System__CreateNonceAccount,
		System__AdvanceNonce,
		System__WithdrawFromNonce,
		System__AuthorizeNonce,
		System__Allocate,
		SplToken__Transfer,
		SplToken__InitializeMint,
		SplToken__InitializeAccount,
		SplToken__CreateToken,
		SplToken__CreateAccount,
		SplToken__Approve,
		SplToken__Revoke,
		SplToken_MintTo,
		SplToken_Burn,
		SplToken_CloseAccount,
		SplToken_FreezeAccount,
		SplToken__TransferChecked,
		Unknown,
```
See https://github.com/imerkle/rosetta-solana-go/blob/master/USAGE.md for examples of request body for every operations
## TODO

 * Add more ops
 * Write tests

## NOTE
This is a go port of https://github.com/imerkle/rosetta-solana originally written in rust.

## License
This project is available open source under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).
