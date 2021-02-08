package models

import (
	"encoding/json" 

	"github.com/sirupsen/logrus"
    "github.com/dgraph-io/dgo"
    "github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

type Node interface {
	Upsert()
}

//Upsert structs to Dgraph
func (n Node) Upsert() {
    n.mutateUpsert()
}

func (msg json.RawMessage) mutateUpsert() {
    mu := &api.Mutation{
	  SetJson: msg,
	  CommitNow: true,
	}
    assigned, err := txn.Mutate(ctx, mu)
    if err != nil {
      logrus.Fatal(err)
    }
}

type Address struct {
	//ID is HumanAddr
	ID           string         `json:"id"`
	IsContract   bool           `json:"is_contract"`
	Signer       Address        `json:"signer,omitempty"`
	CodeID       string         `json:"code_id,omitempty"`
	Views        []ViewingKey   `json:"views,omitempty"`
}

type ViewingKey struct {
	ID     string   `json:"id"`
	Key    string   `json:"key"`
	Signer Address  `json:"signer"`
}

func (addr Address) Upsert() {
    msg, err := json.Marshal(addr)
    if err != nil {
      logrus.Fatal(err)
    }
    mutateUpsert(msg)
}
