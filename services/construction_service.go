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
	"fmt"
	"strconv"
	"strings"

	"github.com/imerkle/rosetta-solana-go/configuration"
	solanago "github.com/imerkle/rosetta-solana-go/solana"
	"github.com/imerkle/rosetta-solana-go/solana/operations"
	"github.com/mitchellh/copystructure"
	"github.com/mr-tron/base58"
	ss "github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	solPTypes "github.com/portto/solana-go-sdk/types"

	"github.com/coinbase/rosetta-sdk-go/types"
)

// ConstructionAPIService implements the server.ConstructionAPIServicer interface.
type ConstructionAPIService struct {
	config *configuration.Configuration
	client *solanago.Client
}

// NewConstructionAPIService creates a new instance of a ConstructionAPIService.
func NewConstructionAPIService(
	cfg *configuration.Configuration,
	client *solanago.Client,
) *ConstructionAPIService {
	return &ConstructionAPIService{
		config: cfg,
		client: client,
	}
}

// ConstructionDerive implements the /construction/derive endpoint.
func (s *ConstructionAPIService) ConstructionDerive(
	ctx context.Context,
	request *types.ConstructionDeriveRequest,
) (*types.ConstructionDeriveResponse, *types.Error) {
	addr := base58.Encode(request.PublicKey.Bytes)
	return &types.ConstructionDeriveResponse{
		AccountIdentifier: &types.AccountIdentifier{
			Address: addr,
		},
	}, nil
}

// ConstructionPreprocess implements the /construction/preprocess
// endpoint.
func (s *ConstructionAPIService) ConstructionPreprocess(
	ctx context.Context,
	request *types.ConstructionPreprocessRequest,
) (*types.ConstructionPreprocessResponse, *types.Error) {
	withNonce, _ := solanago.GetWithNonce(request.Metadata)

	var matchedOperationHashMap map[int64]bool = make(map[int64]bool)

	var SplSystemAccMap map[int64]solanago.SplAccounts = make(map[int64]solanago.SplAccounts)
	for _, op := range request.Operations {
		var cont bool
		var matched *types.Operation
		cont, matched = FindMatch(request.Operations, op, matchedOperationHashMap)
		if cont {
			continue
		}
		if matched != nil && op.Type == solanago.SplToken__TransferWithSystem {
			SplSystemAccMap[op.OperationIdentifier.Index] = solanago.SplAccounts{
				Source:      op.Account.Address,
				Destination: matched.Account.Address,
				Mint:        op.Amount.Currency.Symbol,
			}
			matchedOperationHashMap[op.OperationIdentifier.Index] = true
		}
	}

	return &types.ConstructionPreprocessResponse{
		Options: map[string]interface{}{
			solanago.WithNonceKey:       withNonce,
			solanago.SplSystemAccMapKey: SplSystemAccMap,
		},
	}, nil
}

// ConstructionMetadata implements the /construction/metadata endpoint.
func (s *ConstructionAPIService) ConstructionMetadata(
	ctx context.Context,
	request *types.ConstructionMetadataRequest,
) (*types.ConstructionMetadataResponse, *types.Error) {
	if s.config.Mode != configuration.Online {
		return nil, ErrUnavailableOffline
	}

	var hash string
	var fee ss.FeeCalculator
	withNonce, hasNonce := solanago.GetWithNonce(request.Options)
	if hasNonce {
		acc, _ := s.client.Rpc.GetAccountInfoParsed(ctx, withNonce.Account)
		withNonce.Authority = acc.Data.Nonce.Initialized.Authority
		hash = acc.Data.Nonce.Initialized.BlockHash
		fee = acc.Data.Nonce.Initialized.FeeCalculator
	} else {
		recentBlockhash, _ := s.client.Rpc.GetRecentBlockhash(ctx)
		hash = recentBlockhash.Blockhash
		fee = recentBlockhash.FeeCalculator
	}

	var SplTokenAccMap map[string]solanago.SplAccounts = make(map[string]solanago.SplAccounts)

	if w, ok := request.Options[solanago.SplSystemAccMapKey]; ok {
		w1 := w.(map[string]interface{})
		if err := unmarshalJSONMap(w1, &SplTokenAccMap); err != nil {
			return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
		}

		for k, v := range SplTokenAccMap {

			source, _ := s.client.GetTokenAccountByMint(ctx, v.Source, v.Mint)
			destination, _ := s.client.GetTokenAccountByMint(ctx, v.Destination, v.Mint)
			SplTokenAccMap[k] = solanago.SplAccounts{
				Source:      source,
				Destination: destination,
				Mint:        v.Mint,
			}
		}
	}

	meta, _ := marshalJSONMap(ConstructionMetadata{
		BlockHash:         hash,
		FeeCalculator:     fee,
		SplTokenAccMapKey: SplTokenAccMap,
	})

	return &types.ConstructionMetadataResponse{
		Metadata: meta,
		SuggestedFee: []*types.Amount{
			{
				Value:    strconv.FormatInt(int64(fee.LamportsPerSignature), 10),
				Currency: solanago.Currency,
			},
		},
	}, nil
}
func FindMatch(ops []*types.Operation, op *types.Operation, matchedOperationHashMap map[int64]bool) (bool, *types.Operation) {
	if _, ok := matchedOperationHashMap[op.OperationIdentifier.Index]; ok {
		return true, nil
	}
	var matched *types.Operation = nil
	for _, v := range ops {
		if op.OperationIdentifier.Index == v.OperationIdentifier.Index {
			continue
		}
		if _, ok := matchedOperationHashMap[v.OperationIdentifier.Index]; ok {
			continue
		}
		if v.Type != op.Type {
			continue
		}
		if v.Amount != nil {
			if v.Amount.Currency.Symbol != op.Amount.Currency.Symbol {
				continue
			}
			if solanago.ValueToBaseAmount(v.Amount.Value) != solanago.ValueToBaseAmount(op.Amount.Value) {
				continue
			} else {
				opisNegative := strings.Contains(op.Amount.Value, "-")
				visNegative := strings.Contains(v.Amount.Value, "-")
				if (opisNegative && visNegative) || (!opisNegative && !visNegative) {
					continue
				}
			}
		}
		return false, v
	}
	return false, matched
}

// ConstructionPayloads implements the /construction/payloads endpoint.
func (s *ConstructionAPIService) ConstructionPayloads(
	ctx context.Context,
	request *types.ConstructionPayloadsRequest,
) (*types.ConstructionPayloadsResponse, *types.Error) {
	var instructions []solPTypes.Instruction

	// Convert map to Metadata struct
	var meta ConstructionMetadata

	if err := unmarshalJSONMap(request.Metadata, &meta); err != nil {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}

	var matchedOperationHashMap map[int64]bool = make(map[int64]bool)
	for _, op := range request.Operations {
		var cont bool
		var matched *types.Operation
		cont, matched = FindMatch(request.Operations, op, matchedOperationHashMap)
		if cont {
			continue
		}

		if matched == nil && op.Amount != nil {
			return nil, wrapErr(ErrUnableToParseIntermediateResult, fmt.Errorf("Invalid Operation Request. Please check format"))
		}

		opCopy, err := copystructure.Copy(*op)
		if err != nil {
			return nil, wrapErr(ErrUnclearIntent, fmt.Errorf("Cannot deep copy operations"))

		}
		tp := opCopy.(types.Operation)
		tmpOP := &tp

		if tmpOP.Metadata == nil {
			tmpOP.Metadata = make(map[string]interface{})
		}
		if matched != nil {
			fromOp := tmpOP
			fromAdd := fromOp.Account.Address
			toOp := matched
			toAdd := toOp.Account.Address
			tmpOP.Account = fromOp.Account
			tmpOP.Metadata["source"] = fromAdd
			tmpOP.Metadata["destination"] = toAdd
			tmpOP.Amount = toOp.Amount

			matchedOperationHashMap[fromOp.OperationIdentifier.Index] = true
			matchedOperationHashMap[toOp.OperationIdentifier.Index] = true

		} else {
			matchedOperationHashMap[op.OperationIdentifier.Index] = true
		}
		switch strings.Split(tmpOP.Type, solanago.Separator)[0] {
		case "System":
			s := operations.SystemOperationMetadata{}
			s.SetMeta(tmpOP)
			instructions = append(instructions, (s.ToInstructions(tmpOP.Type))...)

			break
		case "SplToken":
			s := operations.SplTokenOperationMetadata{}
			s.SetMeta(tmpOP, meta.SplTokenAccMapKey)
			instructions = append(instructions, (s.ToInstructions(tmpOP.Type))...)
			break
		case "SplAssociatedTokenAccount":
			s := operations.SplAssociatedTokenAccountOperationMetadata{}
			s.SetMeta(tmpOP)
			instructions = append(instructions, (s.ToInstructions(tmpOP.Type))...)
			break
		default:
			return nil, wrapErr(ErrUnableToParseIntermediateResult, fmt.Errorf("Operation not implemented for construction"))
		}
	}
	signers := solPTypes.GetUniqueSigners(instructions)
	feePayer := common.PublicKeyFromString(signers[0])
	blockHash := meta.BlockHash
	var message solPTypes.Message

	withNonce, hasNonce := solanago.GetWithNonce(request.Metadata)
	if hasNonce {
		message = ss.NewMessageWithNonce(feePayer, instructions, common.PublicKeyFromString(withNonce.Account), common.PublicKeyFromString(withNonce.Authority))
	} else {
		message = solPTypes.NewMessage(feePayer, instructions, blockHash)
	}
	//TODO: use suggestedFee somewhere

	//unsigned signature
	var sig []solPTypes.Signature
	x := make([]byte, 64)
	for i := 0; i < int(message.Header.NumRequireSignatures); i++ {
		sig = append(sig, x)
	}
	tx := solPTypes.Transaction{
		Signatures: sig,
		Message:    message,
	}
	tx.Message.RecentBlockHash = blockHash
	msgBytes, _ := tx.Message.Serialize()
	var signingPayloads []*types.SigningPayload
	for _, sg := range signers {
		signingPayloads = append(signingPayloads, &types.SigningPayload{
			AccountIdentifier: &types.AccountIdentifier{
				Address: sg,
			},
			Bytes:         msgBytes,
			SignatureType: types.Ed25519,
		})
	}

	txUnsigned, err := tx.Serialize()

	if err != nil {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}

	return &types.ConstructionPayloadsResponse{
		UnsignedTransaction: base58.Encode(txUnsigned),
		Payloads:            signingPayloads,
	}, nil
}

func GetSigningKeypairPositions(message solPTypes.Message, pubKeys []common.PublicKey) ([]uint, *types.Error) {
	if len(message.Accounts) < int(message.Header.NumRequireSignatures) {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, fmt.Errorf("invalid positions"))
	}
	signedKeys := message.Accounts[0:message.Header.NumRequireSignatures]
	var positions []uint
	for _, p := range pubKeys {
		index := indexOf(p, signedKeys)
		if index > -1 {
			positions = append(positions, uint(index))
		}
	}
	return positions, nil
}
func indexOf(element common.PublicKey, data []common.PublicKey) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

// ConstructionCombine implements the /construction/combine
// endpoint.
func (s *ConstructionAPIService) ConstructionCombine(
	ctx context.Context,
	request *types.ConstructionCombineRequest,
) (*types.ConstructionCombineResponse, *types.Error) {

	tx, err := solanago.GetTxFromStr(request.UnsignedTransaction)
	if err != nil {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}
	var pubKeys []common.PublicKey
	for _, s := range request.Signatures {
		pubKeys = append(pubKeys, common.PublicKeyFromBytes(s.PublicKey.Bytes))
	}
	positions, errr := GetSigningKeypairPositions(tx.Message, pubKeys)
	if errr != nil {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}
	for i, p := range positions {
		tx.Signatures[p] = request.Signatures[i].Bytes
	}
	signedTx, err := tx.Serialize()
	if err != nil {
		return nil, wrapErr(ErrSignatureInvalid, err)
	}

	return &types.ConstructionCombineResponse{
		SignedTransaction: base58.Encode(signedTx),
	}, nil
}

// ConstructionHash implements the /construction/hash endpoint.
func (s *ConstructionAPIService) ConstructionHash(
	ctx context.Context,
	request *types.ConstructionHashRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {

	tx, err := solanago.GetTxFromStr(request.SignedTransaction)
	if err != nil {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}
	hash := tx.Signatures[0].ToBase58()

	return &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: hash,
		},
	}, nil
}

// ConstructionParse implements the /construction/parse endpoint.
func (s *ConstructionAPIService) ConstructionParse(
	ctx context.Context,
	request *types.ConstructionParseRequest,
) (*types.ConstructionParseResponse, *types.Error) {

	tx, err := solanago.GetTxFromStr(request.Transaction)
	if err != nil {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}

	var signers []*types.AccountIdentifier
	sgns := tx.Message.GetUniqueSigners()
	for _, v := range sgns {
		signers = append(signers, &types.AccountIdentifier{
			Address: v,
		})
	}
	parsedTx, err := solanago.ToParsedTransaction(tx)
	if err != nil {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}

	operations := solanago.GetRosOperationsFromTx(parsedTx, "")

	resp := &types.ConstructionParseResponse{
		Operations:               operations,
		AccountIdentifierSigners: signers,
	}
	return resp, nil
}

// ConstructionSubmit implements the /construction/submit endpoint.
func (s *ConstructionAPIService) ConstructionSubmit(
	ctx context.Context,
	request *types.ConstructionSubmitRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	if s.config.Mode != configuration.Online {
		return nil, ErrUnavailableOffline
	}
	hash, err := s.client.Rpc.SendTransaction(ctx, request.SignedTransaction, ss.SendTransactionConfig{
		SkipPreflight:       false,
		PreflightCommitment: "max",
		Encoding:            "base58",
	})
	if err != nil {
		return nil, wrapErr(ErrBroadcastFailed, err)
	}

	txIdentifier := &types.TransactionIdentifier{
		Hash: hash,
	}
	return &types.TransactionIdentifierResponse{
		TransactionIdentifier: txIdentifier,
	}, nil
}
