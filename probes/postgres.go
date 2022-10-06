package probes

import (
	"bytes"
	"context"
	"database/sql"
	"log"
)

func ProbePostgres(ctx context.Context, id, connStr, query string) ([]byte, error) {
	log.Printf("Probing Postgres with id %q connection %q with query %q", id, connStr, query)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, singleProbeTimeout)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		if isConnectionErr(err) {
			return nil, newRetryableError(err)
		}

		return nil, err
	}

	buff := bytes.NewBuffer([]byte{})
	for rows.Next() {
		b := &sql.RawBytes{}
		if err := rows.Scan(b); err != nil {
			return nil, err
		}

		if _, err := buff.Write(*b); err != nil {
			return nil, err
		}
	}

	return buff.Bytes(), nil
}
