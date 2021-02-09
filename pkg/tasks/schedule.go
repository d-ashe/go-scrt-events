package tasks

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type TaskMsg struct {
	Type    string            `json:"type"`
	Msg     json.RawMessage   `json:"heights"`
}


//
func (task *TaskMsg) HandleTask {
	switch {
		case task.Type == "gather_blocks":
			//
			var cmd GatherBlocksMsg
			errResp := json.Unmarshal(message, &cmd)
			//Call Gather blocks tasks
			done := make(chan struct{})
			blocksOut := make(chan json.RawMessage)
			var wg sync.WaitGroup
			IterBlocks(done, cmd.Heights, blocksOut, wg)
			switch {
			case cmd.Output == "pg":
				logrus.Info("Outputting to Postgresql")
		    	//
			case cmd.Output == "dg":
				logrus.Info("Outputting to Dgraph")
		        //
			case cmd.Output == "gcppubsub":
				logrus.Info("Outputting to gcppubsub")
				//
				models.PublishBlocks(viper.Get("pubsub_project_id"), viper.Get("pubsub_topic_name"), heights)
			}
		case task.Type == "distribute_gather_blocks":
			var cmd DistributeGatherBlocksMsg
			errResp := json.Unmarshal(message, &cmd)
			//Call EmitBlocks to be gathered
		//Can determine which output to use by task.Type
		//
		/*
		
		*/
	}
}

//GatherBlocksMsg uses models.PublishBlocks to gather block results via websocket. 
type GatherBlocksMsg {
	Output  string  `json:"output"`
	Heights int[]   `json:"heights"`
	ChainId string  `json:"chain_id"`
}

//DistributeGatherBlocksMsg has necessary parameters to determine heights necessary to bring given chain-id
// latest data in database. Heights necessary are distributed via GatherBlocksMsg
type DistributeGatherBlocksMsg {
	ChainTip  int  `json:"chain_tip"`
	ChainId string  `json:"chain_id"`
}

