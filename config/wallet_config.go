package config

type WalletConfig struct {
	Name string `hcl:"name,label"`
	// description: Allows optionally use different version of Vega binary for wallet
	VegaBinary *string `hcl:"vega_binary_path,optional"`
	Template   string  `hcl:"template,optional"`
}
