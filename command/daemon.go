package command

import (
	"log"

	cli "github.com/urfave/cli/v2"
)

var DaemonCmd = &cli.Command{
	Name:  "daemon",
	Usage: "Start an indexer daemon, accepting http requests",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "disablep2p",
			Usage:    "Disable libp2p client api for indexer",
			Value:    false,
			Required: false,
		},
	},
	Action: func(cctx *cli.Context) error {

		log.Printf("Version: v%s", Version)

		DataStores()

		ip, err := GetClientIp()
		if err != nil {
			return err
		}
		ClientIP = ip

		ExClientIP = GetExternalIp()

		log.Printf("ExClientIP: %s ClientIP: %s", ExClientIP, ClientIP)

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

		err = tigger(false)
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
