package docsgenerator

import (
	"fmt"
	"go/ast"
	"os"
	"strings"
)

const optionalTag = "optional"

type TypeDocGenerator struct {
	fileTypes map[string]docTypeWithFileContent

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
				var fieldDoc *FieldDoc

				// TODO - field.Tag == nil is not good enough. Need to figure out what to do without tag for private fields on struct.
				if c.IgnoreTag || gen.tagName == "" {
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
		Type:        fi.fieldType,
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
		name = fi.fieldType
	}

	return &FieldDoc{
		// TODO - implement function to extract the name properly
		Name:        name,
		Type:        fi.fieldType,
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

// formatLookupKey formats strings to "packageName.TypeName"
func formatLookupKey(packageName, typeName string) string {
	return fmt.Sprintf("%s.%s", packageName, typeName)
}

type fieldInfo struct {
	fieldType  string
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
		fieldType: fieldType,
	}

	packageName := currentPackageName

	typeSplit := strings.Split(fi.fieldType, ".")
	if len(typeSplit) > 1 {
		packageName = typeSplit[0]
		fi.lookupKey = typeSplit[1]
	} else {
		fi.lookupKey = typeSplit[0]
	}

	switch field.Type.(type) {
	case *ast.Ident:
		fi.lookupKey = fi.fieldType
	case *ast.StarExpr:
		fi.lookupKey = strings.TrimLeft(fi.fieldType, "*")
		fi.fieldType = fi.lookupKey
		if comment.OptionalIf == "" {
			fi.isOptional = true
		}
	case *ast.ArrayType:
		fi.lookupKey = strings.TrimLeft(fi.fieldType, "[]")
		fi.fieldType = fmt.Sprintf("list(%s)", fi.lookupKey)
	}

	fi.lookupKey = formatLookupKey(packageName, fi.lookupKey)

	return fi
}

func getFieldType(fileContent string, field *ast.Field) string {
	typeExpr := field.Type

	start := typeExpr.Pos() - 1
	end := typeExpr.End() - 1

	// grab it in source
	return fileContent[start:end]
}
