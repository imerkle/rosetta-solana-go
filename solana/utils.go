package solanago

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/coinbase/rosetta-sdk-go/types"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/dfuse-io/solana-go"
	"github.com/dfuse-io/solana-go/programs/system"
	"github.com/dfuse-io/solana-go/programs/token"
	"github.com/dfuse-io/solana-go/rpc"
	dfuserpc "github.com/dfuse-io/solana-go/rpc"
	"github.com/iancoleman/strcase"
)

func IsBalanceChanging(opType string) bool {
	a := false
	switch opType {
	case "System__CreateAccount", "System__WithdrawFromNonce", "System__Transfer", "SplToken__Transfer", "SplToken__TransferChecked", "Stake__Split", "Stake__Withdraw", "Vote__Withdraw", "SplToken__InitializeAccount":
		a = true
	}
	return a
}

func getOperationTypeWithProgram(program string, s string) string {
	toPascal := strcase.ToCamel(program)

	newStr := fmt.Sprint(
		toPascal,
		Separator,
		strcase.ToCamel(s),
	)
	return newStr
}
func getOperationType(s string) string {
	x := strings.Split(s, Separator)
	if len(x) < 2 {
		return Unknown
	}
	return getOperationTypeWithProgram(x[0], x[1])
}
func split_at(at int, input []byte) ([]byte, []byte) {
	return input[0:1], input[1:len(input)]
}
func PublicKeyFromBytes(in []byte) (out solana.PublicKey) {
	byteCount := len(in)
	if byteCount == 0 {
		return
	}

	max := 32
	if byteCount < max {
		max = byteCount
	}

	copy(out[:], in[0:max])
	return
}

func parseInstruction(input []byte) TokenParsed {
	tag, rest := split_at(0, input)
	var amount []byte
	var decimals []byte
	var mintAuthority []byte
	var freezeAuthority []byte
	var authorityType []byte
	var newAuthority []byte
	var m byte
	switch tag[0] {
	case 0:
		decimals, rest = split_at(0, rest)
		mintAuthority, rest, _ = unpackPubkey(rest)
		freezeAuthority, _, _ = unpackPubkeyOption(rest)
		break
	case 2:
		m = rest[0]
		break
	case 3:
	case 4:
	case 7:
	case 8:
		amount = rest[0:8]
		break
	case 6:
		authorityType, rest = split_at(0, rest)
		newAuthority, _, _ = unpackPubkeyOption(rest)

		break
	}
	return TokenParsed{
		MintAutority:    PublicKeyFromBytes(mintAuthority),
		FreezeAuthority: PublicKeyFromBytes(freezeAuthority),
		AuthorityType:   PublicKeyFromBytes(authorityType),
		NewAuthority:    PublicKeyFromBytes(newAuthority),
		Decimals:        binary.BigEndian.Uint64(decimals),
		Amount:          binary.LittleEndian.Uint64(amount),
		M:               m,
	}
}

func unpackPubkey(input []byte) ([]byte, []byte, error) {
	if len(input) >= 32 {
		key, rest := split_at(32, input)
		return key, rest, nil
	} else {
		return nil, nil, fmt.Errorf("Invalid instruction")
	}
}

func unpackPubkeyOption(input []byte) ([]byte, []byte, error) {
	f, rest := split_at(0, input)
	switch f[0] {
	case 0:
		return nil, rest, nil
		break
	case 1:
		return unpackPubkey(rest)
		break
	}
	return nil, nil, fmt.Errorf("Invalid instruction")
}
func GetRosOperationsFromTx(tx rpc.TransactionParsed, status string) []*types.Operation {
	//	hash := tx.Transaction.Signatures[0].String()
	opIndex := int64(0)
	var operations []*types.Operation
	for _, ins := range tx.Transaction.Message.Instructions {

		oi := types.OperationIdentifier{
			Index: opIndex,
		}
		opIndex += 1

		if !ins.IsParsed() {

			var inInterface map[string]interface{}
			inrec, _ := json.Marshal(ins)
			json.Unmarshal(inrec, &inInterface)

			operations = append(operations, &types.Operation{
				OperationIdentifier: &oi,
				Type:                Unknown,
				Status:              &status,
				Metadata:            inInterface,
			})
		} else {
			opType := getOperationTypeWithProgram(ins.Program, ins.Parsed.InstructionType)
			if !Contains(OperationTypes, opType) {
				opType = "Unknown"
			}
			jsonString, _ := json.Marshal(ins.Parsed.Info)

			parsedInstructionMeta := ParsedInstructionMeta{}
			var parsedInstructionMetaInterface interface{}
			json.Unmarshal(jsonString, &parsedInstructionMeta)
			json.Unmarshal(jsonString, &parsedInstructionMetaInterface)

			var inInterface map[string]interface{}
			inrec, _ := json.Marshal(parsedInstructionMetaInterface)
			json.Unmarshal(inrec, &inInterface)

			fmt.Println(parsedInstructionMetaInterface)
			if IsBalanceChanging(opType) {
				if parsedInstructionMeta.Mint == "" {
					parsedInstructionMeta.Mint = Symbol
				}
				if parsedInstructionMeta.Decimals == 0 {
					parsedInstructionMeta.Decimals = Decimals
				}
				if parsedInstructionMeta.Amount == "" {
					parsedInstructionMeta.Amount = strconv.FormatInt(int64(parsedInstructionMeta.Lamports), 10)
				}
				currency := types.Currency{
					Symbol:   parsedInstructionMeta.Mint,
					Decimals: int32(parsedInstructionMeta.Decimals),
					Metadata: map[string]interface{}{},
				}
				sender := types.AccountIdentifier{
					Address:  parsedInstructionMeta.Source,
					Metadata: map[string]interface{}{},
				}
				senderAmt := types.Amount{
					Value:    "-" + parsedInstructionMeta.Amount,
					Currency: &currency,
				}
				receiver := types.AccountIdentifier{
					Address:  parsedInstructionMeta.Destination,
					Metadata: map[string]interface{}{},
				}
				receiverAmt := types.Amount{
					Value:    parsedInstructionMeta.Amount,
					Currency: &currency,
				}
				oi2 := types.OperationIdentifier{
					Index: opIndex,
				}
				opIndex += 1

				//for construction test
				delete(inInterface, "amount")
				delete(inInterface, "lamports")
				delete(inInterface, "source")
				delete(inInterface, "destination")

				//sender push
				operations = append(operations, &types.Operation{
					OperationIdentifier: &oi,
					Type:                opType,
					Status:              &status,
					Account:             &sender,
					Amount:              &senderAmt,
					Metadata:            inInterface,
				}, &types.Operation{
					OperationIdentifier: &oi2,
					Type:                opType,
					Status:              &status,
					Account:             &receiver,
					Amount:              &receiverAmt,
					Metadata:            inInterface,
				})
			} else {
				operations = append(operations, &types.Operation{
					OperationIdentifier: &oi,
					Type:                opType,
					Status:              &status,
					Metadata:            inInterface,
				})
			}
		}
	}
	return operations
}
func programFromId(programId string) string {
	program := "unknown"
	switch programId {
	case system.PROGRAM_ID.String():
		program = "system"
		break
	case token.TOKEN_PROGRAM_ID.String():
		program = "spl-token"
		break
	}
	return program
}

func ToRosTxs(txs []dfuserpc.TransactionParsed) []*RosettaTypes.Transaction {
	var rtxs []*RosettaTypes.Transaction
	for _, tx := range txs {
		rtx := ToRosTx(tx)
		rtxs = append(rtxs, &rtx)
	}
	return rtxs
}
func ToRosTx(tx dfuserpc.TransactionParsed) RosettaTypes.Transaction {
	return RosettaTypes.Transaction{
		TransactionIdentifier: &RosettaTypes.TransactionIdentifier{
			Hash: tx.Transaction.Signatures[0].String(),
		},
		Operations: GetRosOperationsFromTx(tx, SuccessStatus),
		Metadata:   map[string]interface{}{},
	}
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
