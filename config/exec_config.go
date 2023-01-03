package config

type ExecConfig struct {
	/*
		description: Name of the service that is going to be used as an identifier when service runs.
		example:
			type: hcl
			value: |
					docker_service "service-name" {
						...
					}
	*/
	Name string `hcl:"name,label"`

	/*
		description: Command that will run
		example:
			type: hcl
			value: |
					cmd = "ganache-cli"
	*/
	Command string `hcl:"cmd,optional"`

	/*
		description: List of arguments that will be added to cmd.
		example:
			type: hcl
			value: |
					args = [
						"--blockTime", "1",
				    	"--chainId", "1440",
					]
	*/
	Args []string `hcl:"args"`

	/*
		description: Allows the user to set environment variables launched process.
		example:
			type: hcl
			value: |
					env = {
						ENV_VAR="value"
						ENV_VAR_2="value-2"
					}
	*/
	Env map[string]string `hcl:"env,optional"`
}
