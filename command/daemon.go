package command

// import (
// 	logging "github.com/ipfs/go-log/v2"
// 	"github.com/urfave/cli"
// )

// var log = logging.Logger("command/storetheindex")

// var DaemonCmd = &cli.Command{
// 	Name:  "daemon",
// 	Usage: "Start an daemon",
// 	Flags: []cli.Flag{},
// 	Action: func(cctx *cli.Context) {
// 		log.Info("Starting daemon servers")
// 		var err error

// 		select {
// 		case <-cctx.Done():
// 			// Command was canceled (ctrl-c)
// 		}

// 		log.Infow("Shutting down daemon")
// 	},
// }
