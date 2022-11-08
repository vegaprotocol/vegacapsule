package docsgenerator

/*
Example of usage:

name: Type Document
type: TypeDoc
description: TypeDoc represents the types of documentation in a file
node: A note to be displayed at the end of the doc
example:

	type: json
	name: JSON usage
	value: { "Name": "name", "Type": "type", "Description": "description" }
*/
type TypeDoc struct {
	// Name represents struct name or field name.
	Name string
	// Type represents struct name or field type.
	Type string
	// Description represents the full description for the item.
	Description string
	// Note is rendered as a note for the example in markdown file.
	Note string
	// Fields contains fields documentation if related item is a struct.
	Fields []FieldDoc
	// Example values for the type.
	Example Example
}

// Doc represents a struct documentation rendered from comments by docgen.
type FieldDoc struct {
	// Name represents struct name or field name.
	Name string
	// Type represents struct name or field type.
	Type string
	// Description represents the full description for the item.
	Description string
	// Note is rendered as a note for the example in markdown file.
	Note string

	Optional   bool
	OptionalIf string
	RequiredIf string
	Default    string

	// Examples: List of example values for the item.
	Examples []Example
	// Values is only used to render valid values list in the documentation.
	Values []string

	// Options renders extra options for this field.
	Options []string

	// lookType is for internal usage only.
	// It is the field type stripped of pointer and array symbols (*, [])
	lookType string
}

type Example struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
	Type  string `yaml:"type"`
}

type Comment struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Note        string    `yaml:"note"`
	Default     string    `yaml:"default"`
	Examples    []Example `yaml:"examples"`
	Example     Example   `yaml:"example"`
	Values      []string  `yaml:"values"`
	// Makes item param if specified param is defined
	OptionalIf  string `yaml:"optional_if"`
	RequiredlIf string `yaml:"required_if"`
}
