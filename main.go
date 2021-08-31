package main

import (
	"context"
	"fmt"

	"github.com/lyswifter/claimpunks/fsm"
	// logging "github.com/ipfs/go-log/v2"
)

// var log = logging.Logger("claim-punks")

var wallets = []string{}
var ClientIP string = ""
var ExternalClientIP string = ""

var PunksInHand = []ResultPunk{}

func init() {
	for i := 0; i < wallet_count; i++ {
		wallets = append(wallets, fmt.Sprintf("%d", i))
	}
}

func main() {
	ctx := context.TODO()

	err := EnvThing()
	if err != nil {
		return
	}

	DataStores()
	icpunk := fsm.SetupFSM(DB, repoPath, workDir, projectName)

	icpunk.Run(ctx)

	// 	timeFormat := "2006-01-02 15:04:05"
	// 	destTime, err := time.ParseInLocation(timeFormat, destTiming, time.UTC)
	// 	if err != nil {
	// 		return
	// 	}

	// 	delay := destTime.Sub(time.Now().UTC())

	// 	log.Printf("time is not reached dest: %v nowUtc: %v now: %v delay: %v", destTime, time.Now().UTC(), time.Now(), delay)

	// 	timer := time.NewTimer(delay)
	// 	tickerOne := time.NewTicker(10 * time.Second)

	// nextStep:
	// 	for {
	// 		select {
	// 		case <-timer.C:
	// 			log.Printf("It's time now to do sth")
	// 			break nextStep
	// 		case <-tickerOne.C:
	// 			go func() error {
	// 				err := reportStatus("waiting time", statOk)
	// 				if err != nil {
	// 					return err
	// 				}
	// 				return nil
	// 			}()
	// 		}
	// 	}

	// Tigger
	icpunk.Tigger(ctx, wallets, ClientIP)
}
