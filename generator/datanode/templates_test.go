package datanode

import (
	"fmt"
	"testing"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/zannen/toml"
)

func TestRun(t *testing.T) {
	tmpl, err := NewConfigTemplate(templateRaw)
	if err != nil {
		t.Error(err)
	}

	c, err := config.DefaultConfig()
	if err != nil {
		t.Error(err)
	}

	confGen, err := NewConfigGenerator(c)
	if err != nil {
		t.Error(err)
	}

	buff, err := confGen.TemplateConfig(types.NodeSet{}, tmpl)
	if err != nil {
		t.Error(err)
	}

	overrideConfig := map[string]interface{}{}

	if _, err := toml.DecodeReader(buff, &overrideConfig); err != nil {
		t.Error(err)
	}

	fmt.Println(ExtractPorts(overrideConfig))
}

func ExtractPorts(m map[string]interface{}) map[int64]string {
	out := map[int64]string{}

	for name, value := range m {
		switch v := value.(type) {
		case map[string]interface{}:
			walkOut := ExtractPorts(v)

			for port, currentName := range walkOut {
				newName := name
				if currentName != "" {
					newName = fmt.Sprintf("%s.%s", name, currentName)
				}

				out[port] = newName
			}
		case int64:
			if name == "Port" {
				out[v] = ""
			}
		}
	}

	return out
}

const templateRaw = `
GatewayEnabled = true

[SQLStore]
  Enabled = true
  [SQLStore.ConnectionConfig]
    Port = 5232
    UseTransactions = true
    Database = "vega{{.NodeNumber}}"
  

[API]
  Level = "Info"
  Port = 30{{.NodeNumber}}7
  CoreNodeGRPCPort = 30{{.NodeNumber}}2

[Pprof]
  Level = "Info"
  Enabled = false
  Port = 6{{.NodeNumber}}60
  ProfilesDir = "{{.NodeHomeDir}}"

[Gateway]
  Level = "Info"
  [Gateway.Node]
    Port = 30{{.NodeNumber}}7
  [Gateway.GraphQL]
    Port = 30{{.NodeNumber}}8
  [Gateway.REST]
    Port = 30{{.NodeNumber}}9
	
[Metrics]
  Level = "Info"
  Timeout = "5s"
  Port = 21{{.NodeNumber}}2
  Enabled = false
[Broker]
  Level = "Info"
  UseEventFile = false
  [Broker.SocketConfig]
    Port = 30{{.NodeNumber}}5
`
