package nomad

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig"
)

func GenerateNodeSetTemplate[T any](templateRaw string, context T) (*bytes.Buffer, error) {
	t, err := template.New("nomad_job.hcl").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template config for nomad job: %w", err)
	}

	buff := bytes.NewBuffer([]byte{})

	if err := t.Execute(buff, context); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buff, nil
}

type PreGenerateTemplateCtx struct {
	Name  string
	Index int
}

func GeneratePreGenerateTemplate(templateRaw string, ctx PreGenerateTemplateCtx) (*bytes.Buffer, error) {
	t, err := template.New("nomad_job.hcl").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template config for nomad job: %w", err)
	}

	buff := bytes.NewBuffer([]byte{})

	if err := t.Execute(buff, ctx); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buff, nil
}
