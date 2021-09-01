package command

import (
	"fmt"

	cli "github.com/urfave/cli/v2"
)

var TiggerCmd = &cli.Command{
	Name:  "tigger",
	Usage: "Start an indexer daemon, accepting http requests",
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

		// tigger

		err := tigger(true, cctx.Int64("count"), int(cctx.Int64("delta")))
		if err != nil {
			return err
		}

		return nil
	},
}
