package ports

import (
	"fmt"

	"code.vegaprotocol.io/vega/paths"
)

const (
	PortConfigKey = "Port"
)

// ExtractPortsFromConfig read TOML config from disc and extracts
// all ports definitin into a map of ports to port name
func ExtractPortsFromConfig(configPath string) (map[int64]string, error) {
	config := map[string]interface{}{}
	if err := paths.ReadStructuredFile(configPath, &config); err != nil {
		return nil, fmt.Errorf("failed to read configuration file at %s: %w", configPath, err)
	}

	return ExtractPorts(config), nil
}

// ExtractPorts extracts all ports definitin into a map of ports to port name.
// Example map: { 2002: "API.Rest", 2003: "API.GRPC" }
func ExtractPorts(m map[string]interface{}) map[int64]string {
	out := map[int64]string{}

	for name, value := range m {
		switch v := value.(type) {
		case map[string]interface{}:
			inOut := ExtractPorts(v)

			for port, currentName := range inOut {
				newName := name
				if currentName != "" {
					newName = fmt.Sprintf("%s.%s", name, currentName)
				}

				out[port] = newName
			}
		case int64:
			if name == PortConfigKey {
				out[v] = ""
			}
		}
	}
	return out
}
