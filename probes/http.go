package probes

import (
	"context"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var httpClient = http.Client{}

func ProbeHTTP(ctx context.Context, id, url string) (bool, error) {
	log.Printf("Probing HTPP with id %q and url %q", id, url)

	ctx, cancel := context.WithTimeout(ctx, singleProbeTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return false, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		if isConnectionErr(err) {
			return false, newRetryableError(err)
		}

		return false, err
	}

	defer resp.Body.Close()

	return resp.StatusCode > 199 && resp.StatusCode < 300, nil
}
