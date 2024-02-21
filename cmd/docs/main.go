package main

import (
	"flag"
	"fmt"
	"strings"

	"code.vegaprotocol.io/vegacapsule/docsgenerator"

	"github.com/cometbft/cometbft/libs/os"
)

var (
	tagName         string
	typeNames       string
	directoryPath   string
	descriptionPath string
)

func init() {
	flag.StringVar(&tagName, "tag-name", "", "name of the tag")
	flag.StringVar(&typeNames, "type-names", "", "comma separated types to be processed")
	flag.StringVar(&directoryPath, "dir-path", "", "directory path of the file to generate docs from")
	flag.StringVar(&descriptionPath, "description-path", "", "path of file with description")
}

func main() {
	flag.Parse()

	if typeNames == "" {
		panic("missing required `type-names` flag")
	}
	if directoryPath == "" {
		panic("missing required `dir-path` flag")
	}
	if descriptionPath == "" {
		panic("missing required `description-path` flag")
	}

	description, err := os.ReadFile(descriptionPath)
	if err != nil {
		panic(err)
	}

	gen, err := docsgenerator.NewTypeDocGenerator(directoryPath, tagName)
	if err != nil {
		panic(err)
	}

	names := strings.Split(typeNames, ",")

	typeDocs, err := gen.Generate(names...)
	if err != nil {
		panic(err)
	}

	fd := docsgenerator.NewFileDoc(
		string(description),
		typeDocs,
	)

	b, err := fd.Encode()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
