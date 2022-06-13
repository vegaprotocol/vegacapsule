package tendermint

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"

	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/Masterminds/sprig"
	"github.com/spf13/viper"
	tmconfig "github.com/tendermint/tendermint/config"
)

type ConfigTemplateContext struct {
	Prefix               string
	TendermintNodePrefix string
	VegaNodePrefix       string
	NodeNumber           int
	NodesCount           int
	NodeSet              types.NodeSet
	nodes                []node
}

func NewConfigTemplate(templateRaw string) (*template.Template, error) {
	t, err := template.New("config.toml").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template config: %w", err)
	}

	return t, nil
}

func nodesByGroupNameLookup(nodes []node) map[string]struct{} {
	m := map[string]struct{}{}
	for _, n := range nodes {
		m[n.groupName] = struct{}{}
	}
	return m
}

func (tc ConfigTemplateContext) NodePeersByGroupName(groupNames ...string) []Peer {
	gns := nodesByGroupNameLookup(tc.nodes)

	peers := []Peer{}
	for _, node := range tc.nodes {
		if node.index == tc.NodeSet.Index {
			continue
		}

		if _, ok := gns[node.groupName]; !ok {
			continue
		}

		peers = append(peers, Peer{
			Index: node.index,
			ID:    node.id,
		})
	}

	return peers
}

func (tc ConfigTemplateContext) NodePeers() []Peer {
	peers := []Peer{}
	for _, node := range tc.nodes {
		if node.index == tc.NodeSet.Index {
			continue
		}

		peers = append(peers, Peer{
			Index: node.index,
			ID:    node.id,
		})
	}

	return peers
}

func (tc ConfigTemplateContext) NodeIDsByGroupName(groupNames ...string) []string {
	gns := nodesByGroupNameLookup(tc.nodes)

	nodeIDs := make([]string, 0, len(tc.nodes))
	for _, node := range tc.nodes {
		if _, ok := gns[node.groupName]; !ok {
			continue
		}

		nodeIDs = append(nodeIDs, node.id)
	}

	return nodeIDs
}

func (tc ConfigTemplateContext) NodeIDs() []string {
	nodeIDs := make([]string, 0, len(tc.nodes))

	for _, node := range tc.nodes {
		nodeIDs = append(nodeIDs, node.id)
	}

	return nodeIDs
}

// TemplateConfig templates the provided template
func (tg *ConfigGenerator) TemplateConfig(ns types.NodeSet, configTemplate *template.Template) (*bytes.Buffer, error) {
	templateCtx := ConfigTemplateContext{
		Prefix:               *tg.conf.Prefix,
		TendermintNodePrefix: *tg.conf.TendermintNodePrefix,
		VegaNodePrefix:       *tg.conf.VegaNodePrefix,
		NodeNumber:           ns.Index,
		NodesCount:           len(tg.nodes),
		NodeSet:              ns,
		nodes:                tg.nodes,
	}

	buff := bytes.NewBuffer([]byte{})

	if err := configTemplate.Execute(buff, templateCtx); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buff, nil
}

// TemplateAndMergeConfig templates provided template and merge it with originally initated Tendermint instance's config
func (tg *ConfigGenerator) TemplateAndMergeConfig(ns types.NodeSet, configTemplate *template.Template) (*bytes.Buffer, error) {
	tempFileName := fmt.Sprintf("tendermint-%s.config", ns.Name)

	f, err := os.CreateTemp("", tempFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary config file %q: %w", tempFileName, err)
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	buff, err := tg.TemplateConfig(ns, configTemplate)
	if err != nil {
		return nil, err
	}

	if err := tg.mergeAndSaveConfig(ns, buff, originalConfigFilePath(ns.Tendermint.HomeDir), f.Name()); err != nil {
		return nil, err
	}

	buffOut := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffOut, f); err != nil {
		return nil, fmt.Errorf("failed to copy content of config file %q: %w", f.Name(), err)
	}

	return buffOut, nil
}

func (tg *ConfigGenerator) OverwriteConfig(ns types.NodeSet, configTemplate *template.Template) error {
	buff, err := tg.TemplateConfig(ns, configTemplate)
	if err != nil {
		return err
	}

	configFilePath := ConfigFilePath(ns.Tendermint.HomeDir)
	return tg.mergeAndSaveConfig(ns, buff, configFilePath, configFilePath)
}

func (tg *ConfigGenerator) mergeAndSaveConfig(
	ns types.NodeSet,
	tmpldConf *bytes.Buffer,
	configPath string,
	saveConfigPath string,
) error {
	// merge
	v := viper.New()
	v.SetConfigFile(configPath)
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file %q: %w", configPath, err)
	}

	if err := v.MergeConfig(tmpldConf); err != nil {
		return fmt.Errorf("failed to merge config override with config file %q: %w", configPath, err)
	}
	conf := tmconfig.DefaultConfig()
	if err := v.Unmarshal(conf); err != nil {
		return fmt.Errorf("failed to unmarshal merged config file %q: %w", configPath, err)
	}

	if err := conf.ValidateBasic(); err != nil {
		return fmt.Errorf("failed to validated merged config file %q: %w", configPath, err)
	}

	// save
	conf.SetRoot(ns.Tendermint.HomeDir)
	tmconfig.WriteConfigFile(saveConfigPath, conf)

	return nil
}
