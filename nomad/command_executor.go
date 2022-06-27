package nomad

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/hashicorp/nomad/api"
)

const (
	NomadStateRunning = "running"
)

type CommandExecutor struct {
	allocationsClient *api.Allocations
}

type remoteCommandRunnerDetails struct {
	NodeSet    types.NodeSet
	Allocation *api.Allocation
	TaskID     string
}

type commandCallback func(pathsMapping types.NetworkPathsMapping) []string

func NewCommandRunner(client *Client) *CommandExecutor {
	return &CommandExecutor{
		allocationsClient: client.API.Allocations(),
	}
}

func (runner *CommandExecutor) executeCallbacks(ctx context.Context, cmdCallbacks []commandCallback, nodeSets []types.NodeSet) (io.Reader, error) {
	allocations, err := runner.filterCommandRunnerAllocations(nodeSets)
	if err != nil {
		return nil, fmt.Errorf("failed to filter command runner alocations: %w", err)
	}

	var result bytes.Buffer
	buffWriter := bufio.NewWriter(&result)

	// TODO: Parallel below executions
	// The exec call blocks until command terminates (or an error occurs), and returns the exit code.
	for _, allocationDetails := range allocations {
		for _, cmdCallback := range cmdCallbacks {
			command := cmdCallback(allocationDetails.NodeSet.RemoteCommandRunner.PathsMapping)

			// No command given for some logic conditions
			if command == nil {
				continue
			}

			if _, err := buffWriter.WriteString(fmt.Sprintf(
				"\nRunning the %v command for the\"%s\" node set\n",
				command, allocationDetails.NodeSet.Name)); err != nil {

				return nil, fmt.Errorf("failed to write message to the out buffer: %w", err)
			}

			exitCode, err := runner.allocationsClient.Exec(ctx,
				allocationDetails.Allocation,
				allocationDetails.TaskID,
				false,
				command,
				strings.NewReader(""),
				buffWriter,
				buffWriter,
				nil,
				&api.QueryOptions{},
			)

			if err != nil {
				return nil, fmt.Errorf("execution of %v failed for the \"%s\" node-set with exitcode \"%d\": %w",
					command, allocationDetails.NodeSet.Name, exitCode, err)
			}
		}
	}
	buffWriter.Flush()
	reader := strings.NewReader(result.String())

	return reader, nil
}

func (runner *CommandExecutor) filterCommandRunnerAllocations(nodeSets []types.NodeSet) ([]*remoteCommandRunnerDetails, error) {
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

func (runner *CommandExecutor) getRemoteCommandRunnerDetails(allocationList []*api.AllocationListStub, nodeSet types.NodeSet) (*remoteCommandRunnerDetails, error) {
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
