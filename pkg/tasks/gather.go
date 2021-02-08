package tasks

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/d-ashe/Seraph/pkg/rwe"
	"github.com/d-ashe/Seraph/pkg/web"
)

//GetStatus requests status? endpoint via websocket, 
//returns chainId, and chainTip
//
func GetStatus() (string, int, error) {
	//chainId, chainTip := 
	//return chainId, chainTip, nil
}

//GetHeights returns all block heights for given chainId
//
//
func GetHeights(chainId string) (int[], error){
	dbSession := rwe.InitDB(viper.Get("db_conn"))

	heights := rwe.QueryHeights(dbSession, chainId)
	return heights
}

//GatherBlocks gathers blockresults via websocket
//These block results are then inserted into Postgresql
func GatherBlocks(heightsIn []int) error {
	var wg sync.WaitGroup

	done := make(chan struct{})
	blocksOutWeb := make(chan types.WsResponse)
	blocksInDB := make(chan types.BlockResultDB)

	dbSession := rwe.InitDB(viper.Get("db_conn"))

	signalDone := func () {
		for x := range blocksInDB {
			inBlock := x.DecodeBlock()
			crud.InsertBlock(dbSession, inBlock)
			blocksInDB <- inBlock
		}
		wg.Done()
		close(done)
	}

	wg.Add(1)
	go GetBlocks(done, viper.Get("node_host"), viper.Get("node_ws_path"), heightsIn,  wg)
	wg.Add(1)
	go signalDone()
	return nil
}