package types

/*
description: Allows define pre start probes on external services.
example:

	type: hcl
	value: |
			pre_start_probe {
				...
			}
*/
type ProbesConfig struct {
	/*
		description: Allows to probe HTTP endpoint.
		example:

			type: hcl
			value: |
					http {
						...
					}
	*/
	HTTP *HTTPProbe `hcl:"http,block" template:""`

	/*
		description: Allows to probe TCP socker.
		example:

			type: hcl
			value: |
					tcp {
						...
					}
	*/
	TCP *TCPProbe `hcl:"tcp,block" template:""`

	/*
		description: Allows to probe Postgres database with a query.
		example:

			type: hcl
			value: |
					postgres {
						...
					}
	*/
	Postgres *PostgresProbe `hcl:"postgres,block" template:""`
}

/*
description: Allows the user to probe HTTP endpoint.
example:

	type: hcl
	value: |
			http {
				url = "http://localhost:8002"
			}
*/
type HTTPProbe struct {
	// description: URL of the HTTP endpoint.
	URL string `hcl:"url" template:""`
}

/*
description: Allows to probe TCP socket.
example:

	type: hcl
	value: |
			tcp {
				address = "localhost:9009"
			}
*/
type TCPProbe struct {
	// description: Address of the TCP socket.
	Address string `hcl:"address" template:""`
}

/*
description: Allows to probe Postgres database.
example:

	type: hcl
	value: |
			postgres {
				connection = "user=vega dbname=vega password=vega port=5232 sslmode=disable"
				query = "select 10 + 10"
			}
*/
type PostgresProbe struct {
	// description: Postgres connection string.
	Connection string `hcl:"connection" template:""`
	// description: Test query.
	Query string `hcl:"query" template:""`
}
