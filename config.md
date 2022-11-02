


# Capsule configuration docs
Capsule configuration is used by vegacapsule CLI network generate and bootstrap commands.
It allows to configue and customise Vega network running on Capsule.

The configuration is using [HCL](https://github.com/hashicorp/hcl) language syntax also used for example by [Terraform](https://www.terraform.io/).

This document explains all possible configuration options in Capsule.



### Root - *Config*

All paramaters from this types are use directly in the config file.
Most of the paramaters here are optional and can be left alone.
Please see the example below.



**Fields**:
<hr />

<div class="dd">

<code>network</code>  *<a href="#networkconfig">NetworkConfig</a>*  - required, block 

</div>
<div class="dt">

Configuration of Vega network and it's dependencies.

</div>

<hr />

<div class="dd">

<code>output_dir</code>  *string*  - optional

</div>
<div class="dt">

Directory path (relative or absolute) where Capsule stores generated folders, files, logs and configurations for network.



Default value: <code>~/.vegacapsule/testnet</code>
</div>

<hr />

<div class="dd">

<code>vega_binary_path</code>  *string*  - optional

</div>
<div class="dt">

Path (relative or absolute) to vega binary that will be used to generate and run the network.


Default value: <code>vega</code>
</div>

<hr />

<div class="dd">

<code>vega_capsule_binary_path</code>  *string*  - optional

</div>
<div class="dt">

Path (relative or absolute) of a Capsule binary. The Capsule binary is used by Nomad to aggregate logs from running jobs
and save them to local disk in Capsule home directory.
See `vegacapsule nomad logscollector` for more info.



Default value: <code>Currently running Capsule instance binary</code>

> This optional paramater is used internally. There should never be need to set it to anything else then default.
</div>

<hr />



**Example**:



```
vega_binary_path = "/path/to/vega"

network "your_network_name" {
  ...
}

```



### *NetworkConfig*

Network configuration allows to customise Vega network into different shapes based on personal needs.
It allows to configure and deploy different Vega nodes setups (validator, full) and their dependencies (like Ethereum or Postgres).
It can run custom Docker images before and after the network nodes has started and much more.



**Fields**:
<hr />

<div class="dd">

<code>name</code>  *string*  - required, label 

</div>
<div class="dt">

Name of the network.
All folders generated are placed in folder with this name.
All Nomad jobs are prefix with the name.


</div>

<hr />

<div class="dd">

<code>genesis_template</code>  *string*  - required | optional if <code>genesis_template_file</code> defined)

</div>
<div class="dt">

Go template of genesis file that will be used to bootrap the Vega network.
[Example of templated mainnet genesis file](https://github.com/vegaprotocol/networks/blob/master/mainnet1/genesis.json)



> It is recomended to use `genesis_template_file` param instead.
In case both `genesis_template` and `genesis_template_file` are defined the `genesis_template`
overrides `genesis_template_file`.



Examples:






```
genesis_template = <<EOH
 {
  "app_state": {
   ...
  }
  ..
 }
EOH

```





</div>

<hr />

<div class="dd">

<code>genesis_template_file</code>  *string*  - optional

</div>
<div class="dt">

Same as `genesis_template` but it allows to link the genesis file template as an external file.




Examples:






```
genesis_template_file = "/your_path/genesis.tmpl"

```





</div>

<hr />

<div class="dd">

<code>ethereum</code>  *<a href="#ethereumconfig">EthereumConfig</a>*  - required, block 

</div>
<div class="dt">

Allows to define Ethereum network configuration.
This is necessery because Vega needs to be connected to [Ethereum bridges](https://docs.vega.xyz/mainnet/api/bridge)
or it cannot function otherwise.




Examples:






```
ethereum {
  ...
}

```





</div>

<hr />

<div class="dd">

<code>smart_contracts_addresses</code>  *string*  - required | optional if <code>smart_contracts_addresses_file</code> defined), optional 

</div>
<div class="dt">

Smart contract addresses are addresses of [Ethereum bridges](https://docs.vega.xyz/mainnet/api/bridge) contracts in JSON format.

These addresses should correspond to the choosen network by [Ethereum network](#EthereumConfig) and
can be used in various different types of templates in Capsule.
[Example of smart contract address from mainnet](https://github.com/vegaprotocol/networks/blob/master/mainnet1/smart-contracts.json).



> It is recomended to use `smart_contracts_addresses_file` param instead.
In case both `smart_contracts_addresses` and `smart_contracts_addresses_file` are defined the `genesis_template`
overrides `smart_contracts_addresses_file`.



Examples:






```
smart_contracts_addresses = <<EOH
 {
  "erc20_bridge": "...",
  "staking_bridge": "...",
  ...
 }
EOH

```





</div>

<hr />

<div class="dd">

<code>smart_contracts_addresses_file</code>  *string*  - optional

</div>
<div class="dt">

Same as `smart_contracts_addresses` but it allows to link the smart contracts as an external file.




Examples:






```
smart_contracts_addresses_file = "/your_path/smart-contratcs.json"

```





</div>

<hr />

<div class="dd">

<code>node_set</code>  *list(<a href="#nodeconfig">NodeConfig</a>)*  - required, block 

</div>
<div class="dt">

Allows to define multiple nodes set and their specific configuration.
A node set is a representation of Vega and Data Node nodes.
This is building unit of the Vega network.




Examples:


**Validators node set**



```
node_set "validator-nodes" {
  ...
}

```



**Full nodes node set**



```
node_set "full-nodes" {
  ...
}

```





</div>

<hr />

<div class="dd">

<code>wallet</code>  *<a href="#walletconfig">WalletConfig</a>*  - optional, block 

</div>
<div class="dt">

Allows to deploy and configure [Vega Wallet](https://docs.vega.xyz/mainnet/tools/vega-wallet) instance.
Wallet will not be deployed if this block is not defined.




Examples:






```
wallet "wallet-name" {
  ...
}

```





</div>

<hr />

<div class="dd">

<code>faucet</code>  *<a href="#faucetconfig">FaucetConfig</a>*  - optional, block 

</div>
<div class="dt">

Allows to deploy and configure [Vega Core Faucet](https://github.com/vegaprotocol/vega/tree/develop/core/faucet#faucet) instance.
Faucet will not be deployed if this block is not defined.




Examples:






```
faucet "faucet-name" {
  ...
}

```





</div>

<hr />

<div class="dd">

<code>pre_start</code>  *<a href="#pstartconfig">PStartConfig</a>*  - optional, block 

</div>
<div class="dt">

Allows to define jobs that should run before the node sets starts.
It can be used for node sets dependencies like databases or mock Ethereum chain etc..




Examples:






```
pre_start {
  docker_service "ganache-1" {
    ...
  }
  docker_service "postgres-1" {
    ...
  }
}

```





</div>

<hr />

<div class="dd">

<code>post_start</code>  *<a href="#pstartconfig">PStartConfig</a>*  - optional, block 

</div>
<div class="dt">

Allows to define jobs that should run after the node sets started.
It can be used for services that depends not already running network like block explorer or console.




Examples:






```
post_start {
  docker_service "bloc-explorer-1" {
    ...
  }
  docker_service "vega-console-1" {
    ...
  }
}

```





</div>

<hr />



**Example**:



```
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



### *EthereumConfig*

Allows to define specific Ethereum network to be used.
It can either some of the [Public networks](https://ethereum.org/en/developers/docs/networks/#public-networks) or
local instance of Ganache.



**Fields**:
<hr />

<div class="dd">

<code>chain_id</code>  *string*  - required

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>network_id</code>  *string*  - required

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>endpoint</code>  *string*  - required

</div>
<div class="dt">



</div>

<hr />



**Example**:



```
ethereum {
  chain_id   = "1440"
  network_id = "1441"
  endpoint   = "http://127.0.0.1:8545/"
}

```



### *NodeConfig*

Represents and allows to configure set of Vega (with Tendermint) and Data Node nodes.
One node set definition can be used by applied to multiple node sets (see `count` field) and it uses
templating to distinguish between different nodes and names/ports and other collisions.



**Fields**:
<hr />

<div class="dd">

<code>name</code>  *string*  - required, label 

</div>
<div class="dt">

Name of the node set.
Nomad that are part of this nodes are prefix with the name.


</div>

<hr />

<div class="dd">

<code>mode</code>  *string*  - required

</div>
<div class="dt">

Determines what mode the node set should run in.



Valid values:


  - <code>validator</code>

  - <code>full</code>
</div>

<hr />

<div class="dd">

<code>count</code>  *int*  - required

</div>
<div class="dt">

Defines how many nodes sets with this exact configuration should be created.


</div>

<hr />

<div class="dd">

<code>node_wallet_pass</code>  *string*  - optional | required if <code>mode=validator</code> defined)

</div>
<div class="dt">

Defines password for automatically generated node wallet assosiated with the created node.

</div>

<hr />

<div class="dd">

<code>ethereum_wallet_pass</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>vega_wallet_pass</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>vega_binary_path</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>use_data_node</code>  *bool*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>visor_binary</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>config_templates</code>  *<a href="#configtemplates">ConfigTemplates</a>*  - required, block 

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>pre_generate</code>  *<a href="#pregenerate">PreGenerate</a>*  - optional, block 

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>pre_start_probe</code>  *types.ProbesConfig*  - optional, block 

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>clef_wallet</code>  *<a href="#clefconfig">ClefConfig</a>*  - optional, block 

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>nomad_job_template</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>nomad_job_template_file</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />



**Example**:



```
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



### *WalletConfig*


**Fields**:
<hr />

<div class="dd">

<code>name</code>  *string*  - required, label 

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>vega_binary_path</code>  *string*  - optional

</div>
<div class="dt">

Allows optionally use different version of Vega binary for wallet

</div>

<hr />

<div class="dd">

<code>template</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />




### *FaucetConfig*


**Fields**:
<hr />

<div class="dd">

<code>name</code>  *string*  - required, label 

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>wallet_pass</code>  *string*  - required

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>template</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />




### *PStartConfig*


**Fields**:
<hr />

<div class="dd">

<code>docker_service</code>  *list(<a href="#dockerconfig">DockerConfig</a>)*  - required, block 

</div>
<div class="dt">



</div>

<hr />




### *ConfigTemplates*


**Fields**:
<hr />

<div class="dd">

<code>vega</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>vega_file</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>tendermint</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>tendermint_file</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>data_node</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>data_node_file</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>visor_run_conf</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>visor_run_conf_file</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>visor_conf</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>visor_conf_file</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />




### *PreGenerate*


**Fields**:
<hr />

<div class="dd">

<code>nomad_job</code>  *list(<a href="#nomadconfig">NomadConfig</a>)*  - required, block 

</div>
<div class="dt">



</div>

<hr />




### *ClefConfig*


**Fields**:
<hr />

<div class="dd">

<code>ethereum_account_addresses</code>  *list(string)*  - required

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>clef_rpc_address</code>  *string*  - required

</div>
<div class="dt">



</div>

<hr />




### *DockerConfig*


**Fields**:
<hr />

<div class="dd">

<code>name</code>  *string*  - required, label 

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>image</code>  *string*  - required

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>cmd</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>args</code>  *list(string)*  - required

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>env</code>  *map[string]string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>static_port</code>  *<a href="#staticport">StaticPort</a>*  - optional, block 

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>auth_soft_fail</code>  *bool*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>resources</code>  *<a href="#resources">Resources</a>*  - optional, block 

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>volume_mounts</code>  *list(string)*  - optional

</div>
<div class="dt">



</div>

<hr />




### *NomadConfig*


**Fields**:
<hr />

<div class="dd">

<code>name</code>  *string*  - required, label 

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>job_template</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>job_template_file</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />




### *StaticPort*


**Fields**:
<hr />

<div class="dd">

<code>to</code>  *int*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>value</code>  *int*  - required

</div>
<div class="dt">



</div>

<hr />




### *Resources*


**Fields**:
<hr />

<div class="dd">

<code>cpu</code>  *int*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>cores</code>  *int*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>memory</code>  *int*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>memory_max</code>  *int*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>disk</code>  *int*  - optional

</div>
<div class="dt">



</div>

<hr />




