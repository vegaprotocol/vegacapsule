package main

import (
	"flag"
	"fmt"

	"code.vegaprotocol.io/vegacapsule/docsgenerator"
)

var (
	tagName  string
	typeName string
	filePath string
)

func init() {
	flag.StringVar(&tagName, "tag-name", "", "name of the tag")
	flag.StringVar(&typeName, "type-name", "", "type to be processed")
	flag.StringVar(&filePath, "filepath", "", "path of the file to generate docs from")
}

func main() {
	flag.Parse()

	if tagName == "" {
		panic("missing required `tag-name` flag")
	}
	if typeName == "" {
		panic("missing required `type-name` flag")
	}
	if filePath == "" {
		panic("missing required `filePath` flag")
	}

	gen, err := docsgenerator.NewTypeDocGenerator(filePath, tagName)
	if err != nil {
		panic(err)
	}

	typeDocs, err := gen.Generate(typeName)
	if err != nil {
		panic(err)
	}

	fd := docsgenerator.NewFileDoc(
		"Capsule configuration",
		"This Capsule configuration file allows to configurate custom Vega network.",
		typeDocs,
	)

	b, err := fd.Encode()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
