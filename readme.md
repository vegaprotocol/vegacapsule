# Vegacapsule

## Quick start

### Pre-start
1. Make sure you have Go 1.17+ installed locally. [Get Go](https://go.dev/doc/install).
```bash
go version
```

1. Make sure Docker is running on your machine. [Get Docker](https://docs.docker.com/get-docker/).
```bash
docker version
```

2. Install vegacapsule
- Clone Capsule repository
```bash
git clone git@github.com:vegaprotocol/vegacapsule.git
cd vegacapsule
```
- Turn off GONOSUMDB for private vega repositories
```bash
export GONOSUMDB="code.vegaprotocol.io/*"
```
- Build Capsule from source
```bash
go install
```
- Validate Capsule installation
```bash
vegacapsule --help
```

3. Install **vega**, **data-node** and **vegawallet** binaries on your machine and.
[Install Vega binaries](install_vega_bins.md).

### Start Capsule Network
1. Start nomad
```bash
vegacapsule nomad
```
2. In another Terminal window run bootstrap command to generate and start new network
```bash
vegacapsule network bootstrap --config-path=config.hcl
```
3. Check Nomad console by opening http://localhost:4646/

## Restoring network from checkpoint
### Bootstrapping a new network

1. First generate the network
```bash
vegacapsule network generate --config-path=config.hcl
```

2. Run restore command to change networks genesis files
```
vegacapsule nodes restore-checkpoint --checkpoint-file PATH_TO_YOUR_CHECKPOINT_FILE
```

3. Lastly the network can be started. It will load it's state from the checkpoint
```
vegacapsule network start
```

### Restoring on existing network

1. Stop the currently running network first (if the network is running)
```bash
vegacapsule network stop
```

2. Reset current network nodes state
```bash
vegacapsule nodes unsafe-reset-all
```

3. Run restore command to change networks genesis files
```
vegacapsule nodes restore-checkpoint --checkpoint-file PATH_TO_YOUR_CHECKPOINT_FILE
```

4. Lastly the network can be started. It will load it's state from the checkpoint
```
vegacapsule network start
```

## Commands

You can see all available commands calling the `vegacapsule --help` command.

### Examples

```bash
# Generate the network config files
vegacapsule network generate -home-path=/var/tmp/veganetwork/testnetwork --config-path=config.hcl

# Starts the network
vegacapsule network start --home-path=/var/tmp/veganetwork/testnetwork

# Stop the network
vegacapsule network stop --home-path=/var/tmp/veganetwork/testnetwork

# Resume the network with previous configurationh
vegacapsule  network start --home-path=/var/tmp/veganetwork/testnetwork

# Destroy the network
vegacapsule network destroy --home-path=/var/tmp/veganetwork/testnetwork
```

### Commands for ethereum network

You can set up the multisig smart contract with the following command:

```bash
vegacapsule ethereum multisig init
```

Procedure executed by the above command:

1. Set threshold to 1
1. Add validators as signers
1. Remove the contract owner from the signers list
1. Set threshold to 667

## Configuration

Capsule can bootstraps network based on configuration. Please see `config.hcl` for examples.

[TODO expand on this]

### Templating

Capsule is using Go's [text/template](https://pkg.go.dev/text/template) templating engine extended by useful functions from [Sprig](http://masterminds.github.io/sprig/) library.

[TODO expand on this]
