# Vegacapsule

## Quick start

### Pre-start
1. Make sure Docker is running on your machine
```bash
docker version
```
2. If you use a private docker image, run `docker login` to authenticate with private package repository (eg. Github).
```bash
cat PATH_TO_YOUR_TOKEN | docker login https://ghcr.io -u "YOUR_USER_NAME" --password-stdin
```

### Start

1. Clone repository
```bash
git clone git@github.com:vegaprotocol/vegacapsule.git
cd vegacapsule
```
2. Build vega capsule from source
```bash
go install
```
3. Start nomad
```bash
vegacapsule nomad
```

4. In another Terminal window run bootstrap command to generate and start new network
```bash
vegacapsule network bootstrap --config-path=config.hcl
```
5. Check Nomad console by opening http://localhost:4646/

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

You can see all available commands callin the `vegacapsule --help` command.

### Available commands

- `network` - Manages network
- `nodes` - Manages nodes sets
- `nomad` - Starts Nomad instance locally
- `state` - Manages vegacapsule state
- `ethereum` - Interacts with the ethereum network

### Commands to control network

To generate network configuration files, use one of the following commands:

- `generate` - generates the network configuration. Capsule puts network all files in the folder, which you set in the config file as the `output_dir` parameter.
- `bootstrap` - generates the network config files and starts the network in the same command. The `generate` command executes both the `generate` and the `start` internally.

All below commands require generated network configuration. If configuration files are missing, an error is returned.

- `start` - starts the network. 
- `stop` - stops the network. The command will not remove any configuration or data files. You can start the network later using the `start` command.
- `destroy` - stops the network, then removes all associated configuration and data files.

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

### Commands to run nomad

A helper command allows you to download and install (if the correct version `nomad` command is missing on your computer) and set up a nomad if you do not already have one. 

- `nomad` - starts (and installs) simple nomad agent to be run locally in `dev` mode. Instead of this command, you can run nomad manually (`nomad agent -dev -bind=0.0.0.0 -config=client.hcl`)

```bash
vegacapsule nomad
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
