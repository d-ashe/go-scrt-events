package cmd

import (
	"fmt"
	"io"
	"os"
	"sync"


	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/go-pg/pg/v10"


	c "github.com/secretanalytics/go-scrt-events/config"
	"github.com/secretanalytics/go-scrt-events/pkg/node"
	"github.com/secretanalytics/go-scrt-events/pkg/types"
	"github.com/secretanalytics/go-scrt-events/pkg/db"
)

var (
	// Used for flags.
	cfgFile string
	v       string
	rootCmd = &cobra.Command{
		Use:   "go-scrt-events",
		Short: "scrt-events quickly bootstraps a postgresql db with the Secret Network blockchain block-results.",
		Long:  `scrt-events quickly bootstraps a postgresql db with the Secret Network blockchain block-results.`,
		Run: func(cmd *cobra.Command, args []string) {
			var configuration c.Configurations
			err := viper.Unmarshal(&configuration)
			if err != nil {
				logrus.Error("Unable to decode into struct, %v", err)
			}
			logrus.Debug("Node host is: ", configuration.Node.Host)
			logrus.Debug("DB conn string is: ", configuration.Database.Conn)
			run(configuration.Database.Conn, configuration.Node.Host, configuration.Node.Path)
		},
	}
)

func emitDone(done chan struct{}, blocksIn chan types.BlockResultDB, chainTip int, wg *sync.WaitGroup) {
	for {
		select {
		case block := <- blocksIn:
		    if block.Height == chainTip {
				close(done)
			}
		default:
			logrus.Debug("Waiting to receive new blocks for insert")
		}
	}
	wg.Done()
}


//emitHeights() is the main request generator for block results, 
//Block heights available in db are compared to chaintip. 
//Blocks heights needed to sent in heightsIn channel to HandleWs()
func emitHeights(dbSession *pg.DB, chainTip int, heightsIn chan int, wg *sync.WaitGroup) {
	//Checks for existence of block height in slice of heights
	contains := func (checkFor int, inSlice []int) bool {
		for i := range inSlice {
			if i == checkFor {
				return true
			}
		}
		return false
	}
	//If no heights for given chain-id then start at 1, else sort ints and start loop at lowest
	heights := db.GetHeights(dbSession, "secret-2")
	//Loop from dbTip to chainTip, if height i not contained in heights, request for block_results at height i will be made
	for i := 1; i <= chainTip; i++ {
		if contains(i, heights) {
			logrus.Debug("Heights contain block ", i)
		} else {
			heightsIn <- i
			logrus.Debug("Requesting height ", i)
		}
	}
	close(heightsIn)
	wg.Done()
}


//run() is the main runner for go-scrt-events.
//
//Waitgroup of goroutines are started which:
//InsertsBlocks() to postgresql
//HandleWs() read/write to websocket
//emitHeights() shares a channel with HandleWs() to determine which block heights to request.
//emitDone() keeps track of results from websockets and postgresql, when all needed heights have been requested. Done is signaled. 
func run(dbConn, host, path string) {
	var wg sync.WaitGroup
	heightsIn := make(chan int)
	blocksOutWeb := make(chan types.BlockResultDB)
	blocksOutDB := make(chan types.BlockResultDB)


	chainTip := make(chan int)
	done := make(chan struct{})

	dbSession := db.InitDB(dbConn)
	logrus.Debug("Node host is: ", host)

	

	wg.Add(1)
	go db.InsertBlocks(done, dbSession, blocksOutWeb, blocksOutDB, &wg)

	wg.Add(1)
	go node.HandleWs(done, host, path, heightsIn, chainTip, blocksOutWeb, &wg)
	

    latestHeight := <- chainTip
	logrus.Info("Latest height is ", latestHeight)

	wg.Add(1)
	logrus.Info("Emitting heights to fetch")
	go emitHeights(dbSession, latestHeight, heightsIn, &wg)

	wg.Add(1)
	go emitDone(done, blocksOutDB, latestHeight, &wg)

	wg.Wait()

}


func ScrtEventsCmd() *cobra.Command {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := setUpLogs(os.Stdout, v); err != nil {
			return err
		}
		return nil
	}
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scrt-events/config.json)")
	rootCmd.PersistentFlags().StringVarP(&v, "verbosity", "v", logrus.WarnLevel.String(), "Log level (debug, info, warn, error, fatal, panic")

	return rootCmd
}

func setUpLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	return nil
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath("$HOME/.scrt-events")
		viper.SetConfigName("config.yml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
