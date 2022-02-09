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
	"strconv"

	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	ss "github.com/portto/solana-go-sdk/client"
)

type Client struct {
	Rpc *ss.Client
}

// NewClient creates a Client that from the provided url and params.
func NewClient(url string) (*Client, error) {
	rpc := ss.NewClient(url)

	return &Client{Rpc: rpc}, nil
}

// Close shuts down the RPC client connection.
func (ec *Client) Close() {
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
	genesis, _ := ec.Rpc.GetGenesisHash(ctx)
	index, _ := ec.Rpc.GetFirstAvailableBlock(ctx)

	bhash, _ := ec.Rpc.GetRecentBlockhash(ctx)
	slot, _ := ec.Rpc.GetSlot(ctx)
	slotTime, _ := ec.Rpc.GetBlockTime(ctx, uint64(slot))
	clusterNodes, _ := ec.Rpc.GetClusterNodes(ctx)
	var peers []*RosettaTypes.Peer
	for _, k := range clusterNodes {
		peers = append(peers, &RosettaTypes.Peer{PeerID: k.Pubkey})
	}
	return &RosettaTypes.BlockIdentifier{
			Hash:  bhash.Blockhash,
			Index: int64(slot),
		},
		convertTime(uint64(slotTime)),
		peers,
		&RosettaTypes.BlockIdentifier{
			Hash:  genesis,
			Index: int64(index),
		},
		nil
}

func (ec *Client) BlockTransaction(
	ctx context.Context,
	blockTransactionRequest *RosettaTypes.BlockTransactionRequest,
) (*RosettaTypes.Transaction, error) {
	tx, err := ec.Rpc.GetConfirmedTransactionParsed(ctx, blockTransactionRequest.TransactionIdentifier.Hash)
	if err != nil {
		return nil, err
	}
	rosTx := ToRosTx(tx.Transaction)
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
			blockResponse, err := ec.Rpc.GetConfirmedBlockParsed(ctx, uint64(*blockIdentifier.Index))
			if err != nil {
				return nil, err
			}
			return &RosettaTypes.Block{
				BlockIdentifier: &RosettaTypes.BlockIdentifier{
					Index: *blockIdentifier.Index,
					Hash:  blockResponse.Blockhash,
				},
				ParentBlockIdentifier: &RosettaTypes.BlockIdentifier{Index: int64(blockResponse.ParentSlot), Hash: blockResponse.PreviousBlockhash},
				Timestamp:             convertTime(uint64(blockResponse.BlockTime)),
				Transactions:          ToRosTxs(blockResponse.Transactions),
				Metadata:              map[string]interface{}{},
			}, nil
		}
	}
	return nil, fmt.Errorf("block fetch error")
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

	//var symbols []string
	if block != nil {
		return nil, fmt.Errorf("block hash balance not supported")
	}

	bal, err := ec.Rpc.GetBalance(ctx, account.Address)
	if err != nil {
		return nil, err
	}
	var balances []*RosettaTypes.Amount
	nativeBalance := &RosettaTypes.Amount{
		Value: fmt.Sprint(bal),
		Currency: &RosettaTypes.Currency{
			Symbol:   Symbol,
			Decimals: Decimals,
			Metadata: nil,
		},
		Metadata: nil,
	}

	//tokenAccs, err := ec.Rpc.GetTokenAccountsByOwner(ctx, account.Address)
	//
	//if err == nil {
	//	for _, tokenAcc := range tokenAccs {
	//		symbol := tokenAcc.Account.Data.Parsed.Info.Mint
	//		b := &RosettaTypes.Amount{
	//			Value: tokenAcc.Account.Data.Parsed.Info.TokenAmount.Amount,
	//			Currency: &RosettaTypes.Currency{
	//				Symbol:   symbol,
	//				Decimals: tokenAcc.Account.Data.Parsed.Info.TokenAmount.Decimals,
	//				Metadata: nil,
	//			},
	//			Metadata: nil,
	//		}
	//		balances = append(balances, b)
	//	}
	//}
	//if len(symbols) == 0 || Contains(symbols, Symbol) {
	balances = append(balances, nativeBalance)
	//}
	slot, err := ec.Rpc.GetSlot(ctx)

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
	var x []interface{}
	if p, ok := request.Parameters["param"]; ok {
		x = []interface{}{p}
	} else {
		x = []interface{}{}
	}

	out, err := ec.Rpc.CallRequest(ctx, request.Method, x)

	if err != nil {
		return nil, fmt.Errorf("rpc call error")
	}

	res := make(map[string]interface{})
	if _, ok := out.([]interface{}); ok {
		res["result"] = out
	} else {
		res["result"] = out.(map[string]interface{})
	}

	return &RosettaTypes.CallResponse{
		Result: res,
	}, nil

}

func (ec *Client) GetTokenAccountByMint(ctx context.Context, owner string, mint string) (string, error) {
	tokenAccs, err := ec.Rpc.GetTokenAccountByMint(ctx, owner, mint)
	if err != nil || len(tokenAccs) == 0 {
		return "", fmt.Errorf("No Token Account Found")
	}
	return tokenAccs[0].Pubkey, nil
}
