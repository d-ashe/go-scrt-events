package cmd

import (
	"fmt"
	"io"
	"os"


	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	c "github.com/secretanalytics/go-scrt-events/config"

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
				logrus.Error("Unable to decode into config struct, %v", err)
			}
			var testHeights []int
			for i := 0; i < 100; i++ {
				append(testHeights, i)
			}
			models.PublishBlocks(testHeights)
		},
	}
)

func ScrtEventsCmd() *cobra.Command {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := setUpLogs(os.Stdout, v); err != nil {
			return err
		}
		return nil
	}
	cobra.OnInitialize(initConfig)
	return rootCmd
}

func setUpLogs(out io.Writer) error {
	logrus.SetOutput(out)
	lvl, err := logrus.ParseLevel("debug")
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	return nil
}

func bindEnvVarToConf(conf, envVar string) {
	err := viper.BindEnv(conf, envVar)
    if err != nil {
		log := "Failed to bind" + envVar + conf
		logrus.Fatal(log)
		logrus.Fatal(err)	
	}
}

func initConfig() {
	viper.AutomaticEnv()

	bindEnvVarToConf("node.host", "NODE_HOST")
	bindEnvVarToConf("node.path", "WS_PATH")

	bindEnvVarToConf("database.conn", "DB_CONN")

	bindEnvVarToConf("pubsub.projectid", "PUBSUB_PROJECT_ID")
	bindEnvVarToConf("pubsub.topicname", "PUBSUB_TOPIC_NAME")


	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
	}
}