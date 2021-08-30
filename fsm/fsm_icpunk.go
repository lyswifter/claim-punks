package fsm

import (
	"context"
	"strconv"

	statemachine "github.com/filecoin-project/go-statemachine"
	"github.com/ipfs/go-datastore"
)

type TaskNumber uint64

func (s TaskNumber) String() string {
	return strconv.FormatUint(uint64(s), 10)
}

type TaskInfo struct {
	ID      TaskNumber
	Wallet  string
	State   TaskState
	IP      string
	TokenID string
}

type TaskState string

const (
	UndefinedSectorState TaskState = ""
	Empty                TaskState = "Empty"
	Start                TaskState = "Start"
	Processing           TaskState = "Processing"
	Finished             TaskState = "Finished"
	Successed            TaskState = "Successed"
	Failed               TaskState = "Failed"
)

type ICPunks struct {
	DB      datastore.Batching
	FSGroup *statemachine.StateGroup

	repopath string

	workdir    string
	projectdir string
}

func SetupFSM(db datastore.Batching, repopath string, workdir string, projectdir string) ICPunks {
	icpunks := ICPunks{}

	icpunks.FSGroup = statemachine.New(db, &icpunks, TaskInfo{})
	icpunks.repopath = repopath
	icpunks.workdir = workdir
	icpunks.projectdir = projectdir

	return icpunks
}

func (ic *ICPunks) Run(ctx context.Context) error {
	err := ic.restartTasks(ctx)
	if err != nil {
		return err
	}

	return nil
}
