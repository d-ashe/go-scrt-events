package models

//StatusResult is used to unmarshall JSONRPC responses for status? endpoint
type StatusResult struct {
	NInfo   NodeInfo      `json:"node_info"`
	SyInfo  SyncInfo      `json:"sync_info"`
	ValInfo ValidatorInfo `json:"validator_info"`
}

type ValidatorInfo struct {
	Address string `json:"address"`
	PubKey PubKeyStatus `json:"pub_key"`
	VotingPower string `json:"voting_power"`
}

type PubKeyStatus struct {
	PType string `json:"type"`
	PValue string `json:"value"`
}

type NodeInfo struct {
	PVersion ProtocolVersion `json:"protocol_version"`
	ID string `json:"id"`
	ListenAddr string `json:"listen_addr"`
	Network string `json:"network"`
	Version string `json:"version"`
	Channels string `json:"channels"`
	Moniker string `json:"moniker"`
	OtherInfo OtherStatus `json:"other"`
}

type OtherStatus struct {
	TxIndex string `json:"tx_index"`
	RpcAddr string `json:"rpc_address"`
}

type ProtocolVersion struct {
	P2P string `json:"p2p"`
	Block string `json:"block"`
	App   string `json:"app"`
}

type SyncInfo struct {
	LatestBlockHash string `json:"latest_block_hash"`
	LatestAppHash   string `json:"latest_app_hash"`
	LatestBlockHeight string `json:"latest_block_height"`
	LatestBlockTime string `json:"latest_block_time"`
	EarliestBlockHash string `json:"earliest_block_hash"`
	EarliestAppHash   string `json:"earliest_app_hash"`
	EarliestBlockHeight string `json:"earliest_block_height"`
	EarliestBlockTime string `json:"earliest_block_time"`
	CatchingUp bool `json:"catching_up"`
}