package solanago

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/coinbase/rosetta-sdk-go/types"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/mr-tron/base58"
	"github.com/portto/solana-go-sdk/assotokenprog"
	ss "github.com/portto/solana-go-sdk/client"
	common "github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/sysprog"
	"github.com/portto/solana-go-sdk/tokenprog"
	solPTypes "github.com/portto/solana-go-sdk/types"

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
	return input[0:1], input[1:]
}
func GetRosOperationsFromTx(tx solPTypes.ParsedTransaction, status string) []*types.Operation {
	//	hash := tx.Transaction.Signatures[0].String()
	opIndex := int64(0)
	var operations []*types.Operation
	for _, ins := range tx.Message.Instructions {
		oi := types.OperationIdentifier{
			Index: opIndex,
		}
		opIndex += 1

		if ins.Parsed == nil {

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

			if IsBalanceChanging(opType) {
				if parsedInstructionMeta.Decimals == 0 {
					parsedInstructionMeta.Decimals = Decimals
				}
				if parsedInstructionMeta.Amount == 0 {
					if parsedInstructionMeta.Lamports == 0 {
						parsedInstructionMeta.Amount, _ = strconv.ParseUint(parsedInstructionMeta.TokenAmount.Amount, 10, 64)
					} else {
						parsedInstructionMeta.Amount = parsedInstructionMeta.Lamports
					}
				}
				var currency types.Currency
				if parsedInstructionMeta.Mint == "" {
					if ins.Program == "system" {
						currency = types.Currency{
							Symbol:   Symbol,
							Decimals: Decimals,
							Metadata: map[string]interface{}{},
						}
					}
				} else {
					currency = types.Currency{
						Symbol:   parsedInstructionMeta.Mint,
						Decimals: int32(parsedInstructionMeta.Decimals),
						Metadata: map[string]interface{}{},
					}
				}

				sender := types.AccountIdentifier{
					Address:  parsedInstructionMeta.Source,
					Metadata: map[string]interface{}{},
				}
				senderAmt := types.Amount{
					Value:    "-" + fmt.Sprint(parsedInstructionMeta.Amount),
					Currency: &currency,
				}
				receiver := types.AccountIdentifier{
					Address:  parsedInstructionMeta.Destination,
					Metadata: map[string]interface{}{},
				}
				receiverAmt := types.Amount{
					Value:    fmt.Sprint(parsedInstructionMeta.Amount),
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
				var account types.AccountIdentifier
				if parsedInstructionMeta.Source != "" {
					account = types.AccountIdentifier{
						Address: parsedInstructionMeta.Source,
					}
				}
				operations = append(operations, &types.Operation{
					OperationIdentifier: &oi,
					Type:                opType,
					Account:             &account,
					Status:              &status,
					Metadata:            inInterface,
				})
			}
		}
	}
	return operations
}

func ToRosTxs(txs []ss.ParsedTransactionWithMeta) []*RosettaTypes.Transaction {
	var rtxs []*RosettaTypes.Transaction
	for _, tx := range txs {
		rtx := ToRosTx(tx.Transaction)
		rtxs = append(rtxs, &rtx)
	}
	return rtxs
}
func ToRosTx(tx solPTypes.ParsedTransaction) RosettaTypes.Transaction {
	return RosettaTypes.Transaction{
		TransactionIdentifier: &RosettaTypes.TransactionIdentifier{
			Hash: tx.Signatures[0],
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
func EncodeBig(bigint *big.Int) string {
	nbits := bigint.BitLen()
	if nbits == 0 {
		return "0x0"
	}
	return fmt.Sprintf("%#x", bigint)
}

func convertTime(time uint64) int64 {
	return int64(time) * 1000
}

func GetWithNonce(m map[string]interface{}) (WithNonce, bool) {
	var withNonce WithNonce
	hasNonce := false
	if w, ok := m[WithNonceKey]; ok {
		j, _ := json.Marshal(w)
		json.Unmarshal(j, &withNonce)
		if len(withNonce.Account) > 0 {
			hasNonce = true
		}
	}
	return withNonce, hasNonce
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
func ToParsedTransaction(tx solPTypes.Transaction) (solPTypes.ParsedTransaction, error) {
	ins := tx.Message.DecompileInstructions()
	var parsedIns []solPTypes.ParsedInstruction
	for _, v := range ins {
		p, err := ParseInstruction(v)
		if err != nil {
			//cannot parse
			p = solPTypes.ParsedInstruction{}
			return solPTypes.ParsedTransaction{}, fmt.Errorf("Cannot parse Instruction")
		}
		parsedIns = append(parsedIns, p)
	}
	var acckeys []string
	var sigs []string
	for _, v := range tx.Message.Accounts {
		acckeys = append(acckeys, v.ToBase58())
	}
	for _, v := range tx.Signatures {
		sigs = append(sigs, v.ToBase58())
	}
	newTx := solPTypes.ParsedTransaction{
		Signatures: sigs,
		Message: solPTypes.ParsedMessage{
			Header:          tx.Message.Header,
			AccountKeys:     acckeys,
			RecentBlockhash: tx.Message.RecentBlockHash,
			Instructions:    parsedIns,
		},
	}
	return newTx, nil
}
func ParseInstruction(ins solPTypes.Instruction) (solPTypes.ParsedInstruction, error) {
	var parsedInstruction solPTypes.ParsedInstruction
	var err error

	switch ins.ProgramID {
	case common.SystemProgramID:
		parsedInstruction, err = sysprog.ParseSystem(ins)
		break
	case common.TokenProgramID:
		parsedInstruction, err = tokenprog.ParseToken(ins)
		break
	case common.SPLAssociatedTokenAccountProgramID:
		parsedInstruction, err = assotokenprog.ParseAssocToken(ins)
		break
	default:
		return parsedInstruction, fmt.Errorf("Cannot parse instruction")
	}
	if err != nil {
		return parsedInstruction, err
	}
	var accs []string
	for _, v := range ins.Accounts {
		accs = append(accs, v.PubKey.ToBase58())
	}
	parsedInstruction.Accounts = accs
	parsedInstruction.Data = base58.Encode(ins.Data[:])
	parsedInstruction.ProgramID = ins.ProgramID.ToBase58()
	parsedInstruction.Program = common.GetProgramName(ins.ProgramID)
	return parsedInstruction, nil
}
func ValueToBaseAmount(valueStr string) uint64 {
	var amount uint64
	valueStr = strings.Replace(valueStr, "-", "", -1)
	amt, err := strconv.ParseInt(valueStr, 10, 64)
	if err == nil {
		amount = uint64(amt)
	}
	return amount
}
