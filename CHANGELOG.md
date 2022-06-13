# Changelog

## Unreleased (0.1.0)

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

### 🐛 Fixes
- [117](https://github.com/vegaprotocol/vegacapsule/pull/117) - fix nil dereference panics in config
- [41](https://github.com/vegaprotocol/vegacapsule/issues/40) - persist the network state after it's generated in bootstrap command
- [86](https://github.com/vegaprotocol/vegacapsule/issues/86) - allow overriding config options that default true with falue
