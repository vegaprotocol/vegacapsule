package main

import (
	"os"

	"code.vegaprotocol.io/vegacapsule/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
