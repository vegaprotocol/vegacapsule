


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


### Fields

<dl>
<dt>
	<code>Addresses</code>  <strong>map[string]<a href="#genesissmartcontract">genesis.SmartContract</a></strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>ChainID</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>NetworkID</code>  <strong>string</strong>  - required
</dt>

<dd>



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
	<code>NodeNumber</code>  <strong>int</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>NodeSet</code>  <strong><a href="#typesnodeset">types.NodeSet</a></strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>NodeHomeDir</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *visor.ConfigTemplateContext*


### Fields

<dl>
<dt>
	<code>NodeSet</code>  <strong><a href="#typesnodeset">types.NodeSet</a></strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *wallet.ConfigTemplateContext*


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
	<code>WalletPrefix</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>Validators</code>  <strong>[]<a href="#typesnodeset">types.NodeSet</a></strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>NonValidators</code>  <strong>[]<a href="#typesnodeset">types.NodeSet</a></strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *nomad.PreGenerateTemplateCtx*


### Fields

<dl>
<dt>
	<code>Name</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>Index</code>  <strong>int</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>LogsDir</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>CapsuleBinary</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *types.NodeSet*


### Fields

<dl>
<dt>
	<code>GroupName</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>Name</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>Mode</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>Index</code>  <strong>int</strong>  - required
</dt>

<dd>

Index is a node set counter over all created node sets.

</dd>

<dt>
	<code>RelativeIndex</code>  <strong>int</strong>  - required
</dt>

<dd>

RelativeIndex is a counter relative to current node set group. Related to GroupName.

</dd>

<dt>
	<code>GroupIndex</code>  <strong>int</strong>  - required
</dt>

<dd>

GroupIndex is a index of the group where this node set belongs to. Related to GroupName.

</dd>

<dt>
	<code>Vega</code>  <strong><a href="#typesveganode">types.VegaNode</a></strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>Tendermint</code>  <strong><a href="#typestendermintnode">types.TendermintNode</a></strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>DataNode</code>  <strong><a href="#typesdatanode">types.DataNode</a></strong>  - optional
</dt>

<dd>



</dd>

<dt>
	<code>Visor</code>  <strong><a href="#typesvisor">types.Visor</a></strong>  - optional
</dt>

<dd>



</dd>

<dt>
	<code>NomadJobRaw</code>  <strong>string</strong>  - optional
</dt>

<dd>



</dd>

<dt>
	<code>PreGenerateJobs</code>  <strong>[]<a href="#typesnomadjob">types.NomadJob</a></strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>PreStartProbe</code>  <strong><a href="#typesprobesconfig">types.ProbesConfig</a></strong>  - optional
</dt>

<dd>



</dd>



</dl>

---


## *genesis.SmartContract*


### Fields

<dl>
<dt>
	<code>Ethereum</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>Vega</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *types.VegaNode*


### Fields

<dl>
<dt>
	<code>GeneratedService</code>  <strong><a href="#typesgeneratedservice">types.GeneratedService</a></strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>Mode</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>NodeWalletPassFilePath</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>NodeWalletInfo</code>  <strong><a href="#typesnodewalletinfo">types.NodeWalletInfo</a></strong>  - optional
</dt>

<dd>



</dd>

<dt>
	<code>BinaryPath</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *types.TendermintNode*


### Fields

<dl>
<dt>
	<code>GeneratedService</code>  <strong><a href="#typesgeneratedservice">types.GeneratedService</a></strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>NodeID</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>GenesisFilePath</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>BinaryPath</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>ValidatorPublicKey</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *types.DataNode*


### Fields

<dl>
<dt>
	<code>GeneratedService</code>  <strong><a href="#typesgeneratedservice">types.GeneratedService</a></strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>BinaryPath</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *types.Visor*


### Fields

<dl>
<dt>
	<code>GeneratedService</code>  <strong><a href="#typesgeneratedservice">types.GeneratedService</a></strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>BinaryPath</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *types.NomadJob*


### Fields

<dl>
<dt>
	<code>ID</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>NomadJobRaw</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *types.ProbesConfig*


### Fields

<dl>
<dt>
	<code>HTTP</code>  <strong><a href="#typeshttpprobe">types.HTTPProbe</a></strong>  - optional
</dt>

<dd>



</dd>

<dt>
	<code>TCP</code>  <strong><a href="#typestcpprobe">types.TCPProbe</a></strong>  - optional
</dt>

<dd>



</dd>

<dt>
	<code>Postgres</code>  <strong><a href="#typespostgresprobe">types.PostgresProbe</a></strong>  - optional
</dt>

<dd>



</dd>



</dl>

---


## *types.GeneratedService*


### Fields

<dl>
<dt>
	<code>Name</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>HomeDir</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>ConfigFilePath</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *types.NodeWalletInfo*


### Fields

<dl>
<dt>
	<code>EthereumPassFilePath</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>EthereumAddress</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>EthereumPrivateKey</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>EthereumClefRPCAddress</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>VegaWalletID</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>VegaWalletPublicKey</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>VegaWalletRecoveryPhrase</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>VegaWalletName</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>VegaWalletPassFilePath</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *types.HTTPProbe*


### Fields

<dl>
<dt>
	<code>URL</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *types.TCPProbe*


### Fields

<dl>
<dt>
	<code>Address</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *types.PostgresProbe*


### Fields

<dl>
<dt>
	<code>Connection</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>Query</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>



</dl>

---


