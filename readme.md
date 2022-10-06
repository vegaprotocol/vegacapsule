# Vega Capsule
Use Vega Capsule to create an instance of the Vega network on your computer to experiment with using the protocol. 
* Become familiar with Vega, and run commands and API scripts in a controlled environment
* Try out liquidity strategies locally before using the public testnet
* Practice with the market creation process, to make sure proposals will be accepted for a vote
* Simulate network conditions ahead of putting forward a network configuration change proposal
* Simulate market conditions or price scenarios without being concerned about unexpected user behaviour

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
git config --global url."git@github.com:vegaprotocol".insteadOf "https://github.com/vegaprotocol"
cd vegacapsule
```
- Build Capsule from source
```bash
go install
```
- Validate Capsule installation
```bash
vegacapsule --help
```

3. #### Install dependepcies
[Install Vega binaries](install_vega_bins.md). Installs **vega**, **data-node** and **vegawallet** binaries on your machine and.

This step can be skipped if network when network is bootstrapped with --install flag. See 

### Start Capsule Network
1. Start nomad
```bash
vegacapsule nomad
```
**Note**: You may need to set the `GOBIN` environment variable to run it.

2. Bootstrap network

In another Terminal window run bootstrap command to generate and start new network.

#### Bootstrap with preinstalled binaries ####
This step requires preinstalled **vega**, **data-node** and **vegawallet** binaries.
Plese refer to [Install Vega binaries](install_vega_bins.md).

```bash
vegacapsule network bootstrap --config-path=net_confs/config.hcl
```

#### Bootstrap with autoinstall ####
Bootstrap with autoinstall will automatically download required binaries as a first step of the process.
Either **--install** or **--install-release-tag** flags can be used. The former installes latest version and the 
latter installes from given release tag.

```bash
vegacapsule network bootstrap --config-path=net_confs/config.hcl --install
```

```bash
vegacapsule network bootstrap --config-path=net_confs/config.hcl --install-release-tag v0.54.0
```

3. Check Nomad console by opening http://localhost:4646/

## Restoring network from checkpoint
### Bootstrapping a new network

1. First generate the network
```bash
vegacapsule network generate --config-path=net_confs/config.hcl
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

## Logs

Logs from all jobs are stored by default to $CAPSULE_HOM/logs. There is a CLI availible to read them.

To read all logs per single job:
```
vegacapsule logs --job-id $JOBID
```

To follow all logs per job:
```
vegacapsule logs --job-id $JOBID --follow
```

For more information please check
```
vegacapsule logs --help
```

## Commands

You can see all available commands calling the `vegacapsule --help` command.

### Examples

```bash
# Generate the network config files
vegacapsule network generate --home-path=/var/tmp/veganetwork/testnetwork --config-path=net_confs/config.hcl

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


### HCL functions available to use in the configuration files.

You can use HCL functions in the config.hcl. List of available functions is available in the source code of the [config/hcl_eval_context.go](config/hcl_eval_context.go).

Example of usage function in the `config.hcl` is:

```hcl
vega_binary_path = "vega"

network "testnet" {
	ethereum {
    chain_id   = "1440"
    network_id = "1441"
    endpoint   = format("https://ropsten.infura.io/v3/%s", env("INFURA_API_KEY"))
  }

  ...
```

You can find more examples in the [nomad documentation](https://www.nomadproject.io/docs/job-specification/hcl2).



[TODO expand on this]

### Templating

Capsule is using Go's [text/template](https://pkg.go.dev/text/template) templating engine extended by useful functions from [Sprig](http://masterminds.github.io/sprig/) library.

[TODO expand on this]


### Troubleshooting

#### Missing the `GOBIN` environment variable 

```
Error: GOBIN environment variable has not been found - please set install-path flag instead
```

The error may happen during the `vegacapsule nomad` command. To solve it, set the environment variable with the following command: `export GOBIN="$HOME/go/bin"`.