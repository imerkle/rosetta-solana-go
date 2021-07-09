// Copyright 2020 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package solanago

import (
	"github.com/coinbase/rosetta-sdk-go/types"
	bin "github.com/dfuse-io/binary"
	"github.com/dfuse-io/solana-go"
	dfuserpc "github.com/dfuse-io/solana-go/rpc"
)

const (
	// NodeVersion is the version of geth we are using.
	NodeVersion = "1.4.17"

	// Blockchain is solanago.
	Blockchain string = "solana"

	// MainnetNetwork is the value of the network
	// in MainnetNetworkIdentifier.
	MainnetNetwork string = "mainnet"

	// TestnetNetwork is the value of the network
	// in TestnetNetworkIdentifier.
	TestnetNetwork string = "testnet"

	// DevnetNetwork is the value of the network
	// in DevnetNetworkIdentifier.
	DevnetNetwork string = "devnet"

	// Symbol is the symbol value
	// used in Currency.
	Symbol = "SOL"

	// Decimals is the decimals value
	// used in Currency.
	Decimals = 9

	// SuccessStatus is the status of any
	// Ethereum operation considered successful.
	SuccessStatus = "SUCCESS"

	// FailureStatus is the status of any
	// Ethereum operation considered unsuccessful.
	FailureStatus = "FAILURE"

	// HistoricalBalanceSupported is whether
	// historical balance is supported.
	HistoricalBalanceSupported = true

	// GenesisBlockIndex is the index of the
	// genesis block.
	GenesisBlockIndex = int64(0)

	Separator    = "__"
	WithNonceKey = "with_nonce"

	MainnetGenesisHash = "5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d"
	TestnetGenesisHash = "4uhcVJyU9pJkvQyS88uRDiswHXSCkY3zQawwpjk2NsNY"
	DevnetGenesisHash  = "EtWTRABZaYq6iMfeYKouRu166VU2xqa1wcaWoxPkrZBG"
)

//op types

const (
	System__Transfer              = "System__Transfer"
	System__CreateAccount         = "System__CreateAccount"
	System__Assign                = "System__Assign"
	System__CreateNonceAccount    = "System__CreateNonceAccount"
	System__AdvanceNonce          = "System__AdvanceNonce"
	System__WithdrawFromNonce     = "System__WithdrawFromNonce"
	System__AuthorizeNonce        = "System__AuthorizeNonce"
	System__Allocate              = "System__Allocate"
	SplToken__Transfer            = "SplToken__Transfer"
	SplToken__InitializeMint      = "SplToken__InitializeMint"
	SplToken__InitializeAccount   = "SplToken__InitializeAccount"
	SplToken__CreateToken         = "SplToken__CreateToken"
	SplToken__CreateAccount       = "SplToken__CreateAccount"
	SplToken__Approve             = "SplToken__Approve"
	SplToken__Revoke              = "SplToken__Revoke"
	SplToken_MintTo               = "SplToken_MintTo"
	SplToken_Burn                 = "SplToken_Burn"
	SplToken_CloseAccount         = "SplToken_CloseAccount"
	SplToken_FreezeAccount        = "SplToken_FreezeAccount"
	SplToken__TransferChecked     = "SplToken__TransferChecked"
	SplToken__CreateAssocTokenAcc = "SplToken__CreateAssocTokenAcc"
	Unknown                       = "Unknown"
)

var (
	// MainnetGenesisBlockIdentifier is the *types.BlockIdentifier
	// of the mainnet genesis block.
	MainnetGenesisBlockIdentifier = &types.BlockIdentifier{
		Hash:  MainnetGenesisHash,
		Index: GenesisBlockIndex,
	}

	// TestnetGenesisBlockIdentifier is the *types.BlockIdentifier
	// of the testnet genesis block.
	TestnetGenesisBlockIdentifier = &types.BlockIdentifier{
		Hash:  TestnetGenesisHash,
		Index: GenesisBlockIndex,
	}
	// TestnetGenesisBlockIdentifier is the *types.BlockIdentifier
	// of the testnet genesis block.
	DevnetGenesisHashBlockIdentifier = &types.BlockIdentifier{
		Hash:  DevnetGenesisHash,
		Index: GenesisBlockIndex,
	}

	// Currency is the *types.Currency for all
	// Ethereum networks.
	Currency = &types.Currency{
		Symbol:   Symbol,
		Decimals: Decimals,
	}

	// OperationTypes are all suppoorted operation types.
	OperationTypes = []string{
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
	}

	// OperationStatuses are all supported operation statuses.
	OperationStatuses = []*types.OperationStatus{
		{
			Status:     SuccessStatus,
			Successful: true,
		},
		{
			Status:     FailureStatus,
			Successful: false,
		},
	}

	// CallMethods are all supported call methods.
	CallMethods = []string{
		"deregisterNode", "validatorExit", "getAccountInfo", "getBalance", "getBlockTime", "getClusterNodes", "getConfirmedBlock", "getConfirmedBlocks", "getConfirmedBlocksWithLimit", "getConfirmedSignaturesForAddress", "getConfirmedSignaturesForAddress2", "getConfirmedTransaction", "getEpochInfo", "getEpochSchedule", "getFeeCalculatorForBlockhash", "getFeeRateGovernor", "getFees", "getFirstAvailableBlock", "getGenesisHash", "getHealth", "getIdentity", "getInflationGovernor", "getInflationRate", "getLargestAccounts", "getLeaderSchedule", "getMinimumBalanceForRentExemption", "getMultipleAccounts", "getProgramAccounts", "getRecentBlockhash", "getSnapshotSlot", "getSignatureStatuses", "getSlot", "getSlotLeader", "getStorageTurn", "getStorageTurnRate", "getSlotsPerSegment", "getStoragePubkeysForSlot", "getSupply", "getTokenAccountBalance", "getTokenAccountsByDelegate", "getTokenAccountsByOwner", "getTokenSupply", "getTotalSupply", "getTransactionCount", "getVersion", "getVoteAccounts", "minimumLedgerSlot", "registerNode", "requestAirdrop", "sendTransaction", "simulateTransaction", "signVote",
	}
)

type TokenParsed struct {
	Decimals        uint64
	Amount          uint64
	MintAutority    solana.PublicKey
	FreezeAuthority solana.PublicKey
	AuthorityType   solana.PublicKey
	NewAuthority    solana.PublicKey
	M               byte
}

type ParsedInstructionMeta struct {
	Authority    string            `json:"authority,omitempty"`
	NewAuthority string            `json:"newAuthority,omitempty"`
	Source       string            `json:"source,omitempty"`
	Destination  string            `json:"destination,omitempty"`
	Mint         string            `json:"mint,omitempty"`
	Decimals     uint8             `json:"decimals,omitempty"`
	TokenAmount  OpMetaTokenAmount `json:"tokenAmount,omitempty"`
	Amount       uint64            `json:"amount,omitempty"`
	Lamports     uint64            `json:"lamports,omitempty"`
	Space        uint64            `json:"space,omitempty"`
}
type OpMetaTokenAmount struct {
	Amount   string  `json:"amount,omitempty"`
	Decimals uint64  `json:"decimals,omitempty"`
	UiAmount float64 `json:"uiAmount,omitempty"`
}

type GetConfirmedBlockResult struct {
	Blockhash         solana.PublicKey             `json:"blockhash"`
	PreviousBlockhash solana.PublicKey             `json:"previousBlockhash"` // could be zeroes if ledger was clean-up and this is unavailable
	ParentSlot        bin.Uint64                   `json:"parentSlot"`
	Transactions      []dfuserpc.TransactionParsed `json:"transactions"`
	Rewards           []dfuserpc.BlockReward       `json:"rewards"`
	BlockTime         bin.Uint64                   `json:"blockTime,omitempty"`
}

type WithNonce struct {
	Account   string `json:"account"`
	Authority string `json:"authority,omitempty"`
}
