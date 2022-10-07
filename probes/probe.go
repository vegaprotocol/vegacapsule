package probes

import (
	"context"
	"fmt"
	"log"
	"time"

	"code.vegaprotocol.io/vegacapsule/types"
	"golang.org/x/sync/errgroup"
)

var (
	totalProbeTimeout  = time.Second * 20
	singleProbeTimeout = time.Second * 2
)

func probe(ctx context.Context, id, probeType string, call func() error) error {
	t := time.NewTicker(time.Second * 2)

	var err error
	for {
		select {
		case <-ctx.Done():
			fmt.Println("context is done jare: ", err)
			return fmt.Errorf("%s: %w", ctx.Err(), err)
		case <-t.C:
			err = call()
			if err == nil {
				return nil
			}

			log.Printf("Probe with id %q and type %q has failed %q", id, probeType, err)
		}
	}
}

func Probe(ctx context.Context, id string, probes types.ProbesConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), totalProbeTimeout)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	if probes.HTTP != nil {
		call := func() error {
			_, err := ProbeHTTP(ctx, id, probes.HTTP.URL)
			return err
		}

		eg.Go(func() error {
			return probe(ctx, id, "HTTP", call)
		})
	}
	if probes.TCP != nil {
		call := func() error {
			_, err := ProbeTCP(ctx, id, probes.TCP.Address)
			return err
		}

		eg.Go(func() error {
			return probe(ctx, id, "TCP", call)
		})
	}
	if probes.Postgres != nil {
		call := func() error {
			_, err := ProbePostgres(ctx, id, probes.Postgres.Connection, probes.Postgres.Query)
			return err
		}

		eg.Go(func() error {
			return probe(ctx, id, "Postgres", call)
		})
	}

	if err := eg.Wait(); err != nil {
		log.Printf("Probe id %q has failed %s", id, err)
		return fmt.Errorf("failed probes with id %q: %w", id, err)
	}

	return nil
}
