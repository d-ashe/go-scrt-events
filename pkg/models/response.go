//WsResponse is used to unmarshall JSONRPC responses
type WsResponse struct {
	JSONRPC string   `json:"jsonrpc"`
	ID int   `json:"id"`
	RespResult json.RawMessage `json:"result"`
}