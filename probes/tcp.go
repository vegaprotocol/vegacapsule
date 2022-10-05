package probes

import (
	"context"
	"fmt"
	"net"
)

func ProbeTCP(ctx context.Context, address string) (bool, error) {
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
