package solanago

import (
	"fmt"
	"testing"

	"github.com/test-go/testify/assert"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestTxParse(t *testing.T) {
	tx, err := GetTxFromStr("64dq82ETBCJ9zzS6cUGqKc8L8bZ2ZTao3wR2nARKFqBywDccMta29VgGVNK2oza3nhoqidoUZczgyfNgmoTuYrdro3UXdwwVh5TVMx2CUzFUGUmmRmsaqJ1QnFxHQCUzhbroCddPPfvjw9edG3v1aetyNRknQtxgjXEjzkgn9EGtY3mo5XoRiw38qmwqACNkdsKqfNCcG5SC9mujtCoLaFXcmnVeAcdLMgBxXsTjv1JtiLpaWsB5g7TcEo2hLHL8sLV7ZiVsn66xA1ZBdAcFsLu572CHKQ8JJkgkX")
	ptx, err := ToParsedTransaction(tx)
	fmt.Println(len(tx.Message.Instructions))
	fmt.Println(ptx.Message.Instructions[0].Data)
	assert.NoError(t, err)
}
