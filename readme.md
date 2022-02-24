# Vegacapsule

## Commands

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
./capsule generate -config-path=config.hcl

# Starts the network
./capsule start [-home-path=/var/tmp/veganetwork/testnetwork]

# Stop the network
./capsule stop [-home-path=/var/tmp/veganetwork/testnetwork]

# Resume the network with previous configurationh
./capsule start [-home-path=/var/tmp/veganetwork/testnetwork]

# Destroy the network
./capsule destroy [-home-path=/var/tmp/veganetwork/testnetwork]
```


## Configuration

Capsule can bootstraps network based on configuration. Please see `config.hcl` for examples.

[TODO expand on this]

### Templating

Capsule is using Go's [text/template](https://pkg.go.dev/text/template) templating engine extended by useful functions from [Sprig](http://masterminds.github.io/sprig/) library.

[TODO expand on this]