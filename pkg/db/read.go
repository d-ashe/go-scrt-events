package db

import (
	"github.com/sirupsen/logrus"
	"github.com/go-pg/pg/v10"

	"github.com/secretanalytics/go-scrt-events/pkg/types"
)

//GetHeights
//Params
//chainId string -> Chain id to query for 
//Returns
//heights []int -> All heights for chainId
func GetHeights(db *pg.DB, chainId string) []int {
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