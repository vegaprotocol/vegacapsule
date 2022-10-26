package docsgenerator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

var markdownTemplate = `
{{ define "fieldExamples" }}
Examples:

{{ range $example := .Examples }}
{{ if $example.Name }}**{{ $example.Name }}**{{ end }}

{{ if eq $example.Type "hcl" }}
{{ $formated_example := formatHCL $example.Value }}
{{ codeBlock $formated_example }}
{{ else }}
{{ codeBlock $example.Value }}
{{ end }}

{{ end }}

{{ end }}

# {{ .Name }}
{{ .Description }}

{{ range $type := .Types }}
### {{ if $type.Name }}{{ $type.Name }} - {{ end }}*{{ $type.Type }}*
{{ if $type.Description -}}
{{ $type.Description }}
{{ end }}

{{ if $type.Fields -}}
**Fields**:
<hr />

{{ range $field := $type.Fields -}}
<div class="dd">

<code>{{ $field.Name }}</code>  *{{ encodeType $field.Type }}* {{ if $field.Optional}} - optional{{else}} - required{{end}}{{ range $opt := $field.Options }}, {{ $opt }} {{ end }}

</div>
<div class="dt">

{{ $field.Description }}

{{ if $field.Values }}
Valid values:

{{ range $value := $field.Values }}
  - <code>{{ $value }}</code>
{{ end -}}
{{ end -}}

{{ if $field.Default }}
Default value: <code>{{ $field.Default }}</code>
{{ end }}

{{- if $field.Note }}
> {{ $field.Note }}
{{ end -}}

{{- if $field.Examples }}
{{ template "fieldExamples" $field }}
{{ end -}}

</div>

<hr />

{{ end }}

{{ if .Example.Value -}}
**Example**:

{{ if eq .Example.Type "hcl" }}
{{ $formated_example := formatHCL .Example.Value }}
{{ codeBlock $formated_example }}
{{ else }}
{{ codeBlock .Example.Value }}
{{ end }}

{{ end -}}

{{ end -}}

{{ end }}`

// FileDoc represents a single go file documentation.
type FileDoc struct {
	// Name will be used in md file name pattern.
	Name string
	// Description file description, supports markdown.
	Description string
	// Types structs defined in the file.
	Types   []*TypeDoc
	Anchors map[string]string

	t *template.Template
}

func NewFileDoc(name, description string, types []*TypeDoc) *FileDoc {
	anchors := map[string]string{}
	for _, t := range types {
		anchors[t.Type] = strings.ToLower(t.Type)
	}

	fd := &FileDoc{
		Name:        name,
		Description: description,
		Anchors:     anchors,
		Types:       types,
	}

	fd.t = template.Must(template.New("docs_markdown.tpl").
		Funcs(template.FuncMap{
			"codeBlock":  codeBlock,
			"encodeType": fd.encodeType,
			"formatHCL":  formatHCL,
		}).
		Parse(markdownTemplate))

	return fd
}

// Encode encodes file documentation as MD file.
func (fd *FileDoc) Encode() ([]byte, error) {
	buf := bytes.Buffer{}

	if err := fd.t.Execute(&buf, fd); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Write dumps documentation string to folder.
func (fd *FileDoc) Write(path, frontmatter string) error {
	data, err := fd.Encode()
	if err != nil {
		return err
	}

	if stat, e := os.Stat(path); !os.IsNotExist(e) {
		if !stat.IsDir() {
			return fmt.Errorf("destination path should be a directory")
		}
	} else {
		if e := os.MkdirAll(path, 0o777); e != nil {
			return e
		}
	}

	f, err := os.Create(filepath.Join(path, fmt.Sprintf("%s.%s", strings.ToLower(fd.Name), "md")))
	if err != nil {
		return err
	}

	if _, err := f.Write([]byte(frontmatter)); err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

var re = regexp.MustCompile(`[A-Za-z\.]+`)

func (fd *FileDoc) encodeType(t string) string {
	for _, s := range re.FindAllString(t, -1) {
		if anchor, ok := fd.Anchors[s]; ok {
			t = strings.ReplaceAll(t, s, formatLink(s, "#"+strings.ReplaceAll(anchor, ".", "")))
		}
	}
	return t
}

func codeBlock(text string) string {
	return "```\n" + text + "\n```"
}

func formatLink(text, link string) string {
	return fmt.Sprintf(`<a href="%s">%s</a>`, link, text)
}

func formatHCL(text string) string {
	return string(hclwrite.Format([]byte(text)))
}
