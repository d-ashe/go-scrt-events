package rwe

import (
	"sync"
	"context"

	"github.com/spf13/viper"
	"github.com/sirupsen/logrus"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/secretanalytics/go-scrt-events/pkg/models"
)


type DBConfig struct {
	Conn string
}

func InitDB() *pg.DB {
	//Parse the connection string
	opt, err := pg.ParseURL(viper.Get("db_conn"))
    if err != nil {
       panic(err)
	}
	//Connect to the db
	db := pg.Connect(opt)
	
	
	ctx := context.Background()

	//Check if db is up
	if errUp := db.Ping(ctx); errUp != nil {
		logrus.Fatal(errUp)
		logrus.Fatal("Failed to ping database")
		panic(errUp)
	}

	//Create the schema
	errSchema := createSchema(db)
    if errSchema != nil {
		logrus.Fatal(errSchema)
		logrus.Fatal("Failed to create schema")
        panic(errSchema)
	}

	//Query version
	var version string
    _, errVer := db.QueryOneContext(ctx, pg.Scan(&version), "SELECT version()")
    if errVer != nil {
		logrus.Fatal(errVer)
		logrus.Fatal("Failed to check version")
		panic(errVer)
	}
	
	logrus.Info("Connected to Postgresql with version: ", version)

	createSchema(db)
	return db
}



func InsertBlock(db *pg.DB, block *types.BlockResultDB) {
	_, err := db.Model(block).Insert()
	if err != nil {
		logrus.Fatal("Failed to insert block: ", err)
		return
	}
	logrus.Debug("Successful insert of block height", block.Height)
	if len(block.Txs) != 0 {
		insertTxs(db, block)
	}
	insertEvents(db, block)
}

func insertTxs(db *pg.DB, block *types.BlockResultDB) {
	for _, tx := range block.Txs {
		tx.BlockId = block.ID
		_, errTx := db.Model(&tx).Insert()
		if errTx != nil {
			logrus.Fatal("Failed to insert Tx: ", errTx)
			return
		}
		for _, ev := range tx.Events {
			ev.TxId = tx.ID
			ev.BlockId = block.ID
			_, errEv := db.Model(&ev).Insert()
		    if errEv != nil {
		    	logrus.Fatal("Failed to insert Event: ", errEv)
		    	return
		    }
		}
	}
}

func insertEvents(db *pg.DB, block *types.BlockResultDB) {
	insertEventList := func(eventsList []types.Event) {
		if len(eventsList) != 0 {
			for _, ev := range eventsList {
				ev.BlockId = block.ID
				_, errEv := db.Model(&ev).Insert()
				if errEv != nil {
					logrus.Fatal("Failed to insert Event: ", errEv)
				}
			}
		}
	}
	insertEventList(block.BeginBlockEvents)
	insertEventList(block.EndBlockEvents)
}

//InsertBlocks receives blocks via blocks channel from HandleWs websocket.
//
func InsertBlocks(done chan struct{}, db *pg.DB, blocksIn chan types.BlockResultDB, blocksOut chan types.BlockResultDB, wg *sync.WaitGroup) {
	defer db.Close()
	for {
		select {
		case <-done:
			return
		case block := <- blocksIn:
			insertBlock(db, &block)
			blocksOut <- block
		}
	}
	wg.Done()
}

//GetDBTip returns the max block height for a given chain-id
//func GetDBTip(db *pg.DB) int {
//	err := db.Model((*types.BlockResult)(nil)).
//		Column("height").
//		ColumnExpr("max(height) AS max_height").
//		Order("height DESC").
//		Select(res)
//}

// createSchema creates database schema for Block, Tx, and Msg models.
func createSchema(db *pg.DB) error {
    models := []interface{}{
		(*types.BlockResultDB)(nil),
		(*types.Tx)(nil),
		(*types.Event)(nil),
    }

    for _, model := range models {
        db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: false,
			IfNotExists: true,
        })
    }
    return nil
}

//QueryHeights
//Params
//chainId string -> Chain id to query for 
//Returns
//heights []int -> All heights for chainId
func QueryHeights(db *pg.DB, chainId string) []int {
	var heights []int
	err := db.Model((*types.BlockResultDB)(nil)).
	Column("height").
	Where("chain_id = ?", chainId).
	Select(&heights)
	if err != nil {
		logrus.Fatal("Failed to get heights: ", err)
		return heights
	}
	logrus.Info(heights)
	return heights
}