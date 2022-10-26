package main

import (
	"flag"
	"fmt"

	"code.vegaprotocol.io/vegacapsule/docsgenerator"
)

var (
	tagName       string
	typeName      string
	directoryPath string
)

var description = `Capsule configuration is used by vegacapsule CLI network generate and bootstrap commands.
It allows to configue and customise Vega network running on Capsule.

The configuration is using [HCL](https://github.com/hashicorp/hcl) language syntax also used for example by [Terraform](https://www.terraform.io/).

This document explains all possible configuration options in Capsule.
`

func init() {
	flag.StringVar(&tagName, "tag-name", "", "name of the tag")
	flag.StringVar(&typeName, "type-name", "", "type to be processed")
	flag.StringVar(&directoryPath, "dir-path", "", "directory path of the file to generate docs from")
}

func main() {
	flag.Parse()

	if tagName == "" {
		panic("missing required `tag-name` flag")
	}
	if typeName == "" {
		panic("missing required `type-name` flag")
	}
	if directoryPath == "" {
		panic("missing required `dir-path` flag")
	}

	gen, err := docsgenerator.NewTypeDocGenerator(directoryPath, tagName)
	if err != nil {
		panic(err)
	}

	typeDocs, err := gen.Generate(typeName)
	if err != nil {
		panic(err)
	}

	fd := docsgenerator.NewFileDoc(
		"Capsule configuration docs",
		description,
		typeDocs,
	)

	b, err := fd.Encode()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
