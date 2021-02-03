package db

import (
	"sync"
	"context"

	"github.com/sirupsen/logrus"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/secretanalytics/go-scrt-events/pkg/types"
)

func insertBlock(db *pg.DB, block *types.BlockResultDB) {
	_, err := db.Model(block).InsertCascade()
	if err != nil {
		logrus.Fatal("Failed to insert block: ", err)
		return
	}
	logrus.Debug("Successful insert")
}

func InsertBlocks(conn string, blocks chan types.BlockResultDB, wg *sync.WaitGroup) {
	db := initDB(conn)
	defer db.Close()
	for block := range blocks {
		insertBlock(db, &block)
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
    }

    for _, model := range models {
        db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: false,
			IfNotExists: true,
        })
    }
    return nil
}


func initDB(conn string) *pg.DB {
	//Parse the connection string
	opt, err := pg.ParseURL(conn)
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
