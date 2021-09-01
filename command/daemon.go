package command

import (
	"fmt"
	"log"

	cli "github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

var DaemonCmd = &cli.Command{
	Name:  "daemon",
	Usage: "Start a claim-punk daemon",
	Flags: []cli.Flag{
		&cli.Int64Flag{
			Name:  "count",
			Usage: "specify repeat count",
			Value: 1000,
		},
		&cli.Int64Flag{
			Name:  "delta",
			Usage: "specify repeat time interval",
			Value: 100,
		},
	},
	Action: func(cctx *cli.Context) error {

		if cctx.Int64("count") == 0 {
			return fmt.Errorf("count must not be zero")
		}

		if cctx.Int64("delta") == 0 {
			return fmt.Errorf("delta must not be zero")
		}

		DataStores()

		eip, err := getCacheExternalIp()
		if err != nil || eip == "" {
			eip = GetExternalIp()
		}

		if eip == "" {
			return xerrors.Errorf("eip is empty")
		}

		ExClientIP = eip
		ClientIP = eip
		log.Printf("ExClientIP: %s", ExClientIP)

		err = saveExternalIp(ExClientIP)
		if err != nil {
			return err
		}

		err = prepareEnv()
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
