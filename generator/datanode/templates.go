package datanode

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"

	"code.vegaprotocol.io/vega/datanode/networkhistory/store"
	"code.vegaprotocol.io/vega/paths"
	"code.vegaprotocol.io/vegacapsule/types"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig"
	"github.com/imdario/mergo"
)

type IPSFPeer struct {
	Index int
	ID    string
}

type ConfigTemplateContext struct {
	NodeHomeDir string
	NodeNumber  int
	NodeSet     types.NodeSet

	nodes []node
}

func (tc ConfigTemplateContext) getNetworkHistoryPeerIDSeed(nodeNumber int) string {
	return fmt.Sprintf("ipfs-seed-%d", nodeNumber)
}

func (tc ConfigTemplateContext) GetNetworkHistoryPeerID(nodeNumber int) string {
	seed := tc.getNetworkHistoryPeerIDSeed(nodeNumber)
	id, err := store.GenerateIdentityFromSeed([]byte(seed))
	if err != nil {
		panic("couldn't create ipfs peer identity")
	}
	return id.PeerID
}

func (tc ConfigTemplateContext) GetNetworkHistoryPrivKey(nodeNumber int) string {
	seed := tc.getNetworkHistoryPeerIDSeed(nodeNumber)
	id, err := store.GenerateIdentityFromSeed([]byte(seed))
	if err != nil {
		panic("couldn't create ipfs peer identity")
	}
	return id.PrivKey
}

func (tc ConfigTemplateContext) IPSFPeers() []IPSFPeer {
	peersIDs := []IPSFPeer{}
	for _, node := range tc.nodes {
		if len(tc.nodes) != 1 && node.index == tc.NodeSet.Index {
			continue
		}

		seed := tc.getNetworkHistoryPeerIDSeed(node.index)
		id, err := store.GenerateIdentityFromSeed([]byte(seed))
		if err != nil {
			panic("couldn't create ipfs peer identity")
		}

		peersIDs = append(peersIDs, IPSFPeer{
			Index: node.index,
			ID:    id.PeerID,
		})
	}

	return peersIDs
}

func NewConfigTemplate(templateRaw string) (*template.Template, error) {
	t, err := template.New("config.toml").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template config for data node: %w", err)
	}

	return t, nil
}

func (dng ConfigGenerator) TemplateConfig(ns types.NodeSet, configTemplate *template.Template) (*bytes.Buffer, error) {
	templateCtx := ConfigTemplateContext{
		NodeNumber:  ns.Index,
		NodeHomeDir: dng.homeDir,
		NodeSet:     ns,
		nodes:       dng.nodes,
	}

	buff := bytes.NewBuffer([]byte{})

	if err := configTemplate.Execute(buff, templateCtx); err != nil {
		return nil, fmt.Errorf("failed to execute template for data node: %w", err)
	}

	return buff, nil
}

func (dng *ConfigGenerator) TemplateAndMergeConfig(ns types.NodeSet, configTemplate *template.Template) (*bytes.Buffer, error) {
	tempFileName := fmt.Sprintf("datanode-%s.config", ns.Name)

	f, err := os.CreateTemp("", tempFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary config file %q: %w", tempFileName, err)
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	buff, err := dng.TemplateConfig(ns, configTemplate)
	if err != nil {
		return nil, err
	}

	// Sometimes the DataNode field may be nil. Especiall when you want to template the data-node config
	// with merge for the validator node.
	if ns.DataNode == nil {
		return nil, fmt.Errorf("failed to merge and save data node configuration: data node is not initialized properly")
	}

	if err := dng.mergeAndSaveConfig(ns, buff, originalConfigFilePath(ns.DataNode.HomeDir), f.Name()); err != nil {
		return nil, err
	}

	buffOut := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffOut, f); err != nil {
		return nil, fmt.Errorf("failed to copy content of config file %q: %w", f.Name(), err)
	}

	return buffOut, nil
}

func (dng ConfigGenerator) OverwriteConfig(ns types.NodeSet, configTemplate *template.Template) error {
	buff, err := dng.TemplateConfig(ns, configTemplate)
	if err != nil {
		return err
	}

	configFilePath := ConfigFilePath(ns.DataNode.HomeDir)
	return dng.mergeAndSaveConfig(ns, buff, configFilePath, configFilePath)
}

func (dng ConfigGenerator) mergeAndSaveConfig(
	ns types.NodeSet,
	tmpldConf *bytes.Buffer,
	configPath string,
	saveConfigPath string,
) error {
	overrideConfig := map[string]interface{}{}

	if _, err := toml.NewDecoder(tmpldConf).Decode(&overrideConfig); err != nil {
		return fmt.Errorf("failed decode override config: %w", err)
	}

	config := map[string]interface{}{}
	if err := paths.ReadStructuredFile(configPath, &config); err != nil {
		return fmt.Errorf("failed to read configuration file at %s: %w", configPath, err)
	}

	if err := mergo.Map(&config, overrideConfig, mergo.WithOverride); err != nil {
		return fmt.Errorf("failed to merge configs: %w", err)
	}

	if err := paths.WriteStructuredFile(saveConfigPath, config); err != nil {
		return fmt.Errorf("failed to write configuration file for data node at %s: %w", saveConfigPath, err)
	}

	return nil
}
