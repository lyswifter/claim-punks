package command

import (
	"bufio"
	"io"
	"os"
)

const repoPath = "~/.claimpunk"

var wallet_count = 30

// var destTiming = "2021-09-01 20:00:00"
var destTiming = "2021-08-31 07:20:00"

var dfxjson = "dfx.json"
var workDir = "Hell"
var projectName = "Punks"

var ClientsCluster map[string][]string = make(map[string][]string)

func init() {
	ClientsCluster["first"] = readline("./clients/first")
	ClientsCluster["seconds"] = readline("./clients/second")
	ClientsCluster["third"] = readline("./clients/third")
	ClientsCluster["last"] = readline("./clients/last")
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

		ret = append(ret, line)
	}

	return ret
}
