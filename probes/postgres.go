package probes

import (
	"bytes"
	"context"
	"database/sql"
)

func ProbePostgres(ctx context.Context, connStr, query string) ([]byte, error) {
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
