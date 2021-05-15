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
	"context"
	"fmt"
	"math/big"
	"strconv"

	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	dfuserpc "github.com/dfuse-io/solana-go/rpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/ybbus/jsonrpc"
)

type Client struct {
	dfuseRpc *dfuserpc.Client
	rpc      jsonrpc.RPCClient
}

// NewClient creates a Client that from the provided url and params.
func NewClient(url string) (*Client, error) {
	dfuseRpc := dfuserpc.NewClient(url)
	rpc := jsonrpc.NewClient(url)

	return &Client{dfuseRpc, rpc}, nil
}

// Close shuts down the RPC client connection.
func (ec *Client) Close() {
}

func (ec *Client) GetGenesisHash(ctx context.Context) (out *string, err error) {

	params := []interface{}{}

	err = ec.rpc.CallFor(&out, "getGenesisHash", params...)
	return
}
func (ec *Client) getFirstAvailableBlock(ctx context.Context) (out *uint64, err error) {

	params := []interface{}{}

	err = ec.rpc.CallFor(&out, "getFirstAvailableBlock", params...)
	return
}

func (ec *Client) getBlockTime(ctx context.Context, slot uint64) (out *uint64, err error) {

	params := []interface{}{slot}

	err = ec.rpc.CallFor(&out, "getBlockTime", params...)
	return
}

type ClusterNodeResult struct {
	Gossip  string `json:gossip`
	Pubkey  string `json:pubkey`
	Rpc     string `json:rpc`
	Tpu     string `json:tpu`
	Version string `json:version`
}

func (ec *Client) getClusterNodes(ctx context.Context) (out []ClusterNodeResult, err error) {

	params := []interface{}{}

	err = ec.rpc.CallFor(&out, "getClusterNodes", params...)
	return
}

// Status returns geth status information
// for determining node healthiness.
func (ec *Client) Status(ctx context.Context) (
	*RosettaTypes.BlockIdentifier,
	int64,
	[]*RosettaTypes.Peer,
	*RosettaTypes.BlockIdentifier,
	error,
) {
	genesis, _ := ec.GetGenesisHash(ctx)
	index, _ := ec.getFirstAvailableBlock(ctx)

	slot, _ := ec.dfuseRpc.GetSlot(ctx, "")
	slotTime, _ := ec.getBlockTime(ctx, uint64(slot))
	clusterNodes, _ := ec.getClusterNodes(ctx)
	var peers []*RosettaTypes.Peer
	for _, k := range clusterNodes {
		peers = append(peers, &RosettaTypes.Peer{PeerID: k.Pubkey})
	}

	return &RosettaTypes.BlockIdentifier{
			Hash:  strconv.FormatInt(int64(slot), 10), //TODO: should be hash not slot
			Index: int64(slot),
		},
		convertTime(*slotTime),
		peers,
		&RosettaTypes.BlockIdentifier{
			Hash:  *genesis,
			Index: int64(*index),
		},
		nil
}

// SendTransaction injects a signed transaction into the pending pool for execution.
func (ec *Client) SendTransaction(ctx context.Context, txString string) (out string, err error) {
	params := []interface{}{txString}

	err = ec.rpc.CallFor(&out, "sendTransaction", params...)
	return
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}

func (ec *Client) GetRecentBlockhash(ctx context.Context) (out *dfuserpc.GetRecentBlockhashResult, err error) {
	fmt.Println(ec.dfuseRpc.GetRecentBlockhash(ctx, ""))
	return ec.dfuseRpc.GetRecentBlockhash(ctx, "")
}
func (ec *Client) GetConfirmedBlock(ctx context.Context, slot uint64, encoding string) (out *GetConfirmedBlockResult, err error) {
	if encoding == "" {
		encoding = "json"
	}
	params := []interface{}{slot, encoding}

	err = ec.rpc.CallFor(&out, "getConfirmedBlock", params...)
	return
}
func (c *Client) GetConfirmedTransaction(ctx context.Context, signature string) (out dfuserpc.TransactionParsed, err error) {
	params := []interface{}{signature, "jsonParsed"}

	err = c.rpc.CallFor(&out, "getConfirmedTransaction", params...)
	return
}
func (ec *Client) BlockTransaction(
	ctx context.Context,
	blockTransactionRequest *RosettaTypes.BlockTransactionRequest,
) (*RosettaTypes.Transaction, error) {
	tx, err := ec.GetConfirmedTransaction(ctx, blockTransactionRequest.TransactionIdentifier.Hash)
	if err != nil {
		return nil, fmt.Errorf("Error tx")
	}
	rosTx := ToRosTx(tx)
	return &rosTx, nil
}

// Block returns a populated block at the *RosettaTypes.PartialBlockIdentifier.
// If neither the hash or index is populated in the *RosettaTypes.PartialBlockIdentifier,
// the current block is returned.
func (ec *Client) Block(
	ctx context.Context,
	blockIdentifier *RosettaTypes.PartialBlockIdentifier,
) (*RosettaTypes.Block, error) {
	if blockIdentifier != nil {
		if blockIdentifier.Index != nil {
			blockResponse, err := ec.GetConfirmedBlock(ctx, uint64(*blockIdentifier.Index), "jsonParsed")
			if err != nil {
				return nil, err
			}
			return &RosettaTypes.Block{
				BlockIdentifier: &RosettaTypes.BlockIdentifier{
					Index: *blockIdentifier.Index,
					Hash:  blockResponse.Blockhash.String(),
				},
				ParentBlockIdentifier: &RosettaTypes.BlockIdentifier{Index: int64(blockResponse.ParentSlot), Hash: blockResponse.PreviousBlockhash.String()},
				Timestamp:             convertTime(uint64(blockResponse.BlockTime)),
				Transactions:          ToRosTxs(blockResponse.Transactions),
				Metadata:              map[string]interface{}{},
			}, nil
		}
	}
	return nil, fmt.Errorf("block fetch error")
}

func convertTime(time uint64) int64 {
	return int64(time) * 1000
}

//token account

type TokenAccountsOwner struct {
	Pubkey  *string  `json:"pubkey"`
	Account *Account `json:"account"`
}
type TokenAmount struct {
	Amount         string  `json:"amount"`
	Decimals       int32   `json:"decimals"`
	UIAmount       float64 `json:"uiAmount"`
	UIAmountString string  `json:"uiAmountString"`
}
type Info struct {
	Delegate        string      `json:"delegate"`
	DelegatedAmount int64       `json:"delegatedAmount"`
	IsInitialized   bool        `json:"isInitialized"`
	IsNative        bool        `json:"isNative"`
	Mint            string      `json:"mint"`
	Owner           string      `json:"owner"`
	TokenAmount     TokenAmount `json:"tokenAmount"`
}
type Parsed struct {
	AccountType string `json:"accountType"`
	Info        Info   `json:"info"`
}
type Data struct {
	Parsed  Parsed `json:"parsed"`
	Program string `json:"program"`
}
type Account struct {
	Data       Data   `json:"data"`
	Executable bool   `json:"executable"`
	Lamports   int64  `json:"lamports"`
	Owner      string `json:"owner"`
	RentEpoch  int64  `json:"rentEpoch"`
}

type Value struct {
	Account Account `json:"account"`
}
type TokenAccountResponse struct {
	dfuserpc.RPCContext
	Value []Value `json:"value"`
}

func (ec *Client) TokenAccounts(account string) (out *TokenAccountResponse, err error) {

	obj := map[string]interface{}{
		"programId": common.TokenProgramID.ToBase58(),
	}
	obj2 := map[string]interface{}{
		"encoding": "jsonParsed",
	}
	params := []interface{}{account, obj, obj2}
	err = ec.rpc.CallFor(&out, "getTokenAccountsByOwner", params...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (ec *Client) GetSlot(commitment client.Commitment) (out uint64, err error) {
	var params []interface{}
	if commitment != "" {
		params = append(params, string(commitment))
	}

	err = ec.rpc.CallFor(&out, "getSlot", params...)
	return
}
func (ec *Client) GetBlockTime() (out uint64, err error) {
	var params []interface{}

	err = ec.rpc.CallFor(&out, "getBlockTime", params...)
	return
}

// Balance returns the balance of a *RosettaTypes.AccountIdentifier
// at a *RosettaTypes.PartialBlockIdentifier.
//
// We must use graphql to get the balance atomically (the
// rpc method for balance does not allow for querying
// by block hash nor return the block hash where
// the balance was fetched).
func (ec *Client) Balance(
	ctx context.Context,
	account *RosettaTypes.AccountIdentifier,
	block *RosettaTypes.PartialBlockIdentifier,
) (*RosettaTypes.AccountBalanceResponse, error) {

	var symbols []string
	if block != nil {
		return nil, fmt.Errorf("block hash balance not supported")
	}

	bal, err := ec.dfuseRpc.GetBalance(ctx, account.Address, "")
	if err != nil {
		return nil, err
	}
	var balances []*RosettaTypes.Amount
	nativeBalance := &RosettaTypes.Amount{
		Value: fmt.Sprint(bal.Value),
		Currency: &RosettaTypes.Currency{
			Symbol:   Symbol,
			Decimals: Decimals,
			Metadata: nil,
		},
		Metadata: nil,
	}

	tokenAccs, err := ec.TokenAccounts(account.Address)

	if err == nil {
		for _, tokenAcc := range tokenAccs.Value {
			symbol := tokenAcc.Account.Data.Parsed.Info.Mint
			b := &RosettaTypes.Amount{
				Value: fmt.Sprint(bal.Value),
				Currency: &RosettaTypes.Currency{
					Symbol:   symbol,
					Decimals: tokenAcc.Account.Data.Parsed.Info.TokenAmount.Decimals,
					Metadata: nil,
				},
				Metadata: nil,
			}
			balances = append(balances, b)
		}
	}
	if len(symbols) == 0 || Contains(symbols, Symbol) {
		balances = append(balances, nativeBalance)
	}
	slot, err := ec.GetSlot(client.CommitmentFinalized)
	//slotTime, err := ec.GetBlockTime();

	return &RosettaTypes.AccountBalanceResponse{

		BlockIdentifier: &RosettaTypes.BlockIdentifier{
			Hash:  strconv.FormatInt(int64(slot), 10),
			Index: int64(slot),
		},
		Balances: balances,
		Metadata: nil,
	}, nil
}

// Call handles calls to the /call endpoint.
func (ec *Client) Call(
	ctx context.Context,
	request *RosettaTypes.CallRequest,
) (*RosettaTypes.CallResponse, error) {

	var out map[string]interface{}
	param := request.Parameters["param"]
	err := ec.rpc.CallFor(&out, request.Method, param)
	if err != nil {
		return nil, fmt.Errorf("rpc call error")
	}
	return &RosettaTypes.CallResponse{
		Result: out,
	}, nil

}
