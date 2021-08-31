package command

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
)

const Version = "0.0.3"

var gdfx = ""
var home = ""
var wallets = []string{}
var ClientIP string = ""
var PunksInHand = []ResultPunk{}
var output chan punks = make(chan punks, 10)
var remainChan chan remaining = make(chan remaining, 10)

func init() {
	for i := 0; i < wallet_count; i++ {
		wallets = append(wallets, fmt.Sprintf("%d", i))
	}

	h, err := homedir.Dir()
	if err != nil {
		return
	}
	home = h
}
