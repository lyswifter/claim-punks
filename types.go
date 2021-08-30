package main

var statErr = "err"
var statOk = "ok"

type StatuePunk struct {
	IP     string `json:"ip"`
	Statue string `json:"status"`
	Type   string `json:"type"`
}

type PemPunk struct {
	IP        string `json:"ip"`
	Wal       string `json:"wal"`
	WalPriKey string `json:"walPriKey"`
}

type ResultPunk struct {
	IP        string `json:"ip"`
	Wal       string `json:"wal"`
	TokenID   string `json:"tokenID"`
	WalPubkey string `json:"walPubKey"`
	WalPriKey string `json:"walPriKey"`
}

type punks struct {
	Wal  string
	Data string
}

type remaining struct {
	Err error
	Wal string
}
