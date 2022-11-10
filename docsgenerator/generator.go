package docsgenerator

import (
	"fmt"
	"go/ast"
	"os"
	"regexp"
	"unicode"
)

const optionalTag = "optional"

var (
	mapRegex = regexp.MustCompile(`map\[(.*)\](.*)`)
)

type TypeDocGenerator struct {
	fileTypes map[string]docTypeWithFileContent

	documentStructsMode bool

	tagName string
}

func NewTypeDocGenerator(dir, tagName string) (*TypeDocGenerator, error) {
	filePaths := findFiles(dir, ".go")

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
			types[formatLookupKey(v.packageName, k)] = v
		}
	}

	var documentStructsMode bool
	if tagName == "" {
		documentStructsMode = true
	}

	return &TypeDocGenerator{
		fileTypes:           types,
		tagName:             tagName,
		documentStructsMode: documentStructsMode,
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
				Type:        gen.formatFieldType(t.packageName, t.Name),
				Description: c.Description,
				Note:        c.Note,
				Example:     c.Example,
				Fields:      []FieldDoc{},
			}

			for _, field := range typeStruct.Fields.List {
				var fieldDoc *FieldDoc

				if gen.documentStructsMode {
					fieldDoc, err = gen.processField(t.packageName, getFieldType(t.fileContent, field), field)
				} else {
					fieldDoc, err = gen.processFieldWithTag(t.packageName, getFieldType(t.fileContent, field), field)
				}

				if err != nil {
					return nil, fmt.Errorf("failed for type %q: %w", name, err)
				}

				if fieldDoc == nil {
					continue
				}

				typeDoc.Fields = append(typeDoc.Fields, *fieldDoc)

				if _, ok := gen.fileTypes[fieldDoc.lookupKey]; !ok {
					continue
				}

				if _, ok := processedTypes[fieldDoc.lookupKey]; !ok {
					// Enqueue next type
					typesNames = append(typesNames, fieldDoc.lookupKey)
					processedTypes[fieldDoc.lookupKey] = struct{}{}
				}

			}

			typeDocs = append(typeDocs, typeDoc)
		}
	}

	return typeDocs, nil
}

func (gen TypeDocGenerator) processFieldWithTag(packageName, fieldType string, field *ast.Field) (*FieldDoc, error) {
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

	fi := getFieldInfo(packageName, fieldType, field, comment, isOptional)

	return &FieldDoc{
		Name:        tag.Name,
		Type:        gen.formatFieldType(packageName, fieldType),
		Description: comment.Description,
		Note:        comment.Note,
		Examples:    comment.Examples,
		OptionalIf:  comment.OptionalIf,
		RequiredIf:  comment.RequiredlIf,
		Optional:    fi.isOptional,
		Default:     comment.Default,
		Options:     options,
		Values:      comment.Values,

		lookupKey: fi.lookupKey,
	}, nil
}

func (gen TypeDocGenerator) processField(packageName, fieldType string, field *ast.Field) (*FieldDoc, error) {
	comment, err := parseComment(field.Doc.Text())
	if err != nil {
		return nil, fmt.Errorf("failed for field %q: %w", fieldType, err)
	}

	fi := getFieldInfo(packageName, fieldType, field, comment, false)

	var name string
	if len(field.Names) > 0 {
		name = field.Names[0].Name
	} else {
		name = fieldType
	}

	return &FieldDoc{
		Name:        name,
		Type:        gen.formatFieldType(packageName, fieldType),
		Description: comment.Description,
		Note:        comment.Note,
		Examples:    comment.Examples,
		OptionalIf:  comment.OptionalIf,
		RequiredIf:  comment.RequiredlIf,
		Optional:    fi.isOptional,
		Default:     comment.Default,
		Values:      comment.Values,

		lookupKey: fi.lookupKey,
	}, nil
}

// TODO - consider implemeting by using actual ast types rather then matching strings
func (gen TypeDocGenerator) formatFieldType(packageName, fieldType string) string {
	// pointer
	if fieldType[0] == '*' {
		return gen.formatFieldType(packageName, fieldType[1:])
	}

	// slice
	if fieldType[0:2] == "[]" {
		return fmt.Sprintf("[]%s", gen.formatFieldType(packageName, fieldType[2:]))
	}

	// map
	if len(fieldType) > 4 && fieldType[0:4] == "map[" {
		key, val := typesFromMap(fieldType)
		return fmt.Sprintf("map[%s]%s", key, gen.formatFieldType(packageName, val))
	}

	// struct - only structs can start with uppercase letter
	if unicode.IsUpper(rune(fieldType[0])) && gen.documentStructsMode {
		return formatLookupKey(packageName, fieldType)
	}

	// scalar or simple struct type without package prepend
	return fieldType
}
