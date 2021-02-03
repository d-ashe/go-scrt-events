package types

import (
	b64 "encoding/base64"
	"strconv"
	"encoding/json" 
	"github.com/sirupsen/logrus"
)

//WsResponse is used to unmarshall JSONRPC responses
type WsResponse struct {
	JSONRPC string   `json:"jsonrpc"`
	ID int   `json:"id"`
	RespResult json.RawMessage `json:"result"`
}

//StatusResult is used to unmarshall JSONRPC responses for status? endpoint
type StatusResult struct {
	NInfo   NodeInfo      `json:"node_info"`
	SyInfo  SyncInfo      `json:"sync_info"`
	ValInfo ValidatorInfo `json:"validator_info"`
}

type ValidatorInfo struct {
	Address string `json:"address"`
	PubKey PubKeyStatus `json:"pub_key"`
	VotingPower string `json:"voting_power"`
}

type PubKeyStatus struct {
	PType string `json:"type"`
	PValue string `json:"value"`
}

type NodeInfo struct {
	PVersion ProtocolVersion `json:"protocol_version"`
	ID string `json:"id"`
	ListenAddr string `json:"listen_addr"`
	Network string `json:"network"`
	Version string `json:"version"`
	Channels string `json:"channels"`
	Moniker string `json:"moniker"`
	OtherInfo OtherStatus `json:"other"`
}

type OtherStatus struct {
	TxIndex string `json:"tx_index"`
	RpcAddr string `json:"rpc_address"`
}

type ProtocolVersion struct {
	P2P string `json:"p2p"`
	Block string `json:"block"`
	App   string `json:"app"`
}

type SyncInfo struct {
	LatestBlockHash string `json:"latest_block_hash"`
	LatestAppHash   string `json:"latest_app_hash"`
	LatestBlockHeight string `json:"latest_block_height"`
	LatestBlockTime string `json:"latest_block_time"`
	EarliestBlockHash string `json:"earliest_block_hash"`
	EarliestAppHash   string `json:"earliest_app_hash"`
	EarliestBlockHeight string `json:"earliest_block_height"`
	EarliestBlockTime string `json:"earliest_block_time"`
	CatchingUp bool `json:"catching_up"`

}

//BlockResult is used to unmarshall JSONRPC responses for block_results?height={n} endpoint
type BlockResult struct {
	Height                string            `json:"height"`
	Txs                   []Tx              `json:"txs_results"`
	BeginBlockEvents      []Event           `json:"begin_block_events"`
	EndBlockEvents        []Event           `json:"end_block_events"`
	ValidatorUpdates      []ValidatorUpdate `json:"validator_updates"`
	ConsensusParamUpdates json.RawMessage 
}

//BlockResult is used to unmarshall JSONRPC responses
type BlockResultDB struct {
	tableName struct{} `pg:"blocks,alias:block"`

	ID                    int    `pg:",pk"`
	ChainId               string

	Height                int            
	Txs                   []Tx              `pg:"rel:has-many,join_fk:id"`
	BeginBlockEvents      []Event           `pg:"rel:has-many,join_fk:id"`
	EndBlockEvents        []Event           `pg:"rel:has-many,join_fk:id"`
	ValidatorUpdates      []ValidatorUpdate `pg:"rel:has-many,join_fk:id"`
	ConsensusParamUpdates json.RawMessage 
}

//Tx is used to unmarshall JSONRPC responses
type Tx struct {
	tableName struct{} `pg:"txs,alias:tx"`

	ID        int    `pg:",pk"`

	Code      int     `json:"code"`
	CodeSpace string  `json:"codespace"`
	Info      string  `json:"info"`
	Data      string  `json:"data"`
	GasWanted string  `json:"gasWanted"`
	GasUsed   string  `json:"gasUsed"`
	Log       string  `json:"log"`
	Events    []Event `json:"events" pg:"rel:has-many,join_fk:id"`
}

//ValidatorUpdate is used to unmarshall JSONRPC responses
type ValidatorUpdate struct {
	tableName struct{} `pg:"validator_updates,alias:validator_update"`
	ID        int      `pg:",pk"`

	PubKey PublicKey `json:"pub_key"`
	Power  string    `json:"power"`
}

//PublicKey is used to unmarshall JSONRPC responses
type PublicKey struct {
	Data string `json:"data"`
}

//Event is used to unmarshall JSONRPC responses
type Event struct {
	tableName struct{} `pg:"events,alias:event"`
	ID        int      `pg:",pk"`

	Type       string      `json:"type"`
	Attributes []Attribute `json:"attributes"`
}

//Attribute is used to unmarshall JSONRPC responses
type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (attr *Attribute) decodeAttribute() (attrOut Attribute){
	decKey, err1 := b64.StdEncoding.DecodeString(attr.Key)
	if err1 != nil {
		logrus.Fatal("read:", err1)
	} else {
		attrOut.Key = string(decKey[:])
	}
	decValue, err2 := b64.StdEncoding.DecodeString(attr.Value)
	if err2 != nil {
		logrus.Fatal("read:", err2)
	} else {
		attrOut.Value = string(decValue[:])
	}
	return attrOut
}

func (ev *Event) decodeEventAttributes() (evOut Event){
	var attrsOut []Attribute
	for _, x := range ev.Attributes {
        attrsOut = append(attrsOut, x.decodeAttribute())
	}
	evOut.Attributes = attrsOut
	evOut.Type = ev.Type
	return evOut
}

func decodeEventList(encEvents []Event) (decEvents []Event) {
	for _, x := range encEvents {
        decEvents = append(decEvents, x.decodeEventAttributes())
	}
	return decEvents
}

func (tx *Tx) decodeTx() {
	if len(tx.Events) != 0 {
		decEvents := decodeEventList(tx.Events)
		tx.Events = decEvents
	}
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
	if len(bl.Txs) != 0 {
		logrus.Debug("Decoding Txs")
		for _, x := range bl.Txs {
			x.decodeTx()
		}
	}
	dbBlock.Txs = bl.Txs
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