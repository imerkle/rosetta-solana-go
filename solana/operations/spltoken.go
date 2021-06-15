package operations

import (
	"encoding/json"
	"strconv"

	"github.com/coinbase/rosetta-sdk-go/types"
	solanago "github.com/imerkle/rosetta-solana-go/solana"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/sysprog"
	"github.com/portto/solana-go-sdk/tokenprog"

	solPTypes "github.com/portto/solana-go-sdk/types"
)

type SplTokenOperationMetadata struct {
	Source          string `json:"source,omitempty"`
	Destination     string `json:"destination,omitempty"`
	Mint            string `json:"mint,omitempty"`
	Authority       string `json:"authority,omitempty"`
	FreezeAuthority string `json:"freeze_authority,omitempty"`
	Amount          uint64 `json:"amount,omitempty"`
	Decimals        uint8  `json:"decimals,omitempty"`
}

func (x *SplTokenOperationMetadata) SetMeta(op *types.Operation) {
	jsonString, _ := json.Marshal(op.Metadata)
	if x.Amount == 0 {
		var amount uint64
		amt, err := strconv.ParseInt(op.Amount.Value, 10, 64)
		if err != nil {
			amount = uint64(amt)
		}
		x.Amount = amount
	}
	if x.Source == "" {
		x.Source = op.Account.Address
	}
	if x.Authority == "" {
		x.Authority = x.Source
	}
	json.Unmarshal(jsonString, &x)
}
func (x *SplTokenOperationMetadata) ToInstructions(opType string) []solPTypes.Instruction {
	var ins []solPTypes.Instruction
	switch opType {
	case solanago.SplToken__InitializeMint:
		ins = append(ins, tokenprog.InitializeMint(x.Decimals, p(x.Mint), p(x.Source), p(x.Authority)))
		break
	case solanago.SplToken__InitializeAccount:
		ins = append(ins, tokenprog.InitializeAccount(p(x.Destination), p(x.Mint), p(x.Source)))

		break
	case solanago.SplToken__CreateToken:
		ins = append(ins, sysprog.CreateAccount(p(x.Source), p(x.Mint), common.TokenProgramID, x.Amount, tokenprog.MintAccountSize))
		ins = append(ins, tokenprog.InitializeMint(x.Decimals, p(x.Mint), p(x.Source), p(x.Authority)))
		break
	case solanago.SplToken__CreateAccount:
		ins = append(ins, sysprog.CreateAccount(p(x.Source), p(x.Destination), common.TokenProgramID, x.Amount, tokenprog.TokenAccountSize))
		ins = append(ins, tokenprog.InitializeAccount(p(x.Destination), p(x.Mint), p(x.Authority)))

		break
	case solanago.SplToken__Approve:
		ins = append(ins, tokenprog.Approve(p(x.Source), p(x.Destination), p(x.Authority), []common.PublicKey{}, x.Amount))
		break
	case solanago.SplToken__Revoke:
		ins = append(ins, tokenprog.Revoke(p(x.Source), p(x.Authority), []common.PublicKey{}))
		break
	case solanago.SplToken_MintTo:
		ins = append(ins, tokenprog.MintTo(p(x.Mint), p(x.Source), p(x.Authority), []common.PublicKey{}, x.Amount))
		break
	case solanago.SplToken_Burn:
		ins = append(ins, tokenprog.Burn(p(x.Source), p(x.Mint), p(x.Authority), []common.PublicKey{}, x.Amount))
		break
	case solanago.SplToken_CloseAccount:
		ins = append(ins, tokenprog.CloseAccount(p(x.Source), p(x.Destination), p(x.Authority), []common.PublicKey{}))
		break
	case solanago.SplToken_FreezeAccount:
		ins = append(ins, tokenprog.ThawAccount(p(x.Source), p(x.Mint), p(x.Authority), []common.PublicKey{}))
		break
	case solanago.SplToken__Transfer:
		ins = append(ins, tokenprog.Transfer(p(x.Source), p(x.Destination), p(x.Authority), []common.PublicKey{}, x.Amount))
		break
	case solanago.SplToken__TransferChecked:
		ins = append(ins, tokenprog.TransferChecked(p(x.Source), p(x.Destination), p(x.Mint), p(x.Authority), []common.PublicKey{}, x.Amount, x.Decimals))
		break

	}
	return ins
}

func p(a string) common.PublicKey {
	return common.PublicKeyFromString(a)
}
