package utils

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig"
)

func GenerateTemplate[T any](templateRaw string, tmplContext T) (*bytes.Buffer, error) {
	t, err := template.New("template.hcl").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template config for nomad job: %w", err)
	}

	buff := bytes.NewBuffer([]byte{})

	if err := t.Execute(buff, tmplContext); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buff, nil
}
