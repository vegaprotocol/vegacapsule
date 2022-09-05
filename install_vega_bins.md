# Vega binaries installation

## Install automatically

There is a feature avalible in Capsule that allows fetching supported binaries automatically.

1. Validate that Capsule is installed
```bash
vegacapsule version
```
### Install binaries to custom path
Install the latest version of the vega binaries to given path:

2. A
```bash
vegacapsule install-bins --install-path YOUR_CUSTOM_PATH
```

Alternatively the binaries can be installed from a specific release tag with install-release-tag flag. To do this follow step 'B'
[Check for releases.](https://github.com/vegaprotocol/vega/releases)

Please note that minimum supported release tag is v0.54.0.

2. B
```bash
vegacapsule install-bins --install-path YOUR_CUSTOM_PATH --install-release-tag SPECIFIC_RELEASE_TAG
```

3. Validate that binaries are accessible from chosen path (YOUR_CUSTOM_PATH) and the versions match the ones from previous cmd output. If not, run step 2 again with `--install-path` flag.
```bash
YOUR_CUSTOM_PATH/vega version
YOUR_CUSTOM_PATH/vegawallet version
YOUR_CUSTOM_PATH/data-node version
```

### Globaly install binaries
1.
```bash
vegacapsule install-bins
```

3. Validate that binaries are accessible trough $PATH. And versions matching the one from previous cmd output. If not, please run step 2 again with --install-path flag.
```bash
vega version
vegawallet version
data-node version
```

## Install manually - build from source (more flexible)

Building from source is a more flexible (recomended for local development) because it gives an option of choosing arbitrary version of the binaries.

**Caveat** - not all binaries versions works with current version of Capsule. For more convenient fast installation please refer to [automatic install](#install-automatically)

Prequsities - this step will require Go 1.18+ installed. [Get Go](https://go.dev/doc/install).
```bash
go version
```

### Vega
All required binaries come from a single git repository. To build them follow the below instructions

1. Clone Vega repository
```bash
git clone git@github.com:vegaprotocol/vega.git
```
2. Enter directory and install from source
```bash
cd vega
go install ./...
```

Alternatively, you can build binaries separately:
```bash
cd vega
go install ./cmd/vega
go install ./cmd/data-node
go install ./cmd/vegawallet
```
3. Verify installation
```bash
vega version
vegawallet version
data-node version
```