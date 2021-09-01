package command

import (
	"fmt"
	"log"

	cli "github.com/urfave/cli/v2"
)

var DaemonCmd = &cli.Command{
	Name:  "daemon",
	Usage: "Start an claim-punk daemon",
	Flags: []cli.Flag{
		&cli.Int64Flag{
			Name:  "count",
			Usage: "specify repeat count",
			Value: 100,
		},
		&cli.Int64Flag{
			Name:  "delta",
			Usage: "specify repeat time interval",
			Value: 50,
		},
	},
	Action: func(cctx *cli.Context) error {

		if cctx.Int64("count") == 0 {
			return fmt.Errorf("count must not be zero")
		}

		if cctx.Int64("delta") == 0 {
			return fmt.Errorf("delta must not be zero")
		}

		log.Printf("Version: v%s", Version)

		DataStores()

		ExClientIP = GetExternalIp()
		ClientIP = ExClientIP
		log.Printf("ExClientIP: %s", ExClientIP)

		err := prepareEnv()
		if err != nil {
			return err
		}

		go func() {
			err := reportStatus("prepareEnv finished", statOk)
			if err != nil {
				return
			}
		}()

		err = prepareProject()
		if err != nil {
			return err
		}

		go func() {
			err := reportStatus("prepareProject finished", statOk)
			if err != nil {
				return
			}
		}()

		err = prepareIdentities(wallets)
		if err != nil {
			return err
		}

		go func() {
			err := reportStatus("prepareIdentities finished", statOk)
			if err != nil {
				return
			}
		}()

		err = tigger(false, cctx.Int64("count"), int(cctx.Int64("delta")))
		if err != nil {
			return err
		}

		go func() {
			err := reportStatus("work done finished", statOk)
			if err != nil {
				return
			}
		}()

		select {
		case <-cctx.Done():
			// Command was canceled (ctrl-c)
		}

		log.Print("Shutting down daemon")
		return nil
	},
}
