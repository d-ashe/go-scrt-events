package models

import (
	"encoding/json" 

	"github.com/sirupsen/logrus"
)


//BlockResult is used to unmarshall JSONRPC responses for block_results?height={n} endpoint
type BlockResult struct {
	Height                string            `json:"height"`
	Txs                   []Tx              `json:"txs_results"`
	BeginBlockEvents      []Event           `json:"begin_block_events"`
	EndBlockEvents        []Event           `json:"end_block_events"`
	ValidatorUpdates      json.RawMessage   `json:"validator_updates"`
	ConsensusParamUpdates json.RawMessage   `json:"consensus_param_updates"`
}

//DecodeBlock() converts WsResponse BlockResult to BlockResultDB 
//----- 
//Base64 decodes event attributes, converts, adds fields for DB insert
//Decodes Txs
//Decodes BeginBlockEvents
//Decodes EndBlockEvents
//Converts height from string to int
//Adds chainId param to returned BlockResultDB struct
func (bl *BlockResult) DecodeBlock(chainId string) (dbBlock BlockResultDB){
	blockId := chainId + "-" + bl.Height
	dbBlock.Txs = decodeTxs(bl.Txs, blockId)
	if len(bl.BeginBlockEvents) != 0 {
		logrus.Debug("Decoding Begin Block Events")
		bl.BeginBlockEvents = decodeEventList(bl.BeginBlockEvents)
	}
	dbBlock.BeginBlockEvents = bl.BeginBlockEvents

	if len(bl.EndBlockEvents) != 0 {
		logrus.Debug("Decoding End Block Events")
		bl.EndBlockEvents = decodeEventList(bl.EndBlockEvents)
	}
	dbBlock.EndBlockEvents = bl.EndBlockEvents

	dbBlock.ChainId = chainId
	
	height, errHeight := strconv.Atoi(bl.Height)
	if errHeight != nil {
		logrus.Fatal("Failed to decode height: string -> int", errHeight)
	} else {
		dbBlock.Height = height
	}
	dbBlock.ValidatorUpdates = bl.ValidatorUpdates
	dbBlock.ConsensusParamUpdates = bl.ConsensusParamUpdates
	return dbBlock
}