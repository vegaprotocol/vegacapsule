package installer

import (
	"fmt"
	"log"
	"runtime"

	"code.vegaprotocol.io/vegacapsule/utils"
)

func formatAssetName(name string) string {
	return fmt.Sprintf("%s-%s-%s.zip", name, runtime.GOOS, runtime.GOARCH)
}

func cpAndChmodxFile(source, destination string) error {
	if err := utils.CpAndChmodxFile(source, destination); err != nil {
		return err
	}
	log.Printf("Successfully copied from %q to %q", source, destination)
	return nil
}
