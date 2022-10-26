


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

<code>genesis_template</code>  *string*  - optional

</div>
<div class="dt">

Template of genesis file that will be used to bootrap the Vega network.



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



</div>

<hr />

<div class="dd">

<code>smart_contracts_addresses</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>smart_contracts_addresses_file</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>ethereum</code>  *<a href="#ethereumconfig">EthereumConfig</a>*  - required, block 

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>node_set</code>  *list(<a href="#nodeconfig">NodeConfig</a>)*  - required, block 

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>wallet</code>  *<a href="#walletconfig">WalletConfig</a>*  - optional, block 

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>faucet</code>  *<a href="#faucetconfig">FaucetConfig</a>*  - optional, block 

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>pre_start</code>  *<a href="#pstartconfig">PStartConfig</a>*  - optional, block 

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>post_start</code>  *<a href="#pstartconfig">PStartConfig</a>*  - optional, block 

</div>
<div class="dt">



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




### *NodeConfig*


**Fields**:
<hr />

<div class="dd">

<code>name</code>  *string*  - required, label 

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>mode</code>  *string*  - required

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>count</code>  *int*  - required

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>node_wallet_pass</code>  *string*  - optional

</div>
<div class="dt">



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




