# Vega Capsule

Vega Capsule allows you to create a local instance of the Vega network on your computer. This can be used to experiment with using the protocol, such as:

* Become familiar with Vega, and run commands and API scripts in a controlled environment
* Try out liquidity strategies locally before using the public testnet
* Practice with the market creation process, to make sure proposals will be accepted for a vote
* Simulate network conditions ahead of putting forward a network configuration change proposal
* Simulate market conditions or price scenarios without being concerned about unexpected user behaviour

## Pre-start

In order to use Vega Capsule you will need to install Go and Docker:

1. Install [Go 1.19 or later](https://go.dev/doc/install) locally on your machine. Check you have the correct version installed using the following command:
```bash
go version
```

2. Install [Docker](https://docs.docker.com/get-docker/) locally on your machine. Check you have the correct version installed using the following command:
```bash
docker version
```

## Quick start

1. Install vegacapsule
- Clone the [Vega Capsule](https://github.com/vegaprotocol/vegacapsule) repository
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

2. Start Nomad
- Nomad comes built in with the installed Vega Capsule binaries 
```bash
vegacapsule nomad
```
> ⚠️ Information: 
> You may need to set the `GOBIN` environment variable to start Nomad. If this is required see the <a href="#common-issues">Common Issues</a> information.

3. Bootstrap with auto-installed dependancies
- Bootstrap with auto-install will automatically download the required Vega binaries during the bootstrapping process. 

- Use **--install** to install the latest version of Vega
```bash
vegacapsule network bootstrap --config-path=net_confs/config.hcl --install
```

- Use **--install-release-tag** to install a given version of Vega
```bash
vegacapsule network bootstrap --config-path=net_confs/config.hcl --install-release-tag v0.54.0
```

> ⚠️ Information:
> The example bootstrap commands use the base default network configuration. The [network configurations](./net_confs) directory contains a number of defaults for various use cases for example with front end dApps or nullblockchain defined. Find out more about the <a href="#configuration">Configuration</a> fields.

4. Check Nomad console is running by opening the [Nomad UI](http://localhost:4646/) in a web browser.

## Configuration and templates

Vega Capsule is highly configurable allowing for networks to be created as desired.

### Configuration

Vega Capsule can bootstrap a network based on pre-determined configuration. See the [configuration documentation](./config.md) to find out more about how to configure your network.

The installed software comes with a number of defaults that can be used, these are in the [network configurations](./net_confs) directory.

### Templating

Capsule is using Go's [text/template](https://pkg.go.dev/text/template) templating engine extended by useful functions from [Sprig](http://masterminds.github.io/sprig/) library. See the [template documentation](./templates.md) to find out more about how to configure your network.

The installed software comes with a number of defaults that can be used, these are in the [node set templates](./node_set_templates) directory.

## Using Vega Capsule
In order to start using Vega Capsule to create assets and markets the following actions should be taken.

### Commands for the Ethereum network

In order for the protocol to authorise function execution, such as deposits and withdrawals, the validators need to be set as signers and the thresholds set on the [multisig smart contract](https://github.com/vegaprotocol/MultisigControl#readme). 

1. Set up the multisig smart contract

```bash
vegacapsule ethereum multisig init
```

The command will execute the following procedures:

1. Set the signature threshold to 1
1. Add the validators as signers to the smart contract
1. Remove the contract owner from the signers list
1. Set the signature threshold to 667


### Wallet

TODO: add basic info on creating a wallet

### Depositing/staking & minting of Ethereum assets

Please note following commands are only for Ethereum assets.
All availible assets can be listed on Data Node REST API under `$DATA_NODE_URL/assets`.

#### Examples

Variables used in examples:

`PUB_KEY` - the wallets public key to deposit or stake to.

`AMOUNT` - the amount the be deposited, staked or minted.

`ASSET_SYMBOL` - symbol of the asset to be deposited, staked or minted. It can be found via the Data Node endpoint above.

`ETH_ADDR` - Ethereum adddress for assets to be minted to.

```bash
# Deposit asset to specific Vega key
vegacapsule ethereum asset deposit --amount $AMOUNT --asset-symbol $ASSET_SYMBOL --pub-key $PUB_KEY

# Stake asset to specific Vega key
vegacapsule ethereum asset stake --amount $AMOUNT --asset-symbol $ASSET_SYMBOL --pub-key $PUB_KEY

# Mint asset to specific Ethereum address
vegacapsule ethereum asset mint --amount $AMOUNT --asset-symbol $ASSET_SYMBOL --to-addr $ETH_ADDR
```

Validating that asset has been deposited on Data Node REST API: `$DATA_NODE_URL/parties/$PUB_KEY/accounts`

Validating that asset has been staked on Data Node REST API: `$DATA_NODE_URL/parties/$PUB_KEY/stake`

## Common Commands

You can see all available commands by calling the `vegacapsule --help` command.

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

> ⚠️ Information:
> Capsule preserves some files when starting and stopping the network, for example any pre-generated keys, the genesis file, and any node configurations in the [network configuration file](https://github.com/vegaprotocol/vegacapsule/tree/main/net_confs). In order to start network with new values in these files, use  the `vegacapsule network destroy` command.

## Troubleshooting

### Logs
Vega Capsule captures the logs from the nodes running on the network. These can be used to investigate what is happening on the network. Should an issue be found with Vega Capsule the logs from the time of the incident should be supplied with any issue raised. 

Logs from all jobs are stored by default to `$CAPSULE_HOME/logs`. There is a CLI availible to access and read them.

- To read all logs per single job:

```bash
vegacapsule logs --job-id $JOBID
```

- To follow all logs per job:

```bash
vegacapsule logs --job-id $JOBID --follow
```

- For more information please check

```bash
vegacapsule logs --help
```

### Common issues
This details commonly seen issues users may face when setting up a Vega Capsule network

#### Missing the `GOBIN` environment variable

When trying to start Nomad with the `vegacapsule nomad` command, this error may be presented:

```
Error: GOBIN environment variable has not been found - please set install-path flag instead
```

To solve it, set the environment variable with the following command:

```bash
export GOBIN="$HOME/go/bin"
```

#### Being on a different branch

When trying to bootstrap the network with the `vegacapsule network bootstrap` command, this error may be presented:

```
Error: failed to start network: failed to start network: failed to start vega network: failed to run node sets: failed to wait for node sets: failed to run testnet-nodeset-xxxxxxx job: starting deadline has been exceeded
```

To solve it, ensure that you are on the `main` branch. This can be checked with the command `git branch`.
