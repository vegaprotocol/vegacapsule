package probes

import (
	"context"
	"fmt"
	"log"
	"net"
)

func ProbeTCP(ctx context.Context, id, address string) (bool, error) {
	log.Printf("Probing TCP with id %q address %q", id, address)

	var d net.Dialer
	ctx, cancel := context.WithTimeout(ctx, singleProbeTimeout)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", address)
	if err != nil {
		if isConnectionErr(err) {
			return false, newRetryableError(err)
		}
		return false, fmt.Errorf("failed to connect to TCP probe address: %w", err)
	}
	conn.Close()

	return true, nil
}
