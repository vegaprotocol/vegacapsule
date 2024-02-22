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

## *config.NodeConfigTemplateContext*

### Fields

<dl>
<dt>
	<code>NodeNumber</code>  <strong>int</strong>  - required
</dt>

<dd>

Absolute index of the generated node.

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

Path to home directory of the Faucet.

</dd>

<dt>
	<code>PublicKey</code>  <strong>string</strong>  - required
</dt>

<dd>

Public key of the Faucet.

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

Ethereum smart contract addresses created by Vega. These can represent bridges or ERC20 tokens.

</dd>

<dt>
	<code>NetworkID</code>  <strong>string</strong>  - required
</dt>

<dd>

Ethereum network ID.

</dd>

<dt>
	<code>ChainID</code>  <strong>string</strong>  - required
</dt>

<dd>

Ethereum chain ID.

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

## *types.NodeSet*

### Fields

<dl>
<dt>
	<code>GroupName</code>  <strong>string</strong>  - required
</dt>

<dd>

Name that represents a group of the same node sets.

</dd>

<dt>
	<code>Name</code>  <strong>string</strong>  - required
</dt>

<dd>

Name of a specific node set in a node sets group.

</dd>

<dt>
	<code>Mode</code>  <strong>string</strong>  - required
</dt>

<dd>

Mode of the node set. Can be `validator` or `full` (full means a non validating node).

</dd>

<dt>
	<code>Index</code>  <strong>int</strong>  - required
</dt>

<dd>

Index is a position and order in which the node set has been generated respective to all other created node sets.
It goes from 0-N where N is the total number of node sets.


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

GroupIndex is an index of the group that this node set belongs to. Related to GroupName.

</dd>

<dt>
	<code>Vega</code>  <strong><a href="#typesveganode">types.VegaNode</a></strong>  - required
</dt>

<dd>

Information about generated Vega node.

</dd>

<dt>
	<code>Tendermint</code>  <strong><a href="#typestendermintnode">types.TendermintNode</a></strong>  - required
</dt>

<dd>

Information about generated Tendermint node.

</dd>

<dt>
	<code>DataNode</code>  <strong><a href="#typesdatanode">types.DataNode</a></strong>  - optional
</dt>

<dd>

Information about generated Data node.

</dd>

<dt>
	<code>Visor</code>  <strong><a href="#typesvisor">types.Visor</a></strong>  - optional
</dt>

<dd>

Information about generated Visor instance.

</dd>

<dt>
	<code>PreGenerateJobs</code>  <strong>[]<a href="#typesnomadjob">types.NomadJob</a></strong>  - required
</dt>

<dd>

Jobs that have been started before the nodes were generated.

</dd>

<dt>
	<code>PreStartProbe</code>  <strong><a href="#typesprobesconfig">types.ProbesConfig</a></strong>  - optional
</dt>

<dd>

Pre start probes.

</dd>

<dt>
	<code>NomadJobRaw</code>  <strong>string</strong>  - optional
</dt>

<dd>

Stores custom Nomad job definition of this node set.

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

Ethereum address.

</dd>

<dt>
	<code>Vega</code>  <strong>string</strong>  - required
</dt>

<dd>

Vega contract ID.

</dd>



</dl>

---

## *types.VegaNode*

Represents generated Vega node.

### Fields

<dl>
<dt>
	<code>GeneratedService</code>  <strong><a href="#typesgeneratedservice">types.GeneratedService</a></strong>  - required
</dt>

<dd>

General information about the node.

</dd>

<dt>
	<code>Mode</code>  <strong>string</strong>  - required
</dt>

<dd>

Mode of the node - `validator` or `full`.

</dd>

<dt>
	<code>NodeWalletPassFilePath</code>  <strong>string</strong>  - required
</dt>

<dd>

Path to generated node wallet passphrase file.


<blockquote>Only present if `mode = validator`.</blockquote>
</dd>

<dt>
	<code>NodeWalletInfo</code>  <strong><a href="#typesnodewalletinfo">types.NodeWalletInfo</a></strong>  - optional
</dt>

<dd>

Information about generated/imported node wallets.


<blockquote>Only present if `mode = validator`.</blockquote>
</dd>

<dt>
	<code>BinaryPath</code>  <strong>string</strong>  - required
</dt>

<dd>

Path to binary used to generate and run the node.

</dd>



</dl>

---

## *types.TendermintNode*

Represents generated Tendermint node.

### Fields

<dl>
<dt>
	<code>GeneratedService</code>  <strong><a href="#typesgeneratedservice">types.GeneratedService</a></strong>  - required
</dt>

<dd>

General information about the node.

</dd>

<dt>
	<code>NodeID</code>  <strong>string</strong>  - required
</dt>

<dd>

ID of the Tendermint node.

</dd>

<dt>
	<code>GenesisFilePath</code>  <strong>string</strong>  - required
</dt>

<dd>

File path of the genesis file used to bootstrap the network.

</dd>

<dt>
	<code>BinaryPath</code>  <strong>string</strong>  - required
</dt>

<dd>

Path to binary used to generate and run the node.

</dd>

<dt>
	<code>ValidatorPublicKey</code>  <strong>string</strong>  - required
</dt>

<dd>

Generated public key of the Tendermint validator.

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

General information about the node.

</dd>

<dt>
	<code>BinaryPath</code>  <strong>string</strong>  - required
</dt>

<dd>

Path to binary used to generate and run the node.

</dd>

<dt>
	<code>UniqueSwarmKey</code>  <strong>string</strong>  - required
</dt>

<dd>

Unique IPFS swarm key for this network

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

General information about Visor.

</dd>

<dt>
	<code>BinaryPath</code>  <strong>string</strong>  - required
</dt>

<dd>

Path to binary used to generate and run the node.

</dd>



</dl>

---

## *types.NomadJob*

Represents a raw Nomad job.

### Fields

<dl>
<dt>
	<code>ID</code>  <strong>string</strong>  - required
</dt>

<dd>

Custom selected ID - name of the job.

</dd>

<dt>
	<code>NomadJobRaw</code>  <strong>string</strong>  - required
</dt>

<dd>

Nomad job definition.

</dd>



</dl>

---

## *types.ProbesConfig*

Allows the user to define pre start probes on external services.

### Fields

<dl>
<dt>
	<code>HTTP</code>  <strong><a href="#typeshttpprobe">types.HTTPProbe</a></strong>  - optional
</dt>

<dd>

Allows the user to probe HTTP endpoint.

</dd>

<dt>
	<code>TCP</code>  <strong><a href="#typestcpprobe">types.TCPProbe</a></strong>  - optional
</dt>

<dd>

Allows the user to probe TCP socker.

</dd>

<dt>
	<code>Postgres</code>  <strong><a href="#typespostgresprobe">types.PostgresProbe</a></strong>  - optional
</dt>

<dd>

Allows the user to probe Postgres database with a query.

</dd>

### Complete example

```hcl
pre_start_probe {
  ...
}

```

</dl>

---

## *types.GeneratedService*

Represents any generated Capsule service.

### Fields

<dl>
<dt>
	<code>Name</code>  <strong>string</strong>  - required
</dt>

<dd>

Name of the service.

</dd>

<dt>
	<code>HomeDir</code>  <strong>string</strong>  - required
</dt>

<dd>

Path to home directory of the service.

</dd>

<dt>
	<code>ConfigFilePath</code>  <strong>string</strong>  - required
</dt>

<dd>

Path to service configuration.

</dd>



</dl>

---

## *types.NodeWalletInfo*

Information about node wallets.

### Fields

<dl>
<dt>
	<code>EthereumAddress</code>  <strong>string</strong>  - required
</dt>

<dd>

Ethereum account address.


<blockquote>Only available when Key Store wallet is used.</blockquote>
</dd>

<dt>
	<code>EthereumPrivateKey</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>EthereumPassFilePath</code>  <strong>string</strong>  - required
</dt>

<dd>

Path to file where Ethereum wallet key is stored.

</dd>

<dt>
	<code>EthereumClefRPCAddress</code>  <strong>string</strong>  - required
</dt>

<dd>

Address of Clef wallet.


<blockquote>Only available when Clef wallet is used.</blockquote>
</dd>

<dt>
	<code>VegaWalletName</code>  <strong>string</strong>  - required
</dt>

<dd>

Name of the Vega wallet.

</dd>

<dt>
	<code>VegaWalletID</code>  <strong>string</strong>  - required
</dt>

<dd>

ID of Vega wallet.

</dd>

<dt>
	<code>VegaWalletPublicKey</code>  <strong>string</strong>  - required
</dt>

<dd>

Public key used from the Vega wallet.

</dd>

<dt>
	<code>VegaWalletRecoveryPhrase</code>  <strong>string</strong>  - required
</dt>

<dd>

Recovery phrase from the Vega wallet.

</dd>

<dt>
	<code>VegaWalletPassFilePath</code>  <strong>string</strong>  - required
</dt>

<dd>

File path of the Vega wallet passphrase.

</dd>



</dl>

---

## *types.HTTPProbe*

Allows the user to probe HTTP endpoint.

### Fields

<dl>
<dt>
	<code>URL</code>  <strong>string</strong>  - required
</dt>

<dd>

URL of the HTTP endpoint.

</dd>

### Complete example

```hcl
http {
  url = "http://localhost:8002"
}

```

</dl>

---

## *types.TCPProbe*

Allows the user to probe TCP socket.

### Fields

<dl>
<dt>
	<code>Address</code>  <strong>string</strong>  - required
</dt>

<dd>

Address of the TCP socket.

</dd>

### Complete example

```hcl
tcp {
  address = "localhost:9009"
}

```

</dl>

---

## *types.PostgresProbe*

Allows the user to probe Postgres database.

### Fields

<dl>
<dt>
	<code>Connection</code>  <strong>string</strong>  - required
</dt>

<dd>

Postgres connection string.

</dd>

<dt>
	<code>Query</code>  <strong>string</strong>  - required
</dt>

<dd>

Test query.

</dd>

### Complete example

```hcl
postgres {
  connection = "user=vega dbname=vega password=vega port=5232 sslmode=disable"
  query      = "select 10 + 10"
}

```

</dl>

---


