


# Capsule configuration docs

Capsule is a tool that allows users to run a custom Vega network simulation locally on single a machine.
This means that it is a incredibly useful tool for anybody who wants to try Vega network without using a real network.

Capsule configuration is used by vegacapsule CLI network to generate and bootstrap commands and can be customised to personal need.
Under the hood Capsule uses this configuration to generate a new network and stores all it's files in a single directory. This directory is then used by [Nomad](https://www.nomadproject.io/) to deploy all generated services from the generation step.

The configuration uses the [HCL](https://github.com/hashicorp/hcl) language syntax, which is also used, for example, by [Terraform](https://www.terraform.io/).

This document explains all possible configuration options in Capsule.


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

Same as `genesis_template` but it allows the user to link the genesis file