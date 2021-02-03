package types


type MsgSend struct {
	Sender string `json:"sender"`
	Recipient string `json:"recipient"`
	// In uscrt
	amount string `json:"amount"`
}

