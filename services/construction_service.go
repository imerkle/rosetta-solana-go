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

	"github.com/imerkle/rosetta-solana-go/configuration"
	solanago "github.com/imerkle/rosetta-solana-go/solana"
	"github.com/mr-tron/base58"
	solPCommon "github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/sysprog"
	"github.com/portto/solana-go-sdk/tokenprog"
	solPTypes "github.com/portto/solana-go-sdk/types"

	"github.com/coinbase/rosetta-sdk-go/parser"
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

type WithNonce struct {
	Account   string
	Authority string
}

// ConstructionPreprocess implements the /construction/preprocess
// endpoint.
func (s *ConstructionAPIService) ConstructionPreprocess(
	ctx context.Context,
	request *types.ConstructionPreprocessRequest,
) (*types.ConstructionPreprocessResponse, *types.Error) {
	//TODO: add preprocess
	return &types.ConstructionPreprocessResponse{
		Options: nil,
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
	recentBlockhash, _ := s.client.GetRecentBlockhash(ctx)

	/*
		TODO: Add nonce
		withNonce := request.Options["with_nonce"]
		if withNonce == nil {
		}else{

		}
	*/
	meta, _ := marshalJSONMap(ConstructionMetadata{
		BlockHash:     recentBlockhash.Value.Blockhash.String(),
		FeeCalculator: recentBlockhash.Value.FeeCalculator,
	})
	return &types.ConstructionMetadataResponse{
		Metadata: meta,
		SuggestedFee: []*types.Amount{
			{
				Value:    strconv.FormatInt(int64(recentBlockhash.Value.FeeCalculator.LamportsPerSignature), 10),
				Currency: solanago.Currency,
			},
		},
	}, nil
}

// ConstructionPayloads implements the /construction/payloads endpoint.
func (s *ConstructionAPIService) ConstructionPayloads(
	ctx context.Context,
	request *types.ConstructionPayloadsRequest,
) (*types.ConstructionPayloadsResponse, *types.Error) {
	var instructions []solPTypes.Instruction
	var signers []string

	for i := 0; i < len(request.Operations); i += 2 {
		op := request.Operations[i]

		descriptions := &parser.Descriptions{
			OperationDescriptions: []*parser.OperationDescription{
				{
					Type: op.Type,
					Account: &parser.AccountDescription{
						Exists: true,
					},
					Amount: &parser.AmountDescription{
						Exists: true,
						Sign:   parser.NegativeAmountSign,
						//Currency: solanago.Currency,
					},
				},
				{
					Type: op.Type,
					Account: &parser.AccountDescription{
						Exists: true,
					},
					Amount: &parser.AmountDescription{
						Exists: true,
						Sign:   parser.PositiveAmountSign,
						//Currency: solanago.Currency,
					},
				},
			},
			ErrUnmatched: true,
		}
		matches, err := parser.MatchOperations(descriptions, request.Operations)

		if err != nil {
			return nil, wrapErr(ErrUnclearIntent, err)
		}

		fromOp, _ := matches[0].First()
		fromAdd := fromOp.Account.Address
		toOp, _ := matches[1].First()
		toAdd := toOp.Account.Address

		var ins solPTypes.Instruction

		switch fromOp.Type {
		case solanago.System__Transfer:
			amt, _ := strconv.ParseInt(toOp.Amount.Value, 10, 64)
			ins = sysprog.Transfer(solPCommon.PublicKeyFromString(fromAdd), solPCommon.PublicKeyFromString(toAdd), uint64(amt))
			if !solanago.Contains(signers, fromAdd) {
				signers = append(signers, fromAdd)
			}
			break
		case solanago.SplToken__Transfer:
			authority := fromOp.Metadata["authority"].(string)
			amt, _ := strconv.ParseInt(toOp.Amount.Value, 10, 64)
			if !solanago.Contains(signers, authority) {
				signers = append(signers, authority)
			}
			ins = tokenprog.Transfer(solPCommon.PublicKeyFromString(fromAdd), solPCommon.PublicKeyFromString(toAdd), solPCommon.PublicKeyFromString(authority), []solPCommon.PublicKey{solPCommon.PublicKeyFromString(authority)}, uint64(amt))
			break
		}
		instructions = append(instructions, ins)
	}
	feePayer := solPCommon.PublicKeyFromString(signers[0])
	// Convert map to Metadata struct
	var meta ConstructionMetadata

	//remove
	/*
		m, _ := s.ConstructionMetadata(ctx, &types.ConstructionMetadataRequest{
			NetworkIdentifier: request.NetworkIdentifier,
		})
		request.Metadata = m.Metadata
	*/
	//remove end

	if err := unmarshalJSONMap(request.Metadata, &meta); err != nil {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}
	blockHash := meta.BlockHash
	//TODO: Nonce message
	message := solPTypes.NewMessage(feePayer, instructions, blockHash)
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
	//remove
	/*
		privb, _ := hex.DecodeString("cb1a134c296fbf309d78fe9378c18bc129e5045fbe92d2ad8577ccc84689d4ef")
		pbk, _ := hex.DecodeString("f22742d48ce6eeb0c062237b04a5b7f57bfeb8803e9287cd8a112320860e307a")
		pk := ed25519.PrivateKey(append(privb, pbk...))
		c, _ := s.ConstructionCombine(ctx, &types.ConstructionCombineRequest{
			NetworkIdentifier:   request.NetworkIdentifier,
			UnsignedTransaction: base58.Encode(txUnsigned),
			Signatures: []*types.Signature{
				{
					SigningPayload: signingPayloads[0],
					PublicKey:      &types.PublicKey{Bytes: pbk, CurveType: types.Edwards25519},
					SignatureType:  types.Ed25519,
					Bytes:          ed25519.Sign(pk, msgBytes),
				},
			},
		})
		xxx, errx := s.ConstructionSubmit(ctx, &types.ConstructionSubmitRequest{
			NetworkIdentifier: request.NetworkIdentifier,
			SignedTransaction: c.SignedTransaction,
		})
		fmt.Println(errx)
		fmt.Println(xxx.TransactionIdentifier.Hash)
	*/
	//remove end

	return &types.ConstructionPayloadsResponse{
		UnsignedTransaction: base58.Encode(txUnsigned),
		Payloads:            signingPayloads,
	}, nil
}

func GetSigningKeypairPositions(message solPTypes.Message, pubKeys []solPCommon.PublicKey) ([]uint, *types.Error) {
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
func indexOf(element solPCommon.PublicKey, data []solPCommon.PublicKey) int {
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

	tx, err := GetTxFromStr(request.UnsignedTransaction)
	if err != nil {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}
	var pubKeys []solPCommon.PublicKey
	for _, s := range request.Signatures {
		pubKeys = append(pubKeys, solPCommon.PublicKeyFromBytes(s.PublicKey.Bytes))
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

	tx, err := GetTxFromStr(request.SignedTransaction)
	if err != nil {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}
	hash := base58.Encode(tx.Signatures[0])

	return &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: hash,
		},
	}, nil
}
func GetTxFromStr(t string) (solPTypes.Transaction, error) {
	signedTx, err := base58.Decode(t)
	if err != nil {
		return solPTypes.Transaction{}, err
	}

	tx, err := solPTypes.TransactionDeserialize(signedTx)
	if err != nil {
		return solPTypes.Transaction{}, err
	}
	return tx, nil
}

// ConstructionParse implements the /construction/parse endpoint.
func (s *ConstructionAPIService) ConstructionParse(
	ctx context.Context,
	request *types.ConstructionParseRequest,
) (*types.ConstructionParseResponse, *types.Error) {

	tx, err := GetTxFromStr(request.Transaction)
	if err != nil {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}

	var signers []*types.AccountIdentifier
	if !request.Signed {
		//TODO: list
		signers = []*types.AccountIdentifier{{
			Address: tx.Message.Accounts[0].ToBase58(),
		},
		}
	}
	//solanago.GetRosOperationsFromTx(tx)

	resp := &types.ConstructionParseResponse{
		Operations:               []*types.Operation{},
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
	hash, err := s.client.SendTransaction(ctx, request.SignedTransaction)
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
