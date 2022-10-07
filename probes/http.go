package probes

import (
	"context"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var httpClient = http.Client{}

func newHTTPProbeErr(url string, err error) error {
	return fmt.Errorf("failed to probe HTPP url %q: %w", url, err)
}

func ProbeHTTP(ctx context.Context, id, url string) error {
	log.Printf("Probing HTPP with id %q and url %q", id, url)

	ctx, cancel := context.WithTimeout(ctx, singleProbeTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return newHTTPProbeErr(url, err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return newHTTPProbeErr(url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode > 299 {
		return newHTTPProbeErr(url, fmt.Errorf("status code is not in a range of 200-299"))
	}

	log.Printf("Probing HTTP with id %q was successfull", id)

	return nil
}
