package docsgenerator

import (
	"fmt"
	"go/ast"
	"go/doc"
	"os"
	"strings"
)

const optionalTag = "optional"

type TypeDocGenerator struct {
	fileContent string
	fileTypes   map[string]*doc.Type

	tagName string
}

func NewTypeDocGenerator(filepath, tagName string) (*TypeDocGenerator, error) {
	fileContentRaw, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", filepath, err)
	}

	fileContent := string(fileContentRaw)
	types, err := extractTypesFromFile(filepath, fileContent)
	if err != nil {
		return nil, fmt.Errorf("failed to extract types from file %q: %w", filepath, err)
	}

	return &TypeDocGenerator{
		fileContent: fileContent,
		fileTypes:   types,
		tagName:     tagName,
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
			return nil, fmt.Errorf("type %s not found in the file", name)
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
				return nil, err
			}

			typeDoc := &TypeDoc{
				Name:        t.Name,
				Type:        t.Name,
				Description: c.Description,
				Note:        c.Note,
				Example:     c.Example,
				Fields:      []FieldDoc{},
			}

			for _, field := range typeStruct.Fields.List {
				fieldDoc, err := gen.processField(field)
				if err != nil {
					return nil, err
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

func (gen TypeDocGenerator) getFieldType(field *ast.Field) string {
	typeExpr := field.Type

	start := typeExpr.Pos() - 1
	end := typeExpr.End() - 1

	// grab it in source
	return gen.fileContent[start:end]
}

func (gen TypeDocGenerator) processField(field *ast.Field) (*FieldDoc, error) {
	// grab it in source
	fieldType := gen.getFieldType(field)

	if field.Tag == nil {
		return nil, nil
	}

	tag, err := parseTag(field.Tag.Value, gen.tagName)
	if err != nil {
		return nil, err
	}

	comment, err := parseComment(field.Doc.Text())
	if err != nil {
		return nil, err
	}

	var options []string
	var isOptional bool

	for _, opt := range tag.Options {
		if opt == optionalTag {
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
		isOptional = true
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
		Optional:    isOptional,
		Options:     options,

		lookType: lookupType,
	}, nil
}
