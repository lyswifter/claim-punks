package command

import (
	"fmt"
	"os/exec"

	"github.com/mitchellh/go-homedir"
)

var gdfx = ""
var home = ""

var wallets = []string{}

var ClientIP string = ""
var ExClientIP string = ""

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

	dfx, err := exec.LookPath("dfx")
	if err != nil {
		return
	}
	gdfx = dfx
}
