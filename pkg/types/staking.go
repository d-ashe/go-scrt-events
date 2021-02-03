package types



type MsgUndelegate struct {
	DelegatorAddress string `json:"sender"`
	Validator string `json:"validator"`
	Amount string `json:"amount"`
}

type MsgDelegate struct {
	DelegatorAddress string `json:"sender"`
	Validator string `json:"validator"`
	Amount string `json:"amount"`
}

type MsgBeginRedelegate struct {
	DelegatorAddress string `json:"sender"`
	ValidatorSrc string `json:"validator_src"`
	ValidatorDst string `json:"validator_dst"`
	Amount string `json:"amount"`
}

type MsgCreateValidator struct {
	ValDescription Description `json:"description"`
	ValidatorAddress string `json:"validator"`
	DelegatorAddress string `json:"sender"`
	PubKey string `json:"pub_key"`
	Value string `json:"amount"`
	CommissionRate Commission `json:"commission"`
	MinSelfDelegation string `json:"min_self_delegation"`
}

type Description struct {
	Details string `json:"details"`
	Moniker string `json:"moniker"`
	Website string `json:"website"`
	Identity string `json:"identity"`
	SecurityContact string `json:"security_contact"`
}

type Commission struct {
	Rate string `json:"rate"`
	MaxRate string `json:"max_rate"`
	MaxChangeRate string `json:"max_change_rate"`
}

type MsgEditValidator struct {
	ValDescription Description `json:"description"`
	ValidatorAddress string `json:"validator"`
	Value string `json:"amount"`
	CommissionRate Commission `json:"commission"`
	MinSelfDelegation string `json:"min_self_delegation"`
}