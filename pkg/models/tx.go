package models

import (
	"encoding/json" 

	"github.com/sirupsen/logrus"
)

//Tx is used to unmarshall JSONRPC responses
type Tx struct {
	ID        int     `json:"id""`
	BlockId   string  `json:block_id`
	Code      int     `json:"code"`
	CodeSpace string  `json:"codespace"`
	Info      string  `json:"info"`
	Data      string  `json:"data"`
	GasWanted string  `json:"gasWanted"`
	GasUsed   string  `json:"gasUsed"`
	Log       string  `json:"log"`
	Events    []Event `json:"events"`
}

func (tx Tx) Upsert() {
    msg, err := json.Marshal(tx)
    if err != nil {
      logrus.Fatal(err)
    }
    mutateUpsert(msg)
}

func (tx *Tx) decodeTx(blockIdIn string) Tx {
	txOut := Tx{
		BlockId: blockIdIn,
	    Code: tx.Code,
	    CodeSpace: tx.CodeSpace,
	    Info: tx.Info,
	    Data: tx.Data,
	    GasUsed: tx.GasUsed,
	    GasWanted: tx.GasWanted
		Log: tx.Log,
	}
	if len(tx.Events) != 0 {
		txOut.Events = decodeEventList(tx.Events)
	}
	return txOut
}

func decodeTxs(txs []Tx, blockId string) []Tx {
	var txsOut []Tx
	if len(txs) != 0 {
		logrus.Debug("Decoding Txs")
		for _, x := range txs {
			txsOut = append(txsOut, x.decodeTx(blockId))
		}
	}
	return txsOut
}
