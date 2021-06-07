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

package services

import (
	"context"

	"github.com/coinbase/rosetta-sdk-go/types"
	ss "github.com/portto/solana-go-sdk/client"
)

// Client is used by the servicers to get block
// data and to submit transactions.
type Client interface {
	Status(context.Context) (
		*types.BlockIdentifier,
		int64,
		*types.SyncStatus,
		[]*types.Peer,
		error,
	)

	Block(
		context.Context,
		*types.PartialBlockIdentifier,
	) (*types.Block, error)

	BlockTransaction(
		context.Context,
		*types.BlockTransactionRequest,
	) (*types.Transaction, error)

	Balance(
		context.Context,
		*types.AccountIdentifier,
		*types.PartialBlockIdentifier,
	) (*types.AccountBalanceResponse, error)

	Call(
		ctx context.Context,
		request *types.CallRequest,
	) (*types.CallResponse, error)
}
type ConstructionMetadata struct {
	BlockHash     string           `json:"blockhash"`
	FeeCalculator ss.FeeCalculator `json:"fee_calculator"`
}

type MetadataWithFee struct {
	Metadata     ConstructionMetadata `json:"metadata"`
	SuggestedFee []*types.Amount      `json:"suggestedFee"`
}
