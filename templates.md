


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


## *config.NodeConfigTemplateContext*


### Fields

<dl>
<dt>
	<code>NodeNumber</code>  <strong>int</strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *datanode.ConfigTemplateContext*


### Fields

<dl>
<dt>
	<code>NodeHomeDir</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>NodeNumber</code>  <strong>int</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>NodeSet</code>  <strong><a href="#typesnodeset">types.NodeSet</a></strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *faucet.ConfigTemplateContext*


### Fields

<dl>
<dt>
	<code>HomeDir</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>PublicKey</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *genesis.TemplateContext*

Template context also includes functions:
- `.GetEthContractAddr "contract_name"` - returns contract address based on name.
- `.GetVegaContractID "contract_name"` - returns contract vega ID based on name.



### Fields

<dl>
<dt>
	<code>Addresses</code>  <strong>map[string]<a href="#genesissmartcontract">genesis.SmartContract</a></strong>  - required
</dt>

<dd>

Ethereum smart contracts addresses managed by Vega. These can represent bridges or ERC20 tokens.

</dd>

<dt>
	<code>NetworkID</code>  <strong>string</strong>  - required
</dt>

<dd>

Ethereum network id.

</dd>

<dt>
	<code>ChainID</code>  <strong>string</strong>  - required
</dt>

<dd>

Ethereum chain id.

</dd>



</dl>

---


## *tendermint.ConfigTemplateContext*


### Fields

<dl>
<dt>
	<code>TendermintNodePrefix</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>VegaNodePrefix</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>NodeNumber</code>  <strong>int</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>NodesCount</code>  <strong>int</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>NodeSet</code>  <strong><a href="#typesnodeset">types.NodeSet</a></strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *vega.ConfigTemplateContext*


### Fields

<dl>
<dt>
	<code>TendermintNodePrefix</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>VegaNodePrefix</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>DataNodePrefix</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>ETHEndpoint</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>NodeMode</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>FaucetPublicKey</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>