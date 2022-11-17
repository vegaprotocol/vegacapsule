package nomad

import (
	"bytes"
	"fmt"
	"text/template"

	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/Masterminds/sprig"
)

func GenerateNodeSetTemplate(templateRaw string, ns types.NodeSet) (*bytes.Buffer, error) {
	t, err := template.New("nomad_job.hcl").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template config for nomad job: %w", err)
	}

	buff := bytes.NewBuffer([]byte{})

	if err := t.Execute(buff, ns); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buff, nil
}

type PreGenerateTemplateCtx struct {
	Name          string
	Index         int
	LogsDir       string
	CapsuleBinary string
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
