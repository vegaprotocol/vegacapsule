package commands

import (
	"bytes"
	"fmt"
	"io"

	"code.vegaprotocol.io/vegacapsule/types"
)

func ResetNodeSetsData(binary string, nss []types.NodeSet) (io.Reader, error) {
	buff := bytes.NewBuffer([]byte{})

	for _, ns := range nss {
		vegaOut, err := VegaUnsafeResetAll(binary, ns.Vega.HomeDir)
		if err != nil {
			return nil, err
		}

		tendOut, err := TendermintUnsafeResetAll(binary, ns.Tendermint.HomeDir)
		if err != nil {
			return nil, err
		}

		buff.WriteString(fmt.Sprintf("Name: %s", ns.Vega.Name))
		buff.Write(vegaOut)
		buff.WriteString(fmt.Sprintf("Name: %s", ns.Tendermint.Name))
		buff.Write(tendOut)
	}

	return buff, nil
}
