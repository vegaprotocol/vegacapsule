package probes

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
)

func newPostgresProbeErr(connStr, query string, err error) error {
	return fmt.Errorf("failed to probe Postgres with connection %q and query %q: %w", connStr, query, err)
}

func ProbePostgres(ctx context.Context, id, connStr, query string) error {
	log.Printf("Probing Postgres with id %q connection %q with query %q", id, connStr, query)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return newPostgresProbeErr(connStr, query, err)
	}

	ctx, cancel := context.WithTimeout(ctx, singleProbeTimeout)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return newPostgresProbeErr(connStr, query, err)
	}

	buff := bytes.NewBuffer([]byte{})
	for rows.Next() {
		b := &sql.RawBytes{}
		if err := rows.Scan(b); err != nil {
			return newPostgresProbeErr(connStr, query, err)
		}

		if _, err := buff.Write(*b); err != nil {
			return newPostgresProbeErr(connStr, query, err)
		}
	}

	log.Printf("Probing Postgres with id %q was successful, result: %s", id, buff.String())

	return nil
}
