package node

import (
	"sync"
	"strconv"
	"net/url"
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/secretanalytics/go-scrt-events/pkg/types"
)

//WsRequest gets passed to websockets write()
type WsRequest struct {
	JSONRPC string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	ID      int      `json:"id"`
}

//NewWsRequest to be passed to HandleWs in main
func NewWsRequest(endpoint string, reqParams []string) WsRequest {
	return WsRequest{JSONRPC: "2.0", Method: endpoint, Params: reqParams, ID: 1}
}

func read(c *websocket.Conn, blocksOut chan types.BlockResultDB, chainTip chan int, done chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-done:
			return
		default:
			var response types.WsResponse
		    _, message, err := c.ReadMessage()
		    if err != nil {
		    	logrus.Fatal("Failed to readmessage from websocket:", err)
		    }
		    errResp := json.Unmarshal(message, &response)
		    if errResp != nil {
		    	logrus.Fatal("Failed to unmarshall result", errResp)
		    }
		    var checkMap map[string]string
		    json.Unmarshal(response.RespResult, &checkMap)
		    _, checkBlockResult := checkMap["begin_block_events"]
		    if checkBlockResult {
		    	logrus.Debug("Response is BlockResult")
		    	var blockOut types.BlockResult
		    	errBlock := json.Unmarshal(response.RespResult, &blockOut)
		    	if errBlock != nil {
		    		logrus.Fatal("Failed to unmarshall Result to BlockResult:", errBlock)
		    	}
				outBlock := blockOut.DecodeBlock("secret-2")
				logrus.Debug("Got block at height ", outBlock.Height)

		    	blocksOut <- outBlock
		    }
		    _, checkStatus := checkMap["sync_info"]
		    if checkStatus {
		    	logrus.Debug("Response is Status")
		    	var status types.StatusResult
		    	errStatus := json.Unmarshal(response.RespResult, &status)
		    	if errStatus != nil {
		    		logrus.Fatal("Failed to unmarshall Result to StatusResult:", errStatus)
		    	}
		    	height, errHeight := strconv.Atoi(status.SyInfo.LatestBlockHeight)
		    	if errHeight != nil {
		    		logrus.Fatal("Failed to convert height string -> int:", errHeight)
		    	}
		    	logrus.Info("Block height", height)
		    	chainTip <- height
		    }
		}
	}
}

func write(c *websocket.Conn, reqs chan WsRequest, done chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case msg := <- reqs:
			err := c.WriteJSON(msg)
			if err != nil {
				logrus.Fatal("write:", err)
			}
			logrus.Debug("Request Made->", msg)
		case <-done:
			logrus.Info("Channel done closed")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logrus.Fatal("write close:", err)
			}
			logrus.Info("Websocket cleanly closed by interrupt")
			return
		}
	}
}

func mergeRequests(reqs ...<-chan WsRequest) chan WsRequest {
	var wg sync.WaitGroup
	out := make(chan WsRequest)
    output := func(c <-chan WsRequest) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }
    wg.Add(len(reqs))
    for _, c := range reqs {
        go output(c)
    }
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}

func emitBlocks(done chan struct{}, blockReqs chan WsRequest, heightsIn chan int, wg *sync.WaitGroup) {
	for {
		select {
		case <-done:
			return
		case height := <- heightsIn:
			params := []string{strconv.Itoa(height)}
		    blockReqs <- NewWsRequest("block_results", params)
		}
	}
	
	wg.Done()
}

func iterBlocks(c *websocket.Conn, done chan struct{}, heightsIn, chainTip chan int, blocksOut chan types.BlockResultDB, wg *sync.WaitGroup) {
	chainTipReq := make(chan WsRequest)
	blockReqs := make(chan WsRequest)

	reqs := mergeRequests(chainTipReq, blockReqs)

	wg.Add(1)
	go read(c, blocksOut, chainTip, done, wg)
	wg.Add(1)
	go write(c, reqs, done, wg)
	var params []string
	chainTipReq <- NewWsRequest("status", params)
	
	wg.Add(1)
	go emitBlocks(done, blockReqs, heightsIn, wg)
}

func HandleWs(done chan struct{}, host, path string, heightsIn, chainTip chan int, blocksOut chan types.BlockResultDB, wg *sync.WaitGroup) {
	u := url.URL{Scheme: "ws", Host: host, Path: path}
	logrus.Debug("connecting to", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logrus.Fatal("dial:", err)
	}

	iterBlocks(c, done, heightsIn, chainTip, blocksOut, wg)
	wg.Done()
}
