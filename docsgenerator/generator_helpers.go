package docsgenerator

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"

	"github.com/fatih/structtag"
	yaml "gopkg.in/yaml.v3"
)

func tabToSpace(input string) string {
	return strings.ReplaceAll(input, "\t", " ")
}

func extractTypesFromFile(fileName, fileContent string) (map[string]*doc.Type, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, fileContent, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	p, err := doc.NewFromFiles(fset, []*ast.File{f}, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to parse docs from file: %w", err)
	}

	types := map[string]*doc.Type{}
	for _, t := range p.Types {
		types[t.Name] = t
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
