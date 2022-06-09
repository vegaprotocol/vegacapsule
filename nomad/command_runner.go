package nomad

import "code.vegaprotocol.io/vegacapsule/types"

type CommandRunner struct {
	client *Client
}

func NewCommandRunner(client *Client) *CommandRunner {
	return &CommandRunner{
		client: client,
	}
}

func (runner *CommandRunner) GetNodes(nodeSets []*types.NodeSet) {
	// jobsCli := runner.client.API.Jobs()
	// jobsCli.Info()
}
