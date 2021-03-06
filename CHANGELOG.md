# Changelog

## Unreleased (v0.2.0)

### 🚨 Breaking changes
- [](https://github.com/vegaprotocol/vegacapsule/issues/xxxx) -

### 🗑️ Deprecation
- [](https://github.com/vegaprotocol/vegacapsule/issues/xxxx) -

### 🛠 Improvements
- [164](https://github.com/vegaprotocol/vegacapsule/issues/164) Update contributor information
- [145](https://github.com/vegaprotocol/vegacapsule/issues/145) Update Nomad version and allow Nomad to be installed to PATH
- [134](https://github.com/vegaprotocol/vegacapsule/issues/134) Add support for Clef and allow templating of some node set config fields
- [139](https://github.com/vegaprotocol/vegacapsule/issues/139) Allow non validator nodes to be iterated during wallet configuration
- [125](https://github.com/vegaprotocol/vegacapsule/issues/125) Update network state when `--update-network` flag is passed to the `template nomad` cmd
- [149](https://github.com/vegaprotocol/vegacapsule/issues/149) Update sentry config to reflect correct architecture
- [191](https://github.com/vegaprotocol/vegacapsule/issues/191) Support built-in Tendermint application and version 0.35.8
- [194](https://github.com/vegaprotocol/vegacapsule/issues/194) Set `skip-timeout-commit` value to true to reduce block times
- [75](https://github.com/vegaprotocol/vegacapsule/issues/75) Add support to import pre-generated keys into vegacapsule network
- [190](https://github.com/vegaprotocol/vegacapsule/issues/190) Support multisig with Clef
- [204](https://github.com/vegaprotocol/vegacapsule/pull/204) Add intro section to readme with more about capsule

### 🐛 Fixes
- [167](https://github.com/vegaprotocol/vegacapsule/issues/167) Fix validators filter in tendermint generator
- [188](https://github.com/vegaprotocol/vegacapsule/issues/188) Support new changes for Ethereum RPC endpoint in Vega configuration
- [209](https://github.com/vegaprotocol/vegacapsule/pull/209) Save tendermint template after merge to given file



## v0.1.0

### 🚨 Breaking changes
- [](https://github.com/vegaprotocol/vegacapsule/issues/xxxx) -

### 🗑️ Deprecation
- [](https://github.com/vegaprotocol/vegacapsule/issues/xxxx) -

### 🛠 Improvements
- [43](https://github.com/vegaprotocol/vegacapsule/issues/39) Add support to download nomad on Apple M1 computers
- [89](https://github.com/vegaprotocol/vegacapsule/issues/89) Add ability to set environment variables for docker jobs
- [88](https://github.com/vegaprotocol/vegacapsule/issues/88) Add ability to map ports for docker jobs
- [60](https://github.com/vegaprotocol/vegacapsule/issues/60) Add support for running null chain network
- [97](https://github.com/vegaprotocol/vegacapsule/issues/97) Add automatic binaries download and improve docs
- [108](https://github.com/vegaprotocol/vegacapsule/issues/108) Add templating commands support
- [114](https://github.com/vegaprotocol/vegacapsule/issues/114) Add support for post_start jobs
- [131](https://github.com/vegaprotocol/vegacapsule/issues/131) Update network binaries: vega&data-node=v0.51.1, vegawallet=v0.15.1
- [120](https://github.com/vegaprotocol/vegacapsule/pull/120) Add support for HCL2 in node-set job template
- [122](https://github.com/vegaprotocol/vegacapsule/issues/122) Add support for sentry nodes and loading node sets templates from files
- [124](https://github.com/vegaprotocol/vegacapsule/issues/124) Allow updating network configurations with templating after network is generated
- [139](https://github.com/vegaprotocol/vegacapsule/issues/139) Allow non validator nodes to be iterated during wallet configuration
- [136](https://github.com/vegaprotocol/vegacapsule/issues/136) New templates for a sentry node with data node setup

### 🐛 Fixes
- [117](https://github.com/vegaprotocol/vegacapsule/pull/117) - fix nil dereference panics in config
- [41](https://github.com/vegaprotocol/vegacapsule/issues/40) - persist the network state after it's generated in bootstrap command
- [86](https://github.com/vegaprotocol/vegacapsule/issues/86) - allow overriding config options that default true with falue
