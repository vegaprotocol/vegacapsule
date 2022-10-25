


This Capsule configuration file allows to configurate custom Vega network.
## Config


### Fields:
<hr />

<div class="dd">

<code>output_dir</code>  *string*  - optional

</div>
<div class="dt">

OutputDir is customisable field

</div>

<hr />

<div class="dd">

<code>vega_binary_path</code>  *string*  - optional

</div>
<div class="dt">

VegaBinary is customisable field

</div>

<hr />

<div class="dd">

<code>vega_capsule_binary_path</code>  *string*  - optional

</div>
<div class="dt">

VegaCapsuleBinary is customisable field

</div>

<hr />

<div class="dd">

<code>prefix</code>  *string*  - optional

</div>
<div class="dt">

Prefix is customisable field

</div>

<hr />

<div class="dd">

<code>node_dir_prefix</code>  *string*  - optional

</div>
<div class="dt">

NodeDirPrefix is customisable field

</div>

<hr />

<div class="dd">

<code>tendermint_node_prefix</code>  *string*  - optional

</div>
<div class="dt">

TendermintNodePrefix is customisable field

</div>

<hr />

<div class="dd">

<code>vega_node_prefix</code>  *string*  - optional

</div>
<div class="dt">

VegaNodePrefix is customisable field

</div>

<hr />

<div class="dd">

<code>data_node_prefix</code>  *string*  - optional

</div>
<div class="dt">

DataNodePrefix is customisable field

</div>

<hr />

<div class="dd">

<code>wallet_prefix</code>  *string*  - optional

</div>
<div class="dt">

WalletPrefix is customisable field

</div>

<hr />

<div class="dd">

<code>faucet_prefix</code>  *string*  - optional

</div>
<div class="dt">

FaucetPrefix is customisable field

</div>

<hr />

<div class="dd">

<code>visor_prefix</code>  *string*  - optional

</div>
<div class="dt">

VisorPrefix is customisable field

</div>

<hr />

<div class="dd">

<code>network</code>  *<a href="#networkconfig">NetworkConfig</a>*  - required, block 

</div>
<div class="dt">

Network is customisable field

</div>

<hr />




## NetworkConfig


### Fields:
<hr />

<div class="dd">

<code>name</code>  *string*  - required, label 

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>genesis_template</code>  *string*  - optional

</div>
<div class="dt">



</div>

<hr />

<div class="dd">

<code>genesis_template_file</code>  *string*  - optional

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

<div class="dd">

<code>node_set</code>  *list(<a href="#nodeconfig">NodeConfig</a>)*  - required, block 

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




## EthereumConfig


### Fields:
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




## WalletConfig


### Fields:
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




## FaucetConfig


### Fields:
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




## PStartConfig


### Fields:
<hr />

<div class="dd">

<code>docker_service</code>  *list(<a href="#dockerconfig">DockerConfig</a>)*  - required, block 

</div>
<div class="dt">



</div>

<hr />




## NodeConfig


### Fields:
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




## DockerConfig


### Fields:
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




## PreGenerate


### Fields:
<hr />

<div class="dd">

<code>nomad_job</code>  *list(<a href="#nomadconfig">NomadConfig</a>)*  - required, block 

</div>
<div class="dt">



</div>

<hr />




## ClefConfig


### Fields:
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




## ConfigTemplates


### Fields:
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




## StaticPort


### Fields:
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




## Resources


### Fields:
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




## NomadConfig


### Fields:
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




