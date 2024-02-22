# Capsule templating docs

Capsule allows templating for genesis file and [node-set](#nodeconfig) configurations like Vega, Tendermint, and Nomad. This is useful for generating configurations specific to a network, or for using one configuration for all node sets.

Capsule uses Go's [text/template](https://pkg.go.dev/text/template) templating engine, extended by useful functions from the [Sprig](http://masterminds.github.io/sprig/) library.

Every template has a [template context](#template-contexts) - a set of (usually runtime generated) variables passed to the template by Capsule
and then used in the template. These template contexts are documented below.

There are some basic templates provided by Capsule and used by some of the provided configurations in the *net_confs* folder in the Vega Capsule GitHub repo.

## Template tool

Capsule includes a tool to test these templates before they get used in [network config](config.md). Plese check `vegacapsule template --help` for more information.

You can test the *template tool* by using some of the provided default templates after the network has been generated.

For example, run command below and compare the outcome with the [template](net_confs/node_set_templates/default/vega_validators.tmpl).

```bash
vegacapsule template node-sets --type vega --path net_confs/node_set_templates/default/vega_validators.tmpl --nodeset-name testnet-nodeset-validators-0-validator
```

## Template contexts
