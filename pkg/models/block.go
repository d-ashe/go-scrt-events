package models

import (
	"sync"
	"context"

	"encoding/json" 

	"github.com/sirupsen/logrus"
	"cloud.google.com/go/pubsub"
)

//BlockResult is used to unmarshall JSONRPC responses for block_results?height={n} endpoint
type BlockResult struct {
	Height                string            `json:"height"`
	BlockId               string            `json:"block_id"`
	Txs                   []Tx              `json:"txs_results"`
	BeginBlockEvents      []Event           `json:"begin_block_events"`
	EndBlockEvents        []Event           `json:"end_block_events"`
	ValidatorUpdates      json.RawMessage   `json:"validator_updates"`
	ConsensusParamUpdates json.RawMessage   `json:"consensus_param_updates"`
}

func (bl *BlockResult) Publish(ctx context.Context, topic pubsub.Topic) {
	block := bl.DecodeBlock("secret-2")
	pubBlock, _ := json.Marshal(block)
	topic.Publish(ctx, &pubsub.Message{Data: pubBlock})
}

/*
//func emitDone(done chan struct{}, blocksIn chan types.BlockResultDB, chainTip int, wg *sync.WaitGroup) {
	for {
		select {
		case block := <- blocksIn:
		    if block.Height == chainTip {
				logrus.Info("SIGNALING DONE - CHAINTIP REACHED")
				close(done)
			}
		}
	}
	wg.Done()
}
///
///
*/
func PublishBlocks(projectID, topicName string, heights int[]) {
	var wg sync.WaitGroup
	dataIn := make(chan json.RawMessage)
	pubBlocks := make(chan json.RawMessage)
	done := make(chan struct{})
	//ctx context.Context, topic pubsub.Topic, done chan struct{}, dataIn chan json.RawMessage, wg *sync.WaitGroup
	ctx, _, topic := rwe.InitTopic(projectID, topicName)
	wg.Add(1)
	go web.IterBlocks(done, heights, pubBlocks, &wg)
	wg.Wait(1)
	//done chan struct{}, c *websocket.Conn, heightsIn []int, blocksOut chan json.RawMessage, wg *sync.WaitGroup
	emitDone := func() {
		defer wg.Done()
		for x := range pubBlocks {
			var block BlockResult
			json.Unmarshal(x, &block)
		}
		close(done)
	}
	wg.Add(1)
	go emitDone()
	wg.Wait()
}

//DecodeBlock() converts WsResponse BlockResult to BlockResultDB 
//----- 
//Base64 decodes event attributes, converts, adds fields for DB insert
//Decodes Txs
//Decodes BeginBlockEvents
//Decodes EndBlockEvents
//Converts height from string to int
//Adds chainId param to returned BlockResultDB struct
func (bl *BlockResult) DecodeBlock(chainId string) (dbBlock BlockResult){
	blockId := chainId + "-" + bl.Height
	dbBlock.Txs = decodeTxs(bl.Txs, blockId)
	if len(bl.BeginBlockEvents) != 0 {
		logrus.Debug("Decoding Begin Block Events")
		bl.BeginBlockEvents = decodeEventList(bl.BeginBlockEvents)
	}
	dbBlock.BeginBlockEvents = bl.BeginBlockEvents

	if len(bl.EndBlockEvents) != 0 {
		logrus.Debug("Decoding End Block Events")
		dbBlock.EndBlockEvents = decodeEventList(bl.EndBlockEvents)
	}

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