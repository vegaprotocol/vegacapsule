# Vegacapsule

## Quick start

### Pre-start
1. Make sure Docker is running on your machine
```bash
docker version
```
2. Run `docker login` to authenticate with private package repository on Github. Github token should be used.
```bash
cat PATH_TO_YOUR_TOKEN | docker login https://ghcr.io -u "YOUR_USER_NAME" --password-stdin
```
3. Test that docker login worked by running
```bash
docker pull ghcr.io/vegaprotocol/devops-infra/ganache:latest
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

4. A) **For M1 users only**!
> Please locate `docker_service "ganache-1"` block in `config.hcl` file and replace image parameter from `ghcr.io/vegaprotocol/devops-infra/ganache:latest` to `ghcr.io/vegaprotocol/devops-infra/ganache:arm64-latest`


4. B) In another Terminal window run bootstrap command to generate and start new network
```bash
vegacapsule bootstrap --config-path=config.hcl
```
5. Check Nomad console by opening http://localhost:4646/

## Commands

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
vegacapsule generate --config-path=config.hcl

# Starts the network
vegacapsule start --home-path=/var/tmp/veganetwork/testnetwork

# Stop the network
vegacapsule stop --home-path=/var/tmp/veganetwork/testnetwork

# Resume the network with previous configurationh
vegacapsule start --home-path=/var/tmp/veganetwork/testnetwork

# Destroy the network
vegacapsule destroy --home-path=/var/tmp/veganetwork/testnetwork
```

### Commands to run nomad

A helper command allows you to download and install (if the correct version `nomad` command is missing on your computer) and set up a nomad if you do not already have one. 

- `nomad` - starts (and installs) simple nomad agent to be run locally in `dev` mode. Instead of this command, you can run nomad manually (`nomad agent -dev -bind=0.0.0.0 -config=client.hcl`)

```bash
vegacapsule nomad
```


## Configuration

Capsule can bootstraps network based on configuration. Please see `config.hcl` for examples.

[TODO expand on this]

### Templating

Capsule is using Go's [text/template](https://pkg.go.dev/text/template) templating engine extended by useful functions from [Sprig](http://masterminds.github.io/sprig/) library.

[TODO expand on this]