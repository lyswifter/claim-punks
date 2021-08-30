package main

import (
	"context"

	statemachine "github.com/filecoin-project/go-statemachine"
)

type TaskInfo struct {
	ID      uint64
	Wallet  string
	State   TaskState
	IP      string
	TokenID string
}

type TaskState string

const (
	Empty      TaskState = "Empty"
	Start      TaskState = "Start"
	Processing TaskState = "Processing"
	Finished   TaskState = "Finished"
	Successed  TaskState = "Successed"
	Failed     TaskState = "Failed"
)

type ICPunks struct {
	FSGroup *statemachine.StateGroup
}

func SetupFSM() {
	icpunks := ICPunks{}

	icpunks.FSGroup = statemachine.New(DB, &icpunks, TaskInfo{})

	go icpunks.Run(context.TODO())
}

func (ic *ICPunks) Run(ctx context.Context) error {
	err := ic.restartTasks(ctx)
	if err != nil {
		return err
	}

	return nil
}
