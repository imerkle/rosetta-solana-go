package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"crypto/ed25519"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/imerkle/rosetta-solana-go/configuration"
	solanago "github.com/imerkle/rosetta-solana-go/solana"
	"gotest.tools/assert"
)

func TestConstructionServiceSpl(t *testing.T) {
	fromToken := &types.AccountIdentifier{
		Address: "95Dq3sXa3omVjiyxBSD6UMrzPYdmyu6CFCw5wS4rhqgV",
	}
	toToken := &types.AccountIdentifier{
		Address: "GyUjMMeZH3PVXp4tk5sR8LgnVaLTvCPipQ3dQY74k75L",
	}
	c := &types.Currency{
		Symbol:   "3fJRYbtSYZo9SYhwgUBn2zjG98ASy3kuUEnZeHJXqREr",
		Decimals: 2,
	}
	m := map[string]interface{}{
		"authority": "HJGPMwVuqrbm7BDMeA3shLkqdHUru337fgytM7HzqTnH",
	}

	fromSystem := &types.AccountIdentifier{
		Address: "HJGPMwVuqrbm7BDMeA3shLkqdHUru337fgytM7HzqTnH",
	}
	toSystem := &types.AccountIdentifier{
		Address: "42jb8c6XpQ6KXxJEHSWPeoFvyrhuiGvcCJQKumdtW78v",
	}
	cSol := &types.Currency{
		Symbol:   solanago.Currency.Symbol,
		Decimals: solanago.Currency.Decimals,
	}
	mNil := map[string]interface{}{}

	ops0 := []*types.Operation{{
		OperationIdentifier: &types.OperationIdentifier{
			Index: 0,
		},
		Type:    solanago.SplAssociatedTokenAccount__Create,
		Account: fromSystem,
		Metadata: map[string]interface{}{
			"wallet": "42jb8c6XpQ6KXxJEHSWPeoFvyrhuiGvcCJQKumdtW78v",
			"mint":   "GmrqGgTJ2mmNDvqaa39NAnzcwyXtm5ntTa41zPTHyc9o",
		},
	}}
	constructionPipe(t, ops0, false)

	ops1 := []*types.Operation{{
		OperationIdentifier: &types.OperationIdentifier{
			Index: 0,
		},
		Type:    solanago.System__Transfer,
		Account: fromSystem,
		Amount: &types.Amount{
			Value:    "-1",
			Currency: cSol,
		},
		Metadata: mNil,
	}, {
		OperationIdentifier: &types.OperationIdentifier{
			Index: 1,
		},
		Type:    solanago.System__Transfer,
		Account: toSystem,
		Amount: &types.Amount{
			Value:    "1",
			Currency: cSol,
		},
		Metadata: mNil,
	}, {
		OperationIdentifier: &types.OperationIdentifier{
			Index: 2,
		},
		Type:    solanago.SplToken__Transfer,
		Account: fromToken,
		Amount: &types.Amount{
			Value:    "-1",
			Currency: c,
		},
		Metadata: m,
	}, {
		OperationIdentifier: &types.OperationIdentifier{
			Index: 3,
		},
		Type:    solanago.SplToken__Transfer,
		Account: toToken,
		Amount: &types.Amount{
			Value:    "1",
			Currency: c,
		},
		Metadata: m,
	},
	}
	constructionPipe(t, ops1, false)

	ops2 := []*types.Operation{
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 0,
			},
			Type:    solanago.SplToken__TransferNew,
			Account: fromToken,
			Amount: &types.Amount{
				Value:    "-1",
				Currency: c,
			},
			Metadata: m,
		}, {
			OperationIdentifier: &types.OperationIdentifier{
				Index: 1,
			},
			Type: solanago.SplToken__TransferNew,
			Account: &types.AccountIdentifier{
				Address: "9o8WJKYkm71RoGdBziUEPnpPCyW3TgaRpagxBDA9qiiY",
			}, //systemaccount
			Amount: &types.Amount{
				Value:    "1",
				Currency: c,
			},
			Metadata: m,
		},
	}
	constructionPipe(t, ops2, false)

	ops3 := []*types.Operation{
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 0,
			},
			Type:    solanago.SplToken__TransferWithSystem,
			Account: fromSystem,
			Amount: &types.Amount{
				Value:    "-1",
				Currency: c,
			},
			Metadata: mNil,
		}, {
			OperationIdentifier: &types.OperationIdentifier{
				Index: 1,
			},
			Type: solanago.SplToken__TransferWithSystem,
			Account: &types.AccountIdentifier{
				Address: "CZDpZ7KeMansnszdEGZ55C4HjGsMSQBzxPu6jqRm6ZrU",
			}, //systemaccount
			Amount: &types.Amount{
				Value:    "1",
				Currency: c,
			},
			Metadata: mNil,
		},
	}
	constructionPipe(t, ops3, false)
}

func constructionPipe(t *testing.T, ops []*types.Operation, submit bool) {

	ctx := context.Background()
	cfg := configuration.Configuration{
		Mode: "ONLINE",
		Network: &types.NetworkIdentifier{
			Blockchain: solanago.Blockchain,
			Network:    solanago.DevnetNetwork,
		},
		GenesisBlockIdentifier: solanago.TestnetGenesisBlockIdentifier,
		GethURL:                "https://api.devnet.solana.com",
		RemoteGeth:             false,
		Port:                   0,
		GethArguments:          "",
	}
	client, _ := solanago.NewClient("https://api.devnet.solana.com")
	constructionAPIService := NewConstructionAPIService(&cfg, client)

	preRes, err := constructionAPIService.ConstructionPreprocess(ctx, &types.ConstructionPreprocessRequest{
		NetworkIdentifier: cfg.Network,
		Operations:        ops,
		Metadata:          map[string]interface{}{},
	})
	var optsjson map[string]interface{}
	unmarshalJSONMap(preRes.Options, &optsjson)
	metaRes, err := constructionAPIService.ConstructionMetadata(ctx, &types.ConstructionMetadataRequest{
		NetworkIdentifier: cfg.Network,
		Options:           optsjson,
	})
	payRes, err := constructionAPIService.ConstructionPayloads(ctx, &types.ConstructionPayloadsRequest{
		NetworkIdentifier: cfg.Network,
		Operations:        ops,
		Metadata:          metaRes.Metadata,
	})
	if err != nil {
		t.Fatal(err)
	}
	var sigs []*types.Signature

	pub := "f22742d48ce6eeb0c062237b04a5b7f57bfeb8803e9287cd8a112320860e307a"
	priv := "cb1a134c296fbf309d78fe9378c18bc129e5045fbe92d2ad8577ccc84689d4ef"

	p, _ := hex.DecodeString(pub)
	pk, _ := hex.DecodeString(priv)

	for _, v := range payRes.Payloads {
		sigs = append(sigs, &types.Signature{
			SigningPayload: &types.SigningPayload{
				Bytes:         v.Bytes,
				SignatureType: v.SignatureType,
			},
			PublicKey: &types.PublicKey{
				Bytes:     p,
				CurveType: types.Edwards25519,
			},
			SignatureType: v.SignatureType,
			Bytes:         ed25519.Sign(ed25519.NewKeyFromSeed(pk), v.Bytes),
		})
	}
	parseRes, err := constructionAPIService.ConstructionParse(
		ctx, &types.ConstructionParseRequest{
			NetworkIdentifier: cfg.Network,
			Transaction:       payRes.UnsignedTransaction,
		},
	)
	s, _ := json.MarshalIndent(ops, "", "\t")
	fmt.Println("Original")
	fmt.Println(string(s))
	s1, _ := json.MarshalIndent(parseRes.Operations, "", "\t")
	fmt.Println("Parsed")
	fmt.Println(string(s1))

	if !(ops[0].Type == solanago.SplToken__TransferNew || ops[0].Type == solanago.SplToken__TransferWithSystem) {
		assert.Equal(t, ops[0].Type, parseRes.Operations[0].Type)
		assert.Equal(t, len(ops), len(parseRes.Operations))

		assert.Equal(t, ops[0].OperationIdentifier.Index, parseRes.Operations[0].OperationIdentifier.Index)

		assert.Equal(t, ops[0].Account.Address, parseRes.Operations[0].Account.Address)
		if ops[0].Amount != nil {
			assert.Equal(t, ops[0].Amount.Value, parseRes.Operations[0].Amount.Value)
			assert.Equal(t, ops[0].Amount.Currency.Symbol, parseRes.Operations[0].Amount.Currency.Symbol)
		}
	} else {
		if ops[0].Type == solanago.SplToken__TransferNew {
			assert.Equal(t, parseRes.Operations[0].Type, solanago.SplAssociatedTokenAccount__Create)
			assert.Equal(t, len(parseRes.Operations), 3)
		}
	}

	combRes, err := constructionAPIService.ConstructionCombine(
		ctx, &types.ConstructionCombineRequest{
			NetworkIdentifier:   cfg.Network,
			UnsignedTransaction: payRes.UnsignedTransaction,
			Signatures:          sigs,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Println(combRes)

	if submit {
		submitRes, err := constructionAPIService.ConstructionSubmit(ctx, &types.ConstructionSubmitRequest{
			NetworkIdentifier: cfg.Network,
			SignedTransaction: combRes.SignedTransaction,
		})
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(submitRes.TransactionIdentifier.Hash)
	}
}
