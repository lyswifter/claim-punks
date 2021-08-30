package main

import (
	"context"

	"github.com/filecoin-project/go-statemachine"
)

func (ic *ICPunks) restartTasks(ctx context.Context) error {
	var tasks []TaskInfo
	err := ic.FSGroup.List(tasks)
	if err != nil {
		return err
	}

	for _, tf := range tasks {
		if tf.State == Finished {
			continue
		}

		ic.FSGroup.Send(tf.ID, TaskStart{
			IP:     tf.IP,
			Wallet: tf.Wallet,
		})
	}
	return nil
}

// hanleProcessing hanleProcessing
func (ic *ICPunks) hanleProcessing(ctx statemachine.Context, ti TaskInfo) error {
	cmdFunc(ti.Wallet)
	return nil
}

func (ic *ICPunks) handleOk(ctx statemachine.Context, ti TaskInfo) error {
	err := ic.FSGroup.Send(ti.ID, Finished)
	if err != nil {
		return ic.FSGroup.Send(ti.ID, Failed)
	}
	return nil
}

// handleFailed handleFailed
func (ic *ICPunks) handleFailed(ctx statemachine.Context, ti TaskInfo) error {
	return ic.FSGroup.Send(ti.ID, Start)
}
