package nomad

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/nomad/api"
)

const (
	downloadingImageMessage      = "Downloading image"
	downloadingArtifactsMessage  = "downloading Artifacts"
	taskReceivedMessage          = "received by client"
	buildingTaskDirectoryMessage = "building task directory"
	eventStartedType             = "Started"
)

type allocationInfo struct {
	taskName  string
	taskState string
	events    []*api.TaskEvent
}

func (ai allocationInfo) started() bool {
	if ai.taskState != Running {
		return false
	}

	if len(ai.events) == 0 {
		return false
	}

	e := ai.events[len(ai.events)-1]

	return e.Type == eventStartedType
}

// finishedWithoutError is required for pre-start tasks that may be
// stopped before the main task
func (ai allocationInfo) finishedWithoutError() bool {
	if ai.taskState != Dead {
		return false
	}

	if len(ai.events) == 0 {
		return false
	}

	e := ai.events[len(ai.events)-1]

	return strings.ToLower(e.Type) == Terminated &&
		e.ExitCode == 0 &&
		!e.FailsTask

}

func (ai allocationInfo) downloadingImage() bool {
	if ai.taskState != Running {
		return false
	}

	if len(ai.events) == 0 {
		return false
	}

	e := ai.events[len(ai.events)-1]

	return strings.Contains(e.DisplayMessage, downloadingImageMessage)
}

func (ai allocationInfo) downloadingArtifacts() bool {
	if ai.taskState != Pending {
		return false
	}

	if len(ai.events) == 0 {
		return false
	}

	e := ai.events[len(ai.events)-1]

	return strings.Contains(strings.ToLower(e.DisplayMessage), downloadingArtifactsMessage)
}

func (ai allocationInfo) buildingTaskDirectory() bool {
	if ai.taskState != Pending {
		return false
	}

	if len(ai.events) == 0 {
		return false
	}

	e := ai.events[len(ai.events)-1]

	return strings.Contains(strings.ToLower(e.DisplayMessage), buildingTaskDirectoryMessage)
}

func (ai allocationInfo) taskReceived() bool {
	if ai.taskState != Pending {
		return false
	}

	if len(ai.events) == 0 {
		return false
	}

	e := ai.events[len(ai.events)-1]

	return strings.Contains(strings.ToLower(e.DisplayMessage), taskReceivedMessage)
}

func (ai allocationInfo) String() string {
	events := make([]string, 0, len(ai.events))
	for _, e := range ai.events {
		events = append(events, fmt.Sprintf("%s: %s", e.Type, e.DisplayMessage))
	}

	return fmt.Sprintf("Task: %q, State: %q, Events: %q", ai.taskName, ai.taskState, strings.Join(events, " -> "))
}

type allocations []allocationInfo

func (allocs allocations) pending() bool {
	for _, a := range allocs {
		if a.downloadingImage() {
			return true
		}

		if a.downloadingArtifacts() {
			return true
		}

		if a.taskReceived() {
			return true
		}

		if a.buildingTaskDirectory() {
			return true
		}
	}

	return false
}

func (allocs allocations) startedOrFinishedWithoutError() bool {
	for _, a := range allocs {
		if !a.started() && !a.finishedWithoutError() {
			return false
		}
	}

	return true
}

func (n *Client) getJobAllocsInfo(ctx context.Context, jobID string) (allocations, error) {
	allocs, _, err := n.API.Jobs().Allocations(jobID, false, &api.QueryOptions{})
	if err != nil {
		return nil, err
	}

	allocsInfo := make([]allocationInfo, 0, len(allocs))

	for _, alloc := range allocs {
		for tsName, ts := range alloc.TaskStates {
			allocsInfo = append(allocsInfo, allocationInfo{
				taskName:  tsName,
				taskState: ts.State,
				events:    ts.Events,
			})
		}
	}

	return allocsInfo, nil
}

func (n *Client) jobTimedOut(ctx context.Context, t *time.Ticker, jobID string) (bool, error) {
	select {
	case <-t.C:
		allocs, err := n.getJobAllocsInfo(ctx, jobID)
		if err != nil {
			return false, err
		}

		if allocs.startedOrFinishedWithoutError() {
			return false, nil
		}

		if allocs.pending() {
			fmt.Println("has not timed out because pending task")
			return false, nil
		}

		for _, alloc := range allocs {
			log.Printf("Job %q has timed out output: %s", jobID, alloc)
		}

		return true, fmt.Errorf("failed to run %s job: starting deadline has been exceeded", jobID)
	default:
		return false, nil
	}
}
