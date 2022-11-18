package config

/*
description: |

	Represents a configuration of a Vega Wallet service.

example:

	type: hcl
	value: |
		wallet "wallet-1" {
			template = <<-EOT
				...
			EOT

		}
*/
type WalletConfig struct {
	/*
		description: Name of the wallet that is going to be use as an identifier when wallet runs.
		example:
			type: hcl
			value: |
					wallet "wallet-name" {
						...
					}
	*/
	Name string `hcl:"name,label"`
	/*
		description: |
					By default wallet config inherits Vega binary from main network config but this paramater allows user to
					define a different Vega binary to be used in wallet.
					This can be used in case differnt wallet version is required.
					A relative or absolute path can be used. If only the binary name is defined, it automatically looks for it in $PATH.
		note: Using wallet version that is not compatible with network will not work - therefore this should be used in advanced cases only.
		example:
			type: hcl
			value: vega_binary_path = "binary_path"
	*/
	VegaBinary *string `hcl:"vega_binary_path,optional"`

	/*
		description: |
			[Go template](templates.md) of a Vega Wallet config.

			The [wallet.ConfigTemplateContext](templates.md#walletconfigtemplatecontext) can be used in the template.
			Example can be found in [default network config](net_confs/config.hcl).
		examples:
			- type: hcl
			  value: |
						template = <<EOH
							...
						EOH

	*/
	Template string `hcl:"template,optional"`
}
