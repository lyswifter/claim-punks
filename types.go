package main

var statErr = "err"
var statOk = "ok"

type StatuePunk struct {
	IP     string `json:"ip"`
	Statue string `json:"status"`
	Type   string `json:"type"`
}
