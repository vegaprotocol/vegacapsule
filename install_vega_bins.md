# Vega binaries installation

## Install automatically

There is a feature avalible in Capsule that allows fetching supported binaries automatically. Your personal Gihub token is required for this step. [Get Github Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)

1. Validate that Capsule is installed
```bash
vegacapsule --help
```
2. Run the install command
```bash
vegacapsule install-deps
```

## Install manually - build from source (more flexible)
Building from source is a more flexible (recomended for local development) because it gives an option of choosing arbitrary version of the binaries.

**Caveat** - not all binaries versions works with current version of Capsule. For more convenient fast installation please refer to [automatic install](#install-automatically)

Prequsities - this step will require Go 1.17+ installed. [Get Go](https://go.dev/doc/install).
```bash
go version
```

### Vega
1. Clone Vega repository
```bash
git clone git@github.com:vegaprotocol/vega.git
```
2. Turn off GONOSUMDB for private vega repositories
```bash
export GONOSUMDB="code.vegaprotocol.io/*"
```
3. Enter directory and install from source
```bash
cd vega
go install ./cmd/vega
```
4. Validate installation
```bash
vega version
```
### Data node
1. Clone Data Node repository
```bash
git clone git@github.com:vegaprotocol/data-node.git
```
2. Turn off GONOSUMDB for private vega repositories
```bash
export GONOSUMDB="code.vegaprotocol.io/*"
```
3. Enter the directory and install from source
```bash
cd data-node
go install ./cmd/vega
```
4. Validate installation
```bash
data-node version
```
### Vegawallet
1. Clone Vega Wallet repository
```bash
git clone git@github.com:vegaprotocol/vegawallet.git
```
2. Turn off GONOSUMDB for private vega repositories
```bash
export GONOSUMDB="code.vegaprotocol.io/*"
```
3. Enter the directory and install from source
```bash
cd vegawallet
go install .
```
4. Validate installation
```bash
vegawallet version
```