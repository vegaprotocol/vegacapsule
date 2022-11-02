package docsgenerator

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"strings"
)

const optionalTag = "optional"

type TypeDocGenerator struct {
	fileTypes map[string]docTypeWithFileContent

	tagName string
}

func NewTypeDocGenerator(dir, tagName string) (*TypeDocGenerator, error) {
	match := fmt.Sprintf(`%s/*.go`, dir)

	filePaths, err := filepath.Glob(match)
	if err != nil {
		return nil, fmt.Errorf("failed to look for files in dir %q: %w", dir, err)
	}

	types := map[string]docTypeWithFileContent{}

	for _, filepath := range filePaths {
		fileContentRaw, err := os.ReadFile(filepath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %q: %w", filepath, err)
		}

		fileContent := string(fileContentRaw)
		fileTypes, err := extractDocTypesFromFile(filepath, fileContent)
		if err != nil {
			return nil, fmt.Errorf("failed to extract types from file %q: %w", filepath, err)
		}

		for k, v := range fileTypes {
			types[k] = v
		}
	}

	return &TypeDocGenerator{
		fileTypes: types,
		tagName:   tagName,
	}, nil
}

func (gen *TypeDocGenerator) Generate(typesNames ...string) ([]*TypeDoc, error) {
	typeDocs := []*TypeDoc{}

	processedTypes := map[string]struct{}{}

	for len(typesNames) > 0 {
		// Dequeue
		name := typesNames[0]
		typesNames = typesNames[1:]

		t, ok := gen.fileTypes[name]
		if !ok {
			return nil, fmt.Errorf("type %s not found in the directory", name)
		}

		for _, s := range t.Decl.Specs {
			typeSpec, ok := s.(*ast.TypeSpec)
			if !ok {
				continue
			}

			typeStruct, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			c, err := parseComment(t.Doc)
			if err != nil {
				return nil, fmt.Errorf("failed for type %q: %w", name, err)
			}

			typeDoc := &TypeDoc{
				Name:        c.Name,
				Type:        t.Name,
				Description: c.Description,
				Note:        c.Note,
				Example:     c.Example,
				Fields:      []FieldDoc{},
			}

			for _, field := range typeStruct.Fields.List {
				fieldDoc, err := gen.processField(getFieldType(t.fileContent, field), field)
				if err != nil {
					return nil, fmt.Errorf("failed for type %q: %w", name, err)
				}

				if fieldDoc == nil {
					continue
				}

				typeDoc.Fields = append(typeDoc.Fields, *fieldDoc)

				if _, ok := gen.fileTypes[fieldDoc.lookType]; !ok {
					continue
				}

				if _, ok := processedTypes[fieldDoc.lookType]; !ok {
					// Enqueue next type
					typesNames = append(typesNames, fieldDoc.lookType)
					processedTypes[fieldDoc.lookType] = struct{}{}
				}

			}

			typeDocs = append(typeDocs, typeDoc)
		}
	}

	return typeDocs, nil
}

func (gen TypeDocGenerator) processField(fieldType string, field *ast.Field) (*FieldDoc, error) {
	if field.Tag == nil {
		return nil, nil
	}

	tag, err := parseTag(field.Tag.Value, gen.tagName)
	if err != nil {
		return nil, err
	}

	comment, err := parseComment(field.Doc.Text())
	if err != nil {
		return nil, fmt.Errorf("failed for field %q: %w", fieldType, err)
	}

	var options []string
	var isOptional bool

	for _, opt := range tag.Options {
		if opt == optionalTag && comment.OptionalIf == "" {
			isOptional = true
			continue
		}

		options = append(options, opt)
	}

	var lookupType string
	switch field.Type.(type) {
	case *ast.Ident:
		lookupType = fieldType
	case *ast.StarExpr:
		lookupType = strings.TrimLeft(fieldType, "*")
		fieldType = lookupType
		if comment.OptionalIf == "" {
			isOptional = true
		}
	case *ast.ArrayType:
		lookupType = strings.TrimLeft(fieldType, "[]")
		fieldType = fmt.Sprintf("list(%s)", lookupType)
	}

	return &FieldDoc{
		Name:        tag.Name,
		Type:        fieldType,
		Description: comment.Description,
		Note:        comment.Note,
		Examples:    comment.Examples,
		OptionalIf:  comment.OptionalIf,
		RequiredIf:  comment.RequiredlIf,
		Optional:    isOptional,
		Default:     comment.Default,
		Options:     options,
		Values:      comment.Values,

		lookType: lookupType,
	}, nil
}

func getFieldType(fileContent string, field *ast.Field) string {
	typeExpr := field.Type

	start := typeExpr.Pos() - 1
	end := typeExpr.End() - 1

	// grab it in source
	return fileContent[start:end]
}
