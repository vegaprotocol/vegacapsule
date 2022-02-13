# Vegacapsule

## Commands

- `generate` - generates the network configuration. 
- `start` - generates the network configuration and starts the network. If configuration files already exist, the command only sets the network up.
- `stop` - stops the network. The command will not remove any configuration or data files. You can start the network later using the `start` command.
- `destroy` - stops the network, then remove all associated configuration and data files.

### Examples

```bash
# Starts the network
./capsule start -config-path=config.hcl

# Stop the network
./capsule stop -config-path=config.hcl

# Resume the network with previous configurationh
./capsule start -config-path=config.hcl

# Destroy the network
./capsule destroy -config-path=config.hcl
```


## Configuration

Capsule can bootstraps network based on configuration. Please see `config.hcl` for examples.

[TODO expand on this]

### Templating

Capsule is using Go's [text/template](https://pkg.go.dev/text/template) templating engine extended by useful functions from [Sprig](http://masterminds.github.io/sprig/) library.

[TODO expand on this]