package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/lyswifter/claimpunks/command"
	"github.com/urfave/cli/v2"
)

// var FirstTiming = "2021-09-01 19:55:00"
// var SecondTiming = "2021-09-01 19:58:00"
// var ThirdTiming = "2021-09-01 20:00:00"

var FirstTiming = "2021-09-01 08:40:00"
var SecondTiming = "2021-09-01 08:43:00"
var ThirdTiming = "2021-09-01 08:45:00"

func init() {
	clientsSep()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up a signal handler to cancel the context
	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
		select {
		case <-interrupt:
			cancel()
			fmt.Println("Received interrupt signal, shutting down...")
			fmt.Println("(Hit ctrl-c again to force-shutdown the daemon.)")
		case <-ctx.Done():
		}
		// Allow any forther SIGTERM or SIGING to kill process
		signal.Stop(interrupt)
	}()

	app := &cli.App{
		Name:    "claim-punk",
		Usage:   "Claim punk node",
		Version: "0.0.8",
		Commands: []*cli.Command{
			command.DaemonCmd,
			command.TiggerCmd,
		},
	}

	log.Printf("Version: v%s", app.Version)

	if err := app.RunContext(ctx, os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func clientsSep() {
	for _, s := range readline("./clients/first") {
		command.TimingMap[s] = FirstTiming
	}
	for _, s := range readline("./clients/second") {
		command.TimingMap[s] = SecondTiming
	}
	for _, s := range readline("./clients/third") {
		command.TimingMap[s] = ThirdTiming
	}

	// log.Printf("command.TimingMap: %+v", command.TimingMap)
}

func readline(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	rd := bufio.NewReader(f)

	var ret = []string{}
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}

		line = strings.Replace(line, "\n", "", -1)

		ret = append(ret, line)
	}

	return ret
}
