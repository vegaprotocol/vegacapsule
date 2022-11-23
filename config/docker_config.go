package config

/*
description: |

	Allows to configure Docker container services that will run before or after the Vega network starts.

example:

	type: hcl
	value: |
			docker_service "ganache-1" {
				image = "vegaprotocol/ganache:latest"
				cmd = "ganache-cli"
				args = [
					"--blockTime", "1",
					"--chainId", "1440",
					"--networkId", "1441",
					"-h", "0.0.0.0",
				]
				static_port {
					value = 8545
					to = 8545
				}
				auth_soft_fail = true
			}
*/
type DockerConfig struct {
	/*
		description: Name of the service that is going to be use as an identifier when service runs.
		example:
			type: hcl
			value: |
					docker_service "service-name" {
						...
					}
	*/
	Name string `hcl:"name,label"`

	/*
		description: Name of publicly available Docker image.
		example:
			type: hcl
			value: |
					image = "vegaprotocol/ganache:latest"
	*/
	Image string `hcl:"image"`

	/*
		description: Command that will run at the image startup.
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
		description: Allows to set environment varibles for the container.
		example:
			type: hcl
			value: |
					env = {
						ENV_VAR="value"
						ENV_VAR_2="value-2"
					}
	*/
	Env map[string]string `hcl:"env,optional"`

	/*
		description: Allows to open a static port from container to host.
		example:
			type: hcl
			value: |
					static_port {
						value = 5232
						to = 5432
					}
	*/
	StaticPort *StaticPort `hcl:"static_port,block"`

	/*
		description: Defines whether or not the task fails on an auth failure.
		note: Should be always `true` for public images.
		example:
			type: hcl
			value: |
				auth_soft_fail = true
	*/
	AuthSoftFail bool `hcl:"auth_soft_fail,optional"`

	/*
		description: Allows to to define minimun required hardware resources for the container.
		note: In most cases the default values (not defined) should be sufficient.
		example:
			type: hcl
			value: |
					resources {
						cpu    = 100
						memory = 100
						memory_max = 300
					}
	*/
	Resources *Resources `hcl:"resources,block"`

	VolumeMounts []string `hcl:"volume_mounts,optional"`
}

/*
description: Represents static port mapping from host to container.
example:

	type: hcl
	value: |
			static_port {
				value = 8001
				to = 8002
			}
*/
type StaticPort struct {
	// description: Represents port value on the host.
	Value int `hcl:"value"`
	// description: Represents port value inside of the container.
	To int `hcl:"to,optional"`
}

/*
description: Allows to define hardware resoucers requirements
example:

	type: hcl
	value: |
			resources {
				cpu    = 100
				memory = 100
				memory_max = 300
			}
*/
type Resources struct {
	// description: Minimum required CPU in MHz
	CPU *int `hcl:"cpu,optional"`
	// description: Num of minimum required CPU cores
	Cores *int `hcl:"cores,optional"`
	// description: Minimum required RAM in Mb
	MemoryMB *int `hcl:"memory,optional"`
	// description: Maximum allowed RAM in Mb
	MemoryMaxMB *int `hcl:"memory_max,optional"`
	// description: Minimum required disk space in Mb
	DiskMB *int `hcl:"disk,optional"`
}
