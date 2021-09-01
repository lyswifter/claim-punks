package command

import (
	cli "github.com/urfave/cli/v2"
)

var TiggerCmd = &cli.Command{
	Name:  "tigger",
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

		// tigger

		err := tigger(true)
		if err != nil {
			return err
		}

		return nil
	},
}
