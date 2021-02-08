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

func (task TaskMsg) HandleTask {
	switch {
		case task.Type == "gather_blocks":
			var cmd GatherBlocksMsg
			errResp := json.Unmarshal(message, &cmd)
			//Call Gather blocks tasks
		case task.Type == "chain_tip":
			var cmd ChainTipMsg
			errResp := json.Unmarshal(message, &cmd)
			//Call EmitBlocks to be gathered
	}
}

type GatherBlocksMsg {
	Heights int[]   `json:"heights"`
	ChainId string  `json:"chain_id"`
}

type ChainTipMsg {
	ChainTip  int  `json:"chain_tip"`
	ChainId string  `json:"chain_id"`
}