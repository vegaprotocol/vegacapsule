package nomad

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"text/template"

	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/hashicorp/nomad/api"
)

const (
	NomadStateRunning = "running"
)

type CommandRunner struct {
	allocationsClient *api.Allocations
}

type command struct {
	binary string
	args   []string
}

func (cmd *command) ToSlice() []string {
	return append([]string{cmd.binary}, cmd.args...)
}

type remoteCommandRunnerDetails struct {
	NodeSet    types.NodeSet
	Allocation *api.Allocation
	TaskID     string
	Command    command
}

func NewCommandRunner(client *Client) *CommandRunner {
	return &CommandRunner{
		allocationsClient: client.API.Allocations(),
	}
}

func (runner *CommandRunner) Execute(ctx context.Context, binary string, args []string, nodeSets []types.NodeSet) (io.Reader, error) {
	allocations, err := runner.filterCommandRunnerAllocations(nodeSets)
	if err != nil {
		return nil, fmt.Errorf("failed to filter command runner alocations: %w", err)
	}

	for idx, allocationDetails := range allocations {
		command, err := remoteCommandForNodeSet(allocationDetails.NodeSet, binary, args)
		if err != nil {
			return nil, fmt.Errorf("failed to prepare command for nodeset \"%s\": %w", allocationDetails.NodeSet.Name, err)
		}

		allocations[idx].Command = *command
	}

	var result bytes.Buffer
	buffWriter := bufio.NewWriter(&result)
	// buffReader := bufio.NewReader(&result)
	// TODO: Parallel below executions
	// The exec call blocks until command terminates (or an error occurs), and returns the exit code.
	for _, allocationDetails := range allocations {
		buffWriter.WriteString(fmt.Sprintf(
			"\nRunning the %v command for the\"%s\" node set\n",
			allocationDetails.Command.ToSlice(), allocationDetails.NodeSet.Name))
		exitCode, err := runner.allocationsClient.Exec(ctx,
			allocationDetails.Allocation,
			allocationDetails.TaskID,
			false,
			allocationDetails.Command.ToSlice(),
			strings.NewReader(""),
			buffWriter,
			buffWriter,
			nil,
			&api.QueryOptions{},
		)

		if err != nil {
			return nil, fmt.Errorf("execution of %v failed for the \"%s\" node-set with exitcode \"%d\": %w",
				allocationDetails.Command.ToSlice(), allocationDetails.NodeSet.Name, exitCode, err)
		}
	}
	buffWriter.Flush()
	reader := strings.NewReader(result.String())
	return reader, nil
}

func remoteCommandForNodeSet(nodeSet types.NodeSet, binary string, args []string) (*command, error) {
	newBinary, err := applyPathsMappingToString(nodeSet.RemoteCommandRunner.PathsMapping, binary)
	if err != nil {
		return nil, fmt.Errorf("failed to apply paths mapping for binary(\"%s\"): %w", binary, err)
	}

	newArgs := []string{}
	for _, oldArg := range args {
		arg, err := applyPathsMappingToString(nodeSet.RemoteCommandRunner.PathsMapping, oldArg)
		if err != nil {
			return nil, fmt.Errorf("failed to apply paths mapping for one of the arguments(\"%s\"): %w", oldArg, err)
		}
		newArgs = append(newArgs, arg)
	}

	return &command{
		binary: newBinary,
		args:   newArgs,
	}, nil
}

func applyPathsMappingToString(mapping types.NetworkPathsMapping, arg string) (string, error) {
	tmpl, err := template.New("cmdArg").Parse(arg)
	if err != nil {
		return "", fmt.Errorf("failed to parse argument: %w", err)
	}

	buff := bytes.NewBufferString("")
	if err = tmpl.Execute(buff, mapping); err != nil {
		return "", fmt.Errorf("failed to execute templating for string \"%s\": %w", arg, err)
	}

	return buff.String(), nil
}

func (runner *CommandRunner) filterCommandRunnerAllocations(nodeSets []types.NodeSet) ([]*remoteCommandRunnerDetails, error) {
	if err := validateNodeSets(nodeSets); err != nil {
		return nil, fmt.Errorf("failed to run command on remote nomad cluster: %w", err)
	}

	networkAllocations, _, err := runner.allocationsClient.List(&api.QueryOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nomad allocations: %w", err)
	}

	result := []*remoteCommandRunnerDetails{}
	for _, nodeSet := range nodeSets {
		alloc, err := runner.getRemoteCommandRunnerDetails(networkAllocations, nodeSet)
		if err != nil {
			return nil, fmt.Errorf("failed to get command runner allocation for \"%s\" node set: %w", nodeSet.Name, err)
		}

		result = append(result, alloc)
	}

	return result, nil
}

func (runner *CommandRunner) getRemoteCommandRunnerDetails(allocationList []*api.AllocationListStub, nodeSet types.NodeSet) (*remoteCommandRunnerDetails, error) {
	var foundAllocationStub *api.AllocationListStub

	for _, allocStub := range allocationList {
		if !strings.HasPrefix(allocStub.JobID, nodeSet.RemoteCommandRunner.Name) {
			continue
		}

		foundAllocationStub = allocStub
		break
	}

	if foundAllocationStub == nil {
		return nil, fmt.Errorf("the remote command runner allocation not found")
	}

	alloc, _, err := runner.allocationsClient.Info(foundAllocationStub.ID, &api.QueryOptions{})

	if err != nil {
		return nil, fmt.Errorf("failed to get info for nomad allocation: %w", err)
	}

	var firstTaskID string
	for taskID := range alloc.TaskStates {
		firstTaskID = taskID
		break
	}

	if firstTaskID == "" {
		return nil, fmt.Errorf("the remote command runner allocation has no task defined")
	}

	taskState, taskStateOK := alloc.TaskStates[firstTaskID]
	if !taskStateOK {
		return nil, fmt.Errorf("failed to get task state details for \"%s\" task", firstTaskID)
	}

	if taskState.State != NomadStateRunning || taskState.Failed {
		return nil, fmt.Errorf("task is not running or failed")
	}

	return &remoteCommandRunnerDetails{
		NodeSet:    nodeSet,
		Allocation: alloc,
		TaskID:     firstTaskID,
	}, nil
}

func validateNodeSets(nodeSets []types.NodeSet) error {
	for _, nodeSet := range nodeSets {
		if nodeSet.RemoteCommandRunner == nil {
			return fmt.Errorf("the remote command runner is not specified for one or more node sets")
		}
	}

	return nil
}
