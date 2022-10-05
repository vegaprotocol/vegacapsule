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

func probe(ctx context.Context, call func() error) error {
	t := time.NewTicker(time.Second * 2)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			fmt.Println("probing.........")
			if err := call(); err != nil {
				return err
			}
		}
	}
}

func Probe(ctx context.Context, probes types.ProbesConfig) error {
	log.Printf("running probes for: %+v", probes)
	ctx, cancel := context.WithTimeout(context.Background(), totalProbeTimeout)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	if probes.HTTP != nil {
		fmt.Println("--------- running http probe")
		call := func() error {
			_, err := ProbeHTTP(ctx, probes.HTTP.URL)
			return err
		}

		eg.Go(func() error {
			return probe(ctx, call)
		})
	}
	if probes.TCP != nil {
		fmt.Println("--------- running tcp probe")
		call := func() error {
			_, err := ProbeTCP(ctx, probes.TCP.Address)
			return err
		}

		eg.Go(func() error {
			return probe(ctx, call)
		})
	}
	if probes.Postgres != nil {
		fmt.Println("--------- running postgres probe")
		call := func() error {
			_, err := ProbePostgres(ctx, probes.Postgres.Connection, probes.Postgres.Query)
			return err
		}

		eg.Go(func() error {
			return probe(ctx, call)
		})
	}

	return eg.Wait()
}
