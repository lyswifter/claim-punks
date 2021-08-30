package main

import (
	"fmt"
	"os"

	"github.com/lyswifter/claimpunks/fsm"
	gen "github.com/whyrusleeping/cbor-gen"
)

func main() {
	err := gen.WriteMapEncodersToFile("./fsm/cbor_gen.go", "fsm",
		fsm.TaskInfo{},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
