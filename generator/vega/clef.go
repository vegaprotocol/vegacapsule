package vega

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func waitForClef(url string, payload string, timeout time.Duration) (err error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return err
	}

	req.Header = map[string][]string{
		"Content-Type": {"application/json"},
	}

	httpClient := http.Client{Timeout: time.Second * 5}

	testCall := func() error {
		log.Printf("wating for Clef %q to start", url)

		res, err := httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send request to Clef %q: %w", url, err)
		}

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read Clef %q response: %w", url, err)
		}

		var jsonOut struct {
			Result []string `json:"result"`
		}

		log.Printf("received respose from Clef %q %s", url, b)

		if err := json.Unmarshal(b, &jsonOut); err != nil {
			return fmt.Errorf("failed to unmarshal Clef %q response: %w", url, err)
		}

		return nil
	}

	for tmt := time.After(timeout); ; {
		select {
		case <-tmt:
			return fmt.Errorf("wating for %s has timed out: %w", url, err)
		default:
			err = testCall()
			if err == nil {
				return nil
			}
			time.Sleep(time.Second * 1)
		}
	}
}
