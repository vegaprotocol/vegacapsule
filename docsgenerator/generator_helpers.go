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
	// nolint
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

// formatLookupKey formats strings to "packageName.TypeName"
func formatLookupKey(packageName, typeName string) string {
	return fmt.Sprintf("%s.%s", packageName, typeName)
}

type fieldInfo struct {
	lookupKey  string
	isOptional bool
}

func getFieldInfo(
	currentPackageName, fieldType string,
	field *ast.Field,
	comment *Comment,
	isOptional bool,
) fieldInfo {
	fi := fieldInfo{
		isOptional: isOptional,
	}

	packageName := currentPackageName

	typeSplit := strings.Split(fieldType, ".")
	if len(typeSplit) > 1 {
		packageName = typeSplit[0]
		fi.lookupKey = typeSplit[1]
	} else {
		fi.lookupKey = typeSplit[0]
	}

	switch field.Type.(type) {
	case *ast.ArrayType:
		fi.lookupKey = strings.TrimLeft(fieldType, "[]")
	case *ast.MapType:
		fi.lookupKey = valueTypeFromMap(fieldType)
	case *ast.Ident:
		fi.lookupKey = fieldType
	case *ast.StarExpr:
		fi.lookupKey = strings.TrimLeft(fi.lookupKey, "*")
		if comment.OptionalIf == "" {
			fi.isOptional = true
		}
	}

	fi.lookupKey = formatLookupKey(packageName, fi.lookupKey)

	return fi
}

func typesFromMap(mapType string) (string, string) {
	sm := mapRegex.FindStringSubmatch(mapType)
	return sm[1], sm[2]
}

func valueTypeFromMap(mapType string) string {
	_, val := typesFromMap(mapType)
	return val
}

func getFieldType(fileContent string, field *ast.Field) string {
	typeExpr := field.Type

	start := typeExpr.Pos() - 1
	end := typeExpr.End() - 1

	// grab it in source
	return fileContent[start:end]
}
