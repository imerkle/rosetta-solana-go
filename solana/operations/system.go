package operations

import (
	"encoding/json"

	"github.com/coinbase/rosetta-sdk-go/types"
	solanago "github.com/imerkle/rosetta-solana-go/solana"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/sysprog"
	solPTypes "github.com/portto/solana-go-sdk/types"
)

type SystemOperationMetadata struct {
	Source       string `json:"source,omitempty"`
	Destination  string `json:"destination,omitempty"`
	Space        uint64 `json:"space,omitempty"`
	Lamports     uint64 `json:"lamports,omitempty"`
	NewAuthority string `json:"new_authority,omitempty"`
	Authority    string `json:"authority,omitempty"`
}

func (x *SystemOperationMetadata) SetMeta(op *types.Operation) {
	jsonString, _ := json.Marshal(op.Metadata)
	if x.Lamports == 0 {
		x.Lamports = solanago.ValueToBaseAmount(op.Amount.Value)
	}
	if x.Source == "" {
		x.Source = op.Account.Address
	}
	if x.Authority == "" {
		x.Authority = x.Source
	}
	json.Unmarshal(jsonString, &x)
}
func (x *SystemOperationMetadata) ToInstructions(opType string) []solPTypes.Instruction {
	var ins []solPTypes.Instruction
	switch opType {
	case solanago.System__CreateAccount:
		ins = append(ins, sysprog.CreateAccount(p(x.Source), p(x.Destination), common.TokenProgramID, x.Lamports, x.Space))
		break
	case solanago.System__Assign:
		ins = append(ins, sysprog.Assign(p(x.Source), common.TokenProgramID))
		break
	case solanago.System__Transfer:
		ins = append(ins, sysprog.Transfer(p(x.Source), p(x.Destination), x.Lamports))
		break
	case solanago.System__CreateNonceAccount:
		ins = append(ins, sysprog.CreateAccount(p(x.Source), p(x.Destination), common.SystemProgramID, x.Lamports, sysprog.NonceAccountSize))
		ins = append(ins, solPTypes.Instruction{
			Accounts: []solPTypes.AccountMeta{
				{PubKey: p(x.Destination), IsSigner: false, IsWritable: true},
				{PubKey: common.SysVarRecentBlockhashsPubkey, IsSigner: false, IsWritable: false},
				{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
			},
			ProgramID: common.SystemProgramID,
			Data:      sysprog.InitializeNonceAccount(p(x.Destination), p(x.Authority)).Data,
		})

		break
	case solanago.System__AdvanceNonce:
		ins = append(ins, sysprog.AdvanceNonceAccount(p(x.Destination), p(x.Authority)))
		break
	case solanago.System__WithdrawFromNonce:
		ins = append(ins, sysprog.WithdrawNonceAccount(p(x.Source), p(x.Authority), p(x.Destination), x.Lamports))
		break
	case solanago.System__AuthorizeNonce:
		ins = append(ins, sysprog.AuthorizeNonceAccount(p(x.Destination), p(x.Authority), p(x.NewAuthority)))
		break
	case solanago.System__Allocate:
		ins = append(ins, sysprog.Allocate(p(x.Source), x.Space))
		break
	}
	return ins
}
