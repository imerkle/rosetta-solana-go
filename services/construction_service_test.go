package services

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"

	"crypto/ed25519"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/imerkle/rosetta-solana-go/configuration"
	solanago "github.com/imerkle/rosetta-solana-go/solana"
)

func TestConstructionServiceSpl(t *testing.T) {
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

	from := &types.AccountIdentifier{
		Address: "95Dq3sXa3omVjiyxBSD6UMrzPYdmyu6CFCw5wS4rhqgV",
	}
	to := &types.AccountIdentifier{
		Address: "GyUjMMeZH3PVXp4tk5sR8LgnVaLTvCPipQ3dQY74k75L",
	}
	c := &types.Currency{
		Symbol:   "3fJRYbtSYZo9SYhwgUBn2zjG98ASy3kuUEnZeHJXqREr",
		Decimals: 2,
	}
	m := map[string]interface{}{
		"authority": "HJGPMwVuqrbm7BDMeA3shLkqdHUru337fgytM7HzqTnH",
	}

	from1 := &types.AccountIdentifier{
		Address: "HJGPMwVuqrbm7BDMeA3shLkqdHUru337fgytM7HzqTnH",
	}
	to1 := &types.AccountIdentifier{
		Address: "42jb8c6XpQ6KXxJEHSWPeoFvyrhuiGvcCJQKumdtW78v",
	}
	c1 := &types.Currency{
		Symbol:   solanago.Currency.Symbol,
		Decimals: solanago.Currency.Decimals,
	}
	m1 := map[string]interface{}{}

	ops := []*types.Operation{&types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: 0,
		},
		Type:    "System__Transfer",
		Account: from1,
		Amount: &types.Amount{
			Value:    "-1",
			Currency: c1,
		},
		Metadata: m1,
	}, &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: 1,
		},
		Type:    "System__Transfer",
		Account: to1,
		Amount: &types.Amount{
			Value:    "1",
			Currency: c1,
		},
		Metadata: m1,
	},
		&types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 2,
			},
			Type:    "SplToken__Transfer",
			Account: from,
			Amount: &types.Amount{
				Value:    "-1",
				Currency: c,
			},
			Metadata: m,
		}, &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 3,
			},
			Type:    "SplToken__Transfer",
			Account: to,
			Amount: &types.Amount{
				Value:    "1",
				Currency: c,
			},
			Metadata: m,
		},
	}
	preRes, err := constructionAPIService.ConstructionPreprocess(ctx, &types.ConstructionPreprocessRequest{
		NetworkIdentifier: cfg.Network,
		Operations:        ops,
		Metadata:          map[string]interface{}{},
	})
	metaRes, err := constructionAPIService.ConstructionMetadata(ctx, &types.ConstructionMetadataRequest{
		NetworkIdentifier: cfg.Network,
		Options:           preRes.Options,
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
	//assert.Equal(t, ops, parseRes.Operations)
	fmt.Println(len(parseRes.Operations))
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
	submitRes, err := constructionAPIService.ConstructionSubmit(ctx, &types.ConstructionSubmitRequest{
		NetworkIdentifier: cfg.Network,
		SignedTransaction: combRes.SignedTransaction,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(submitRes.TransactionIdentifier.Hash)
}
