package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func executeBinary(binaryPath string, args []string, v interface{}) ([]byte, error) {
	command := exec.Command(binaryPath, args...)

	var stdOut, stErr bytes.Buffer
	command.Stdout = &stdOut
	command.Stderr = &stErr

	if err := command.Run(); err != nil {
		return nil, fmt.Errorf("%s: %s", stErr.String(), err.Error())
	}

	if v == nil {
		return stdOut.Bytes(), nil
	}

	if err := json.Unmarshal(stdOut.Bytes(), v); err != nil {
		// TODO Maybe failback to text parsing instead??
		return nil, err
	}

	return nil, nil
}

func createSmartContractsDir(OutputDir string, smartContractsDir string) error {
	scd := path.Join(OutputDir, smartContractsDir)

	if err := os.MkdirAll(scd, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func ganacheCheck() bool {

	for {

		time.Sleep(1 * time.Second)
		postBody, _ := json.Marshal(map[string]string{
			"method": "web3_clientVersion",
		})
		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post("http://127.0.0.1:8545/", "application/json", responseBody)
		if err != nil {
			log.Println("ganache not yet ready", err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			continue
		}

		if strings.Contains(string(body), "EthereumJS") {
			log.Println("ganache is ready")
			return true
		}
		continue
	}

}
