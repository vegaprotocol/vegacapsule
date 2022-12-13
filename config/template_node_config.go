package config

import (
	"bytes"
	"reflect"
)

type NodeConfigTemplateContext struct {
	// description: Absolute index of the generated node.
	NodeNumber int
}

func TemplateNodeConfig(templateContext NodeConfigTemplateContext, n NodeConfig) (*NodeConfig, error) {
	tmplFunc := func(templateRaw string) (*bytes.Buffer, error) {
		return executeConfigTemplate(templateRaw, templateContext)
	}

	if err := TemplateStruct(reflect.ValueOf(&n), tmplFunc); err != nil {
		return nil, err
	}

	return &n, nil
}
