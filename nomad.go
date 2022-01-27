package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/nomad/jobspec"
)

func parseJobFiles(dir string) {

	files, err := os.ReadDir(dir)
	path, _ := filepath.Abs(dir)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		p := filepath.Join(path, file.Name())

		_, err = jobspec.ParseFile(p)
		if err != nil {
			fmt.Println(err)
		}
	}
}
