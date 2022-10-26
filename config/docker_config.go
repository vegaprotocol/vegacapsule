package config

type DockerConfig struct {
	Name         string            `hcl:"name,label"`
	Image        string            `hcl:"image"`
	Command      string            `hcl:"cmd,optional"`
	Args         []string          `hcl:"args"`
	Env          map[string]string `hcl:"env,optional"`
	StaticPort   *StaticPort       `hcl:"static_port,block"`
	AuthSoftFail bool              `hcl:"auth_soft_fail,optional"`
	Resources    *Resources        `hcl:"resources,block"`
	VolumeMounts []string          `hcl:"volume_mounts,optional"`
}

type StaticPort struct {
	To    int `hcl:"to,optional"`
	Value int `hcl:"value"`
}

type Resources struct {
	CPU         *int `hcl:"cpu,optional"`
	Cores       *int `hcl:"cores,optional"`
	MemoryMB    *int `hcl:"memory,optional"`
	MemoryMaxMB *int `hcl:"memory_max,optional"`
	DiskMB      *int `hcl:"disk,optional"`
}
