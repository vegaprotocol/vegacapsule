package nomad

import (
	"bytes"
	"fmt"
	"text/template"

	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/Masterminds/sprig"
)

func GenerateTemplate(templateRaw string, ns types.NodeSet) (string, error) {
	t, err := template.New("nomad_job.hcl").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return "", fmt.Errorf("failed to parse template config for nomad job: %w", err)
	}

	buff := bytes.NewBuffer([]byte{})

	if err := t.Execute(buff, ns); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buff.String(), nil
}
