package types


type MsgStoreCode struct {
	Sender string `json:"sender" yaml:"sender"`
	// WASMByteCode can be raw or gzip compressed
	WASMByteCode []byte `json:"wasm_byte_code" yaml:"wasm_byte_code"`
	// Source is a valid absolute HTTPS URI to the contract's source code, optional
	Source string `json:"source" yaml:"source"`
	// Builder is a valid docker image name with tag, optional
	Builder string `json:"builder" yaml:"builder"`
	// InstantiatePermission to apply on contract creation, optional
	// InstantiatePermission *AccessConfig `json:"instantiate_permission,omitempty" yaml:"instantiate_permission"`
}

type MsgInstantiateContract struct {
	Sender string `json:"sender" yaml:"sender"`
	// Admin is an optional address that can execute migrations
	// Admin string `json:"admin,omitempty" yaml:"admin"`
	// This field is only used for callbacks constructed with this message type
	CallbackCodeHash  string    `json:"callback_code_hash" yaml:"callback_code_hash"`
	CodeID            uint64    `json:"code_id" yaml:"code_id"`
	Label             string    `json:"label" yaml:"label"`
	InitMsg           []byte    `json:"init_msg" yaml:"init_msg"`
	InitFunds         string `json:"init_funds" yaml:"init_funds"`
	CallbackSignature []byte    `json:"callback_sig" yaml:"callback_sig"` // Optional
}

type MsgExecuteContract struct {
	Sender            string `json:"sender" yaml:"sender"`
	Contract          string `json:"contract" yaml:"contract"`
	Msg               []byte         `json:"msg" yaml:"msg"`
	CallbackCodeHash  string         `json:"callback_code_hash" yaml:"callback_code_hash"`
	SentFunds         string      `json:"sent_funds" yaml:"sent_funds"`
	CallbackSignature []byte         `json:"callback_sig" yaml:"callback_sig"` // Optional
}