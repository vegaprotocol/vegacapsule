


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


## Root - *Config*

All parameters from this types are used directly in the config file.
Most of the parameters here are optional and can be left alone.
Please see the example below.



### Fields

<dl>
<dt>
	<code>network</code>  <strong><a href="#networkconfig">NetworkConfig</a></strong>  - required, block 
</dt>

<dd>

Configuration of Vega network and its dependencies.

</dd>

<dt>
	<code>output_dir</code>  <strong>string</strong>  - optional
</dt>

<dd>

Directory path (relative or absolute) where Capsule stores generated folders, files, logs and configurations for network.



Default value: <code>~/.vegacapsule/testnet</code>
</dd>

<dt>
	<code>vega_binary_path</code>  <strong>string</strong>  - optional
</dt>

<dd>

Path (relative or absolute) to vega binary that will be used to generate and run the network.


Default value: <code>vega</code>
</dd>

<dt>
	<code>vega_capsule_binary_path</code>  <strong>string</strong>  - optional
</dt>

<dd>

Path (relative or absolute) of a Capsule binary. The Capsule binary is used to aggregate logs from running jobs
and save them to local disk in Capsule home directory.
See `vegacapsule nomad logscollector` for more info.



Default value: <code>Currently running Capsule instance binary</code>

<blockquote>This optional parameter is used internally. There should never be any need to set it to anything other than default.</blockquote>
</dd>



### Complete example



```hcl
vega_binary_path = "/path/to/vega"

network "your_network_name" {
  ...
}

```


</dl>

---


## *NetworkConfig*

Network configuration allows a user to customise the Capsule Vega network into different shapes based on personal needs.
It also allows the configuration and deployment of different Vega nodes' setups (validator, full) and their dependencies (like Ethereum or Postgres).
It can run custom Docker images before and after the network nodes have started and much more.



### Fields

<dl>
<dt>
	<code>name</code>  <strong>string</strong>  - required, label 
</dt>

<dd>

Name of the network.
All folders generated are placed in the folder with this name.
All Nomad jobs are prefix with the name.


</dd>

<dt>
	<code>genesis_template</code>  <strong>string</strong>  - required | optional if <code>genesis_template_file</code> defined
</dt>

<dd>

[Go template](templates.md) of genesis file that will be used to bootrap the Vega network.
[Example of templated mainnet genesis file](https://github.com/vegaprotocol/networks/blob/master/mainnet1/genesis.json).

The [GenesisTemplateContext](templates.md#genesistemplatecontext) can be used in the template. Example [example](net_confs/genesis.tmpl).



<blockquote>It is recommended that you use `genesis_template_file` param instead.
If both `genesis_template` and `genesis_template_file` are defined, then `genesis_template`
overrides `genesis_template_file`.
</blockquote>

<br />

#### <code>genesis_template</code> example







```hcl
genesis_template = <<EOH
 {
  "app_state": {
   ...
  }
  ..
 }
EOH

```





</dd>

<dt>
	<code>genesis_template_file</code>  <strong>string</strong>  - optional
</dt>

<dd>

Same as `genesis_template` but it allows the user to link the genesis file template as an external file.



<br />

#### <code>genesis_template_file</code> example







```hcl
genesis_template_file = "/your_path/genesis.tmpl"

```





</dd>

<dt>
	<code>ethereum</code>  <strong><a href="#ethereumconfig">EthereumConfig</a></strong>  - required, block 
</dt>

<dd>

Allows the user to define Ethereum network configuration.
This is necessary because Vega needs to be connected to [Ethereum bridges](https://docs.vega.xyz/mainnet/api/bridge)
or it cannot function.



<br />

#### <code>ethereum</code> example







```hcl
ethereum {
  ...
}

```





</dd>

<dt>
	<code>smart_contracts_addresses</code>  <strong>string</strong>  - required | optional if <code>smart_contracts_addresses_file</code> defined, optional 
</dt>

<dd>

Smart contract addresses are addresses of [Ethereum bridges](https://docs.vega.xyz/mainnet/api/bridge) contracts in JSON format.

These addresses should correspond to the chosen network in [Ethereum network](#EthereumConfig) and
can be used in various types of templates in Capsule.
[Example of smart contract address from mainnet](https://github.com/vegaprotocol/networks/blob/master/mainnet1/smart-contracts.json).



<blockquote>It is recommended that you use the `smart_contracts_addresses_file` param instead.
If both `smart_contracts_addresses` and `smart_contracts_addresses_file` are defined, then `genesis_template`
overrides `smart_contracts_addresses_file`.
</blockquote>

<br />

#### <code>smart_contracts_addresses</code> example







```hcl
smart_contracts_addresses = <<EOH
 {
  "erc20_bridge": "...",
  "staking_bridge": "...",
  ...
 }
EOH

```





</dd>

<dt>
	<code>smart_contracts_addresses_file</code>  <strong>string</strong>  - optional
</dt>

<dd>

Same as `smart_contracts_addresses` but it allows you to link the smart contracts as an external file.



<br />

#### <code>smart_contracts_addresses_file</code> example







```hcl
smart_contracts_addresses_file = "/your_path/smart-contratcs.json"

```





</dd>

<dt>
	<code>node_set</code>  <strong>[]<a href="#nodeconfig">NodeConfig</a></strong>  - required, block 
</dt>

<dd>

Allows a user to define multiple nodes sets and their specific configurations.
A node set is a representation of Vega and Data Node nodes.
This is the essential building block of the Vega network.



<br />

#### <code>node_set</code> example



**Validators node set**



```hcl
node_set "validator-nodes" {
  ...
}

```



**Full nodes node set**



```hcl
node_set "full-nodes" {
  ...
}

```





</dd>

<dt>
	<code>wallet</code>  <strong><a href="#walletconfig">WalletConfig</a></strong>  - optional, block 
</dt>

<dd>

Allows for deploying and configuring the [Vega Wallet](https://docs.vega.xyz/mainnet/tools/vega-wallet) instance.
Wallet will not be deployed if this block is not defined.



<br />

#### <code>wallet</code> example







```hcl
wallet "wallet-name" {
  ...
}

```





</dd>

<dt>
	<code>faucet</code>  <strong><a href="#faucetconfig">FaucetConfig</a></strong>  - optional, block 
</dt>

<dd>

Allows for deploying and configuring the [Vega Core Faucet](https://github.com/vegaprotocol/vega/tree/develop/core/faucet#faucet) instance, for supplying builtin assets.
Faucet will not be deployed if this block is not defined.



<br />

#### <code>faucet</code> example







```hcl
faucet "faucet-name" {
  ...
}

```





</dd>

<dt>
	<code>pre_start</code>  <strong><a href="#pstartconfig">PStartConfig</a></strong>  - optional, block 
</dt>

<dd>

Allows the user to define jobs that should run before the node sets start.
It can be used for node sets' dependencies, like databases, mock Ethereum chain, etc..



<br />

#### <code>pre_start</code> example







```hcl
pre_start {
  docker_service "ganache-1" {
    ...
  }
  docker_service "postgres-1" {
    ...
  }
}

```





</dd>

<dt>
	<code>post_start</code>  <strong><a href="#pstartconfig">PStartConfig</a></strong>  - optional, block 
</dt>

<dd>

Allows the user to define jobs that should run after the node sets start.
It can be used for services that depend on a network that is already running, like block explorer or Console.



<br />

#### <code>post_start</code> example







```hcl
post_start {
  docker_service "bloc-explorer-1" {
    ...
  }
  docker_service "vega-console-1" {
    ...
  }
}

```





</dd>



### Complete example



```hcl
network "testnet" {
  ethereum {
    ...
  }

  pre_start {
    ...
  }

  genesis_template_file          = "..."
  smart_contracts_addresses_file = "..."

  node_set "validator-nodes" {
    ...
  }

  node_set "full-nodes" {
    ...
  }
}

```


</dl>

---


## *EthereumConfig*

Allows the user to define the specific Ethereum network to be used.
It can either be one of the [public networks](https://ethereum.org/en/developers/docs/networks/#public-networks) or
a local instance of Ganache.



### Fields

<dl>
<dt>
	<code>chain_id</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>network_id</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>endpoint</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>



### Complete example



```hcl
ethereum {
  chain_id   = "1440"
  network_id = "1441"
  endpoint   = "http://127.0.0.1:8545/"
}

```


</dl>

---


## *NodeConfig*

Represents, and allows the user to configure, a set of Vega (with Tendermint) and Data Node nodes.
One node set definition can be used by applied to multiple node sets (see `count` field) and it uses
templating to distinguish between different nodes and names/ports and other collisions.



### Fields

<dl>
<dt>
	<code>name</code>  <strong>string</strong>  - required, label 
</dt>

<dd>

Name of the node set.
Nomad instances that are part of these nodes are prefixed with this name.


</dd>

<dt>
	<code>mode</code>  <strong>string</strong>  - required
</dt>

<dd>

Determines what mode the node set should run in.



Valid values:

<ul>

<li><code>validator</code></li>

<li><code>full</code></li>
</ul>
</dd>

<dt>
	<code>count</code>  <strong>int</strong>  - required
</dt>

<dd>

Defines how many node sets with this exact configuration should be created.


</dd>

<dt>
	<code>node_wallet_pass</code>  <strong>string</strong>  - optional | required if <code>mode=validator</code> defined
</dt>

<dd>

Defines the password for the automatically generated node wallet associated with the created node.

</dd>

<dt>
	<code>ethereum_wallet_pass</code>  <strong>string</strong>  - optional | required if <code>mode=validator</code> defined
</dt>

<dd>

Defines password for automatically generated Ethereum wallet in node wallet.

</dd>

<dt>
	<code>vega_wallet_pass</code>  <strong>string</strong>  - optional | required if <code>mode=validator</code> defined
</dt>

<dd>

Defines password for automatically generated Vega wallet in node wallet.

</dd>

<dt>
	<code>use_data_node</code>  <strong>bool</strong>  - optional
</dt>

<dd>

Whether or not Data Node should be deployed on node set.

</dd>

<dt>
	<code>visor_binary</code>  <strong>string</strong>  - optional
</dt>

<dd>

Path to [Visor](https://github.com/vegaprotocol/vega/tree/develop/visor) binary.
If defined, Visor is automatically used to deploy Vega and Data nodes.
The relative or absolute path can be used, if only the binary name is defined it automatically looks for it in $PATH.


</dd>

<dt>
	<code>config_templates</code>  <strong><a href="#configtemplates">ConfigTemplates</a></strong>  - required, block 
</dt>

<dd>

Templates that can be used for configurations of Vega and Data nodes, Tendermint and other services.

</dd>

<dt>
	<code>vega_binary_path</code>  <strong>string</strong>  - optional
</dt>

<dd>

Allows user to define a Vega binary to be used in specific node set only.
A relative or absolute path can be used. If only the binary name is defined, it automatically looks for it in $PATH.
This can help with testing different version compatibilities or a protocol upgrade.



<blockquote>Using versions that are not compatible could break the network - therefore this should be used in advanced cases only.</blockquote>
</dd>

<dt>
	<code>pre_generate</code>  <strong><a href="#pregenerate">PreGenerate</a></strong>  - optional, block 
</dt>

<dd>

Allows to run a custom services before the node set is generated.
This can be very useful when generation of the node set might have some extenal depenency for example
a [Clef wallet](https://geth.ethereum.org/docs/clef/introduction).



<blockquote>Clef wallet is a good example - since generating a validator node set requires Ethereum key
to be generated, Clef can be started before the generation starts so that Capsule can generate
the Ethereum key inside of during the generation process.
</blockquote>
</dd>

<dt>
	<code>pre_start_probe</code>  <strong>types.ProbesConfig</strong>  - optional, block 
</dt>

<dd>

Allows to run checks that has to be fulfilled before the node starts.


<blockquote>This can be useful for checking wheter some depedent services has already started or not.
For example databases, mocked services etc.
</blockquote>
</dd>

<dt>
	<code>clef_wallet</code>  <strong><a href="#clefconfig">ClefConfig</a></strong>  - optional, block 
</dt>

<dd>

[Clef](https://geth.ethereum.org/docs/clef/introduction) is one of the
[supported Ethereum wallets](https://docs.vega.xyz/mainnet/node-operators/setup-validator#using-clef) for Vega node.
Capsule supports using Clef and can import pre-generated Ethereum keys from Clef during node set
generation process automatically.

By configuring this paramater Capsule will automatically generate Ethereum keys in Clef and tells Vega to use them.
An example Capsule config setup with Clef can be seen [here](net_confs/config_clef.hcl).


</dd>

<dt>
	<code>nomad_job_template</code>  <strong>string</strong>  - optional
</dt>

<dd>

[Go template](templates.md) of custom Nomad job for node set.

By default Capsule uses predefined Nomad jobs to run the node set on Nomad.
This parameter allows to provide custom Nomad job to represent the generated node set.

The [types.NodeSet](templates.md#types.nodeset) can be used in the template.

Using custom Nomad jobs for node set can break Capsule function properly,
very detailed knowledge is required - therefore it is recommened to leave this parameter
that should be used in advanced cases only.



<blockquote>It is recommended that you use `nomad_job_template_file` param instead.
If both `nomad_job_template` and `nomad_job_template_file` are defined, then `vega`
overrides `nomad_job_template_file`.
</blockquote>
</dd>

<dt>
	<code>nomad_job_template_file</code>  <strong>string</strong>  - optional
</dt>

<dd>

Same as `nomad_job_template` but it allows the user to link the Nomad job template as an external file.



<br />

#### <code>nomad_job_template_file</code> example







```hcl
nomad_job_template_file = "/your_path/vega_config.tmpl"

```





</dd>



### Complete example



```hcl
node_set "validators" {
  count = 2
  mode  = "validator"

  node_wallet_pass     = "n0d3w4ll3t-p4ssphr4e3"
  vega_wallet_pass     = "w4ll3t-p4ssphr4e3"
  ethereum_wallet_pass = "ch41nw4ll3t-3th3r3um-p4ssphr4e3"

  config_templates {
    vega_file       = "./path/vega_validator.tmpl"
    tendermint_file = "./path/tendermint_validator.tmpl"
  }
}

```


</dl>

---


## *WalletConfig*


### Fields

<dl>
<dt>
	<code>name</code>  <strong>string</strong>  - required, label 
</dt>

<dd>



</dd>

<dt>
	<code>vega_binary_path</code>  <strong>string</strong>  - optional
</dt>

<dd>

Allows the user to optionally use a different version of Vega binary for wallet

</dd>

<dt>
	<code>template</code>  <strong>string</strong>  - optional
</dt>

<dd>



</dd>



</dl>

---


## *FaucetConfig*


### Fields

<dl>
<dt>
	<code>name</code>  <strong>string</strong>  - required, label 
</dt>

<dd>



</dd>

<dt>
	<code>wallet_pass</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>template</code>  <strong>string</strong>  - optional
</dt>

<dd>



</dd>



</dl>

---


## *PStartConfig*


### Fields

<dl>
<dt>
	<code>docker_service</code>  <strong>[]<a href="#dockerconfig">DockerConfig</a></strong>  - required, block 
</dt>

<dd>



</dd>



</dl>

---


## *ConfigTemplates*

Allow to add configuration template for certain services deployed by Capsule.
Learn more about how configuration templating work here



### Fields

<dl>
<dt>
	<code>vega</code>  <strong>string</strong>  - required | optional if <code>vega_file</code> defined, optional 
</dt>

<dd>

[Go template](templates.md) of Vega config.

The [vega.ConfigTemplateContext](templates.md#vegaconfigtemplatecontext) can be used in the template. Example [example](net_confs/node_set_templates/default/vega_validators.tmpl).



<blockquote>It is recommended that you use `vega_file` param instead.
If both `vega` and `vega_file` are defined, then `vega`
overrides `vega_file`.
</blockquote>

<br />

#### <code>vega</code> example







```hcl
vega = <<EOH
 ...
EOH

```





</dd>

<dt>
	<code>vega_file</code>  <strong>string</strong>  - optional
</dt>

<dd>

Same as `vega` but it allows the user to link the Vega config template as an external file.



<br />

#### <code>vega_file</code> example







```hcl
vega_file = "/your_path/vega_config.tmpl"

```





</dd>

<dt>
	<code>tendermint</code>  <strong>string</strong>  - required | optional if <code>tendermint_file</code> defined, optional 
</dt>

<dd>

[Go template](templates.md) of Tendermint config.

The [tendermint.ConfigTemplateContext](templates.md#tendermintconfigtemplatecontext) can be used in the template. Example [example](net_confs/node_set_templates/default/tendermint_validators.tmpl).



<blockquote>It is recommended that you use `tendermint_file` param instead.
If both `tendermint` and `tendermint_file` are defined, then `tendermint`
overrides `tendermint_file`.
</blockquote>

<br />

#### <code>tendermint</code> example







```hcl
tendermint = <<EOH
 ...
EOH

```





</dd>

<dt>
	<code>tendermint_file</code>  <strong>string</strong>  - optional
</dt>

<dd>

Same as `tendermint` but it allows the user to link the Tendermint config template as an external file.



<br />

#### <code>tendermint_file</code> example







```hcl
tendermint_file = "/your_path/tendermint_config.tmpl"

```





</dd>

<dt>
	<code>data_node</code>  <strong>string</strong>  - required | optional if <code>data_node_file</code> defined, optional 
</dt>

<dd>

[Go template](templates.md) of Data Node config.

The [datanode.ConfigTemplateContext](templates.md#datanodeconfigtemplatecontext) can be used in the template. Example [example](net_confs/node_set_templates/default/data_node_full_external_postgres.tmpl).



<blockquote>It is recommended that you use `data_node_file` param instead.
If both `data_node` and `data_node_file` are defined, then `data_node`
overrides `data_node_file`.
</blockquote>
</dd>

<dt>
	<code>data_node_file</code>  <strong>string</strong>  - optional
</dt>

<dd>

Same as `data_node` but it allows the user to link the Data Node config template as an external file.


</dd>

<dt>
	<code>visor_run_conf</code>  <strong>string</strong>  - required | optional if <code>visor_run_conf_file</code> defined, optional 
</dt>

<dd>

[Go template](templates.md) of Visor genesis run config.

The [visor.ConfigTemplateContext](templates.md#visorconfigtemplatecontext) can be used in the template. Example [example](net_confs/node_set_templates/default/visor_run.tmpl).

Current Vega binary is automatically copied to the Visor genesis folder by Capsule
so it can be used from this template.



<blockquote>It is recommended that you use `visor_run_conf_file` param instead.
If both `visor_run_conf` and `visor_run_conf_file` are defined, then `visor_run_conf`
overrides `visor_run_conf_file`.
</blockquote>
</dd>

<dt>
	<code>visor_run_conf_file</code>  <strong>string</strong>  - optional
</dt>

<dd>

Same as `visor_run_conf` but it allows the user to link the Visor genesis run config template as an external file.


</dd>

<dt>
	<code>visor_conf</code>  <strong>string</strong>  - required | optional if <code>visor_conf_file</code> defined, optional 
</dt>

<dd>

[Go template](templates.md) of Visor config.

The [visor.ConfigTemplateContext](templates.md#visorconfigtemplatecontext) can be used in the template. Example [example](net_confs/node_set_templates/default/visor_config.tmpl).



<blockquote>It is recommended that you use `visor_conf_file` param instead.
If both `visor_conf` and `visor_conf_file` are defined, then `visor_conf`
overrides `visor_conf_file`.
</blockquote>
</dd>

<dt>
	<code>visor_conf_file</code>  <strong>string</strong>  - optional
</dt>

<dd>

Same as `visor_conf` but it allows the user to link the Visor genesis run config template as an external file.


</dd>



</dl>

---


## *PreGenerate*
Allows to define service that will run before generation step.


### Fields

<dl>
<dt>
	<code>nomad_job</code>  <strong>[]<a href="#nomadconfig">NomadConfig</a></strong>  - required, block 
</dt>

<dd>

Allows to define raw [Nomad jobs](https://developer.hashicorp.com/nomad/docs/job-specification).

</dd>



### Complete example



```hcl
pre_generate {
  nomad_job "clef" {
    ...
  }
}

```


</dl>

---


## *ClefConfig*

Allows to configure connetion to [Clef](https://geth.ethereum.org/docs/clef/introduction) Ethereum wallet.



### Fields

<dl>
<dt>
	<code>ethereum_account_addresses</code>  <strong>[]string</strong>  - required
</dt>

<dd>

List of Clef pre-generated Ethereum addresses that can be used by node set.



<blockquote>There should be enough available addresses for each node set.
So when node set has `count = 2` there has to be minimum 2 addresses defined
similarly when `count = 4` there has to be minimum 4 addresses defined etc.
</blockquote>
</dd>

<dt>
	<code>clef_rpc_address</code>  <strong>string</strong>  - required
</dt>

<dd>

Address of running Clef instance

</dd>



### Complete example



```hcl
clef_wallet {
  ethereum_account_addresses = ["0xc0ffee254729296a45a3885639AC7E10F9d54979", "0x999999cf1046e68e36E1aA2E0E07105eDDD1f08E"]
  clef_rpc_address           = "http://localhost:8555"
}

```


</dl>

---


## *DockerConfig*


### Fields

<dl>
<dt>
	<code>name</code>  <strong>string</strong>  - required, label 
</dt>

<dd>



</dd>

<dt>
	<code>image</code>  <strong>string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>cmd</code>  <strong>string</strong>  - optional
</dt>

<dd>



</dd>

<dt>
	<code>args</code>  <strong>[]string</strong>  - required
</dt>

<dd>



</dd>

<dt>
	<code>env</code>  <strong>map[string]string</strong>  - optional
</dt>

<dd>



</dd>

<dt>
	<code>static_port</code>  <strong><a href="#staticport">StaticPort</a></strong>  - optional, block 
</dt>

<dd>



</dd>

<dt>
	<code>auth_soft_fail</code>  <strong>bool</strong>  - optional
</dt>

<dd>



</dd>

<dt>
	<code>resources</code>  <strong><a href="#resources">Resources</a></strong>  - optional, block 
</dt>

<dd>



</dd>

<dt>
	<code>volume_mounts</code>  <strong>[]string</strong>  - optional
</dt>

<dd>



</dd>



</dl>

---


## *NomadConfig*

Allows to configure a [Nomad job](https://developer.hashicorp.com/nomad/docs/job-specification) definition to be run on Capsule.



### Fields

<dl>
<dt>
	<code>name</code>  <strong>string</strong>  - required, label 
</dt>

<dd>

Name of the Nomad job.


</dd>

<dt>
	<code>job_template</code>  <strong>string</strong>  - required | optional if <code>job_template_file</code> defined, optional 
</dt>

<dd>

[Go template](templates.md) of a Nomad job template.

The [nomad.PreGenerateTemplateCtx](templates.md#nomadpregeneratetemplatectx) can be used in the template. Example [example](jobs/clef.tmpl).



<blockquote>It is recommended that you use `job_template_file` param instead.
If both `job_template` and `job_template_file` are defined, then `job_template`
overrides `job_template_file`.
</blockquote>

<br />

#### <code>job_template</code> example







```hcl
job_template = <<EOH
 ...
EOH

```





</dd>

<dt>
	<code>job_template_file</code>  <strong>string</strong>  - optional
</dt>

<dd>

Same as `job_template` but it allows the user to link the Nomad job template as an external file.



<br />

#### <code>job_template_file</code> example







```hcl
job_template_file = "/your_path/nomad-job.tmpl"

```





</dd>



### Complete example



```hcl
nomad_job "clef" {
  job_template = "/path-to/nomad-job.tmpl"
}

```


</dl>

---


## *StaticPort*


### Fields

<dl>
<dt>
	<code>to</code>  <strong>int</strong>  - optional
</dt>

<dd>



</dd>

<dt>
	<code>value</code>  <strong>int</strong>  - required
</dt>

<dd>



</dd>



</dl>

---


## *Resources*


### Fields

<dl>
<dt>
	<code>cpu</code>  <strong>int</strong>  - optional
</dt>

<dd>



</dd>

<dt>
	<code>cores</code>  <strong>int</strong>  - optional
</dt>

<dd>



</dd>

<dt>
	<code>memory</code>  <strong>int</strong>  - optional
</dt>

<dd>



</dd>

<dt>
	<code>memory_max</code>  <strong>int</strong>  - optional
</dt>

<dd>



</dd>

<dt>
	<code>disk</code>  <strong>int</strong>  - optional
</dt>

<dd>



</dd>



</dl>

---


