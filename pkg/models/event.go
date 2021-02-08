package models

import (
	"encoding/json" 

	"github.com/sirupsen/logrus"
)

//Event is used to unmarshall JSONRPC responses
type Event struct {
	Id          int        `json:"id"`
	BlockId     string        `json:"block_id"`
	TxId        int        `json:"tx_id"`
	Type        string     `json:"type"`
	Attributes []Attribute `json:"attributes"`
}

func (ev Event) Upsert() {
    msg, err := json.Marshal(ev)
    if err != nil {
      logrus.Fatal(err)
    }
    mutateUpsert(msg)
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

func (ev *Event) decodeEvent(blockIdIn string) (evOut Event){
	var attrsOut []Attribute
	for _, x := range ev.Attributes {
        attrsOut = append(attrsOut, x.decodeAttribute())
	}
	evOut.Attributes = attrsOut
	evOut.Type = ev.Type
	evOut.BlockId = blockIdIn
	return evOut
}

func decodeEventList(encEvents []Event, blockIdIn string) (decEvents []Event) {
	for _, x := range encEvents {
        decEvents = append(decEvents, x.decodeEvent(blockIdIn))
	}
	return decEvents
}
