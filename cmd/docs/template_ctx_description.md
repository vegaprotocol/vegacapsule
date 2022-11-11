# Capsulte templating docs

Capsule allows templating for genesis file and [node-sets](#nodeconfig) configurations
like Vega, Tendermint, and Nomad. This is useful for generating configurations specific to a network
or using one configuration for all node set.

Capsule is using Go's [text/template](https://pkg.go.dev/text/template) templating engine extended by useful functions from [Sprig](http://masterminds.github.io/sprig/) library.

Every single template gets it's [template context](#template-contexts) - a set of (usually runtime generated) variables pass to the template by Capsule
that can be use in the template. These template contexts are documented below.

There are some basic templates provided by Capsule and use by some provided configurations in *net_confs* folder.

## Template tool
There is a useful tool as par of Capsule to test these templates before they get used in [network config](config.md).
Plese check `vegacapsule template --help`.

You can test the *template tool* by using some of the provided default templates after the network has been generated.

For example try to run command below and compare the outcome with [the template](net_confs/node_set_templates/default/vega_validators.tmpl).
```bash
vegacapsule template node-sets --type vega --path net_confs/node_set_templates/default/vega_validators.tmpl --nodeset-name testnet-nodeset-validators-0-validator
```

## Template contexts