package web

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

func read(done chan struct{}, c *websocket.Conn, responsesOut chan json.RawMessage, wg *sync.WaitGroup) {
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
			responsesOut <- response.RespResult
		    }
		}
	}
}

func write(done chan struct{}, c *websocket.Conn, requestsIn chan WsRequest, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case msg := <- requestsIn:
			err := c.WriteJSON(msg)
			if err != nil {
				logrus.Fatal("write:", err)
			}
			//logrus.Debug("Request Made->", msg)
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

func emitBlocks(blockReqs chan WsRequest, heightsIn []int, wg *sync.WaitGroup) {
	defer wg.Done()
	for x := range heightsIn {
		params := []string{strconv.Itoa(x)}
		blockReqs <- NewWsRequest("block_results", params)
		}
}

func iterRequests(done chan struct{}, c *websocket.Conn, requestsIn chan WsRequest, responsesOut chan json.RawMessage, wg *sync.WaitGroup) {
	wg.Add(1)
	go read(done, c, blocksOut, wg)

	wg.Add(1)
	go write(done, c, reqs, wg)
}

func iterBlocks(done chan struct{}, c *websocket.Conn, heightsIn []int, blocksOut chan json.RawMessage, wg *sync.WaitGroup) {
	blockReqs := make(chan WsRequest)

	go iterRequests(done, c, blockReqs, blocksOut, wg)

	wg.Add(1)
	go emitBlocks(blockReqs, heightsIn, wg)
}

func initWs(host, path string) *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: host, Path: path}
	logrus.Debug("connecting to", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logrus.Fatal("dial:", err)
	}
	return c
}
