package config

/*
description: |

	Represents a configuration of a Vega Faucet service.

example:

	type: hcl
	value: |
		faucet "faucet-1" { {
			wallet_pass = "wallet_pass"
			template = <<-EOT
				...
			EOT
		}
*/
type FaucetConfig struct {
	/*
		description: Name of the faucet that is going to be use as an identifier when facet runs.
		example:
			type: hcl
			value: |
					wallet "wallet-name" {
						...
					}
	*/
	Name string `hcl:"name,label"`

	/*
		description: Passphrase for the wallet.
		example:
			type: hcl
			value: wallet_pass = "passphrase"
	*/
	Pass string `hcl:"wallet_pass"`

	/*
		description: |
			[Go template](templates.md) of a Vega Faucet config.

			The [faucet.ConfigTemplateContext](templates.md#faucetconfigtemplatecontext) can be used in the template.
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
