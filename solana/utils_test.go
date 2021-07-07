package solanago

import (
	"testing"

	"github.com/test-go/testify/assert"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestValueToBaseAmount(t *testing.T) {
	amt := ValueToBaseAmount("1")
	assert.Equal(t, uint64(1), amt)
}
func TestTxParse(t *testing.T) {
	tx, err := GetTxFromStr("64dq82ETBCJ9zzS6cUGqKc8L8bZ2ZTao3wR2nARKFqBywDccMta29VgGVNK2oza3nhoqidoUZczgyfNgmoTuYrdro3UXdwwVh5TVMx2CUzFUGUmmRmsaqJ1QnFxHQCUzhbroCddPPfvjw9edG3v1aetyNRknQtxgjXEjzkgn9EGtY3mo5XoRiw38qmwqACNkdsKqfNCcG5SC9mujtCoLaFXcmnVeAcdLMgBxXsTjv1JtiLpaWsB5g7TcEo2hLHL8sLV7ZiVsn66xA1ZBdAcFsLu572CHKQ8JJkgkX")
	_, err = ToParsedTransaction(tx)
	assert.NoError(t, err)
}
