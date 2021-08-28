package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

const wallet_count = 20

var gdfx = ""
var home = ""

var dfxjson = "dfx.json"
var workDir = "Hell"
var projectName = "Punks"

var wallets = []string{}
var remainingWallets map[string]struct{}

var PunksInHand = []ResultPunk{}
var ClientIP string = ""

type ResultPunk struct {
	IP        string `json:"ip"`
	Wal       string `json:"wal"`
	TokenID   string `json:"tokenID"`
	WalPubkey string `json:"walPubKey"`
	WalPriKey string `json:"walPriKey"`
}

func init() {
	remainingWallets = make(map[string]struct{})
	for i := 0; i < wallet_count; i++ {
		wallets = append(wallets, fmt.Sprintf("%d", i))
	}
}

func main() {
	var err error

	defer func() {
		if err != nil {
			go func() {
				err := reportStatus(err.Error(), statErr)
				if err != nil {
					return
				}
			}()
			log.Fatalf("process exit: %s", err.Error())
		}
	}()

	ip, erra := GetClientIp()
	if erra != nil {
		return
	}
	err = erra
	ClientIP = ip

	err = prepareEnv()
	if err != nil {
		return
	}

	go func() {
		err := reportStatus("prepareEnv finished", statOk)
		if err != nil {
			return
		}
	}()

	err = prepareProject()
	if err != nil {
		return
	}

	go func() {
		err := reportStatus("prepareProject finished", statOk)
		if err != nil {
			return
		}
	}()

	err = prepareIdentities(wallets)
	if err != nil {
		return
	}

	go func() {
		err := reportStatus("prepareIdentities finished", statOk)
		if err != nil {
			return
		}
	}()

	err = tigger()
	if err != nil {
		return
	}

	go func() {
		err := reportStatus("work done finished", statOk)
		if err != nil {
			return
		}
	}()

}

func prepareEnv() error {
	dfx, err := exec.LookPath("dfx")
	if err != nil {
		return err
	}

	if dfx == "" {
		//need download
		cmd := exec.Command("bash", "sh -ci $(curl -fsSL https://sdk.dfinity.org/install.sh)")
		data, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("failed to call pipeCommands(): %v", err)
		}
		log.Printf("output: %s", data)
	}

	gdfx = dfx

	log.Printf("dfx cmd is ok: %s", dfx)
	return nil
}

func prepareProject() error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	home = dirname

	err = os.Chdir(home)
	if err != nil {
		return err
	}

	_, err = os.Stat(workDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(workDir, 0755)
		if err != nil {
			return err
		}
	}

	err = os.Chdir(workDir)
	if err != nil {
		return err
	}

	if _, err = os.Stat(path.Join(projectName, dfxjson)); os.IsNotExist(err) {
		cmdNew := exec.Command(gdfx, "new", projectName)
		err = cmdNew.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func prepareIdentities(wals []string) error {
	err := os.Chdir(path.Join(home, workDir, projectName))
	if err != nil {
		return err
	}

	cmd := exec.Command(gdfx, "identity", "list")
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	strArr := strings.Split(string(out), "\n")

	for _, wal := range wals {
		var isExist bool = false
		for _, str := range strArr {
			if str == wal {
				log.Printf("identity(%s) already exist", wal)
				isExist = true
				break
			}
		}

		if isExist {
			continue
		}

		cmd := exec.Command(gdfx, "identity", "new", wal)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func tigger() error {
	err := os.Chdir(path.Join(home, workDir, projectName))
	if err != nil {
		return err
	}

	type punks struct {
		Wal  string
		Data string
	}

	type remaining struct {
		Err error
		wal string
	}

	var output chan punks = make(chan punks)
	var remainChan chan remaining = make(chan remaining)

	for _, wal := range wallets {
		go func(wl string) error {
			start := time.Now()
			cmd := exec.Command(gdfx, "--identity", wl, "canister", "--network", "ic", "call", "qcg3w-tyaaa-aaaah-qakea-cai", "name")

			data, err := cmd.Output()
			if err != nil {
				log.Printf("cmd: %v err: %s", cmd, err.Error())
				remainChan <- remaining{
					Err: err,
					wal: wl,
				}
				return nil
			}

			output <- punks{
				Wal:  wl,
				Data: string(data),
			}

			log.Printf("cmd %s finished, took: %s", wl, time.Since(start).String())
			return nil
		}(wal)
	}

	ticker := time.NewTicker(30 * time.Second)

loop:
	for {
		select {
		case ou := <-output:
			ret := ResultPunk{
				IP:      ClientIP,
				Wal:     ou.Wal,
				TokenID: ou.Data,
			}
			PunksInHand = append(PunksInHand, ret)

			go func() {
				err := reportResult(ret)
				if err != nil {
					log.Fatalf("reportResult: %v err: %s", ret, err.Error())
					return
				}
			}()

			delete(remainingWallets, ou.Wal)
			log.Printf("Successfully, wallet: %s, punk: %v", ou.Wal, ou.Data)

			if len(PunksInHand) == len(wallets) {
				break loop
			}
		case re := <-remainChan:
			remainingWallets[re.wal] = struct{}{}
			log.Fatalf("Temp encounter problem: %s err: %s", re.wal, re.Err.Error())

			ret, err := tiggerWith(re.wal)
			if err != nil {
				continue
			}

			output <- punks{
				Wal:  re.wal,
				Data: string(ret),
			}
		case <-ticker.C:
			go func() {
				err := reportStatus("working busy", statOk)
				if err != nil {
					return
				}
			}()
		}
	}

	log.Println("finished")

	return nil
}

func tiggerWith(wallet string) (string, error) {
	start := time.Now()
	cmd := exec.Command(gdfx, "--identity", wallet, "canister", "--network", "ic", "call", "qcg3w-tyaaa-aaaah-qakea-cai", "name")

	data, err := cmd.Output()
	if err != nil {
		log.Fatalf("cmd: %v err: %s", cmd, err.Error())
		return "", nil
	}

	log.Printf("retigger cmd %s finished, took: %s", wallet, time.Since(start).String())
	return string(data), nil
}
