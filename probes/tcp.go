package probes

import (
	"context"
	"fmt"
	"log"
	"net"
)

func ProbeTCP(ctx context.Context, id, address string) error {
	log.Printf("Probing TCP with id %q address %q", id, address)

	var d net.Dialer
	ctx, cancel := context.WithTimeout(ctx, singleProbeTimeout)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("failed to probe TCP address %q: %w", address, err)
	}
	conn.Close()

	log.Printf("Probing TCP with id %q was successfull", id)

	return nil
}
