package docsgenerator

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/fatih/structtag"
	yaml "gopkg.in/yaml.v3"
)

func tabToSpace(input string) string {
	return strings.ReplaceAll(input, "\t", " ")
}

type docTypeWithFileContent struct {
	*doc.Type
	fileContent string
	packageName string
}

func extractDocTypesFromFile(fileName, fileContent string) (map[string]docTypeWithFileContent, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, fileContent, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	p, err := doc.NewFromFiles(fset, []*ast.File{f}, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to parse docs from file: %w", err)
	}

	types := map[string]docTypeWithFileContent{}
	for _, t := range p.Types {
		types[t.Name] = docTypeWithFileContent{
			Type:        t,
			fileContent: fileContent,
			packageName: f.Name.Name,
		}
	}

	return types, nil
}

func parseComment(comment string) (*Comment, error) {
	c := &Comment{}
	if err := yaml.Unmarshal([]byte(tabToSpace(comment)), c); err != nil {
		return nil, fmt.Errorf("failed to parse yaml comment: %w", err)
	}

	return c, nil
}

func parseTag(tags, tagName string) (*structtag.Tag, error) {
	parsed, err := structtag.Parse(strings.ReplaceAll(tags, "`", ""))
	if err != nil {
		return nil, fmt.Errorf("failed to parse tags %q: %w", tags, err)
	}

	tag, err := parsed.Get(tagName)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag %q: %w", tagName, err)
	}

	return tag, nil
}

func findFiles(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}
