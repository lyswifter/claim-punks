package fsm

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/filecoin-project/go-statemachine"
)

func (ic *ICPunks) Tigger(ctx context.Context, wallets []string, ip string) error {
	if len(wallets) == 0 {
		return fmt.Errorf("no wallet provides")
	}

	for i, wlt := range wallets {
		// if tf, err := ic.GetTaskInfo(uint64(i)); err == nil && tf != nil {
		// 	log.Printf("task: %d wlt: %s exist", tf.ID, tf.Wallet)
		// 	continue
		// }

		if i > 0 {
			continue
		}

		log.Printf("start task: %d wlt: %s", i, wlt)

		task := TaskStart{
			ID:     TaskNumber(i),
			IP:     ip,
			Wallet: wlt,
		}

		ic.FSGroup.Send(task.ID, task)

		log.Printf("start task for wallet %s", wlt)
	}
	return nil
}

func (ic *ICPunks) restartTasks(ctx context.Context) error {
	tasks, err := ic.ListTasks()
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
	ret, err := ic.cmdFuncs(ti.Wallet)
	if err != nil {
		return ctx.Send(TaskFailed{})
	}

	log.Printf("hanleProcessing: %s for: %s-%s", ret, ti.IP, ti.Wallet)
	ctx.Send(TaskFinished{TokenID: ret})
	return nil
}

func (ic *ICPunks) handleOk(ctx statemachine.Context, ti TaskInfo) error {

	// go reportStatus("work done finished", statOk)

	return ctx.Send(TaskFinished{})
}

// handleFailed handleFailed
func (ic *ICPunks) handleFailed(ctx statemachine.Context, ti TaskInfo) error {
	return ctx.Send(TaskFailed{
		error: fmt.Errorf("handle task: %d err", ti.ID),
	})
}

func (ic *ICPunks) cmdFuncs(wl string) (string, error) {
	start := time.Now()
	err := os.Chdir(path.Join(ic.repopath, ic.workdir, ic.projectdir))
	if err != nil {
		return "", err
	}

	dfx, err := exec.LookPath("dfx")
	if err != nil {
		return "", err
	}

	cmd := exec.Command(dfx, "--identity", wl, "canister", "--network", "ic", "call", "qcg3w-tyaaa-aaaah-qakea-cai", "name")

	data, err := cmd.Output()
	if err != nil {
		log.Fatalf("cmd: %v err: %s", cmd, err.Error())
		return "", err
	}

	log.Printf("cmd %s finished, took: %s", wl, time.Since(start).String())
	return string(data), nil
}

func (ic *ICPunks) ListTasks() ([]TaskInfo, error) {
	var tasks []TaskInfo
	if err := ic.FSGroup.List(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (ic *ICPunks) GetTaskInfo(tid uint64) (*TaskInfo, error) {
	var out TaskInfo
	err := ic.FSGroup.Get(tid).Get(&out)
	return &out, err
}
