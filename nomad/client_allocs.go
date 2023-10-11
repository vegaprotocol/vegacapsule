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
	downloadingImageMessage      = "downloading image"
	imagePullProgressMessage     = "image pull progress"
	downloadingArtifactsMessage  = "downloading artifacts"
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

// lastEventMessageContains lowercase the message and check if message contains given text.
func (ai allocationInfo) lastEventMessageContains(text string) bool {
	if len(ai.events) == 0 {
		return false
	}

	e := ai.events[len(ai.events)-1]

	return strings.Contains(strings.ToLower(e.DisplayMessage), text)
}

func (ai allocationInfo) downloadingImage() bool {
	return ai.lastEventMessageContains(downloadingImageMessage) || ai.lastEventMessageContains(imagePullProgressMessage)
}

func (ai allocationInfo) downloadingArtifacts() bool {
	return ai.lastEventMessageContains(downloadingArtifactsMessage)
}

func (ai allocationInfo) buildingTaskDirectory() bool {
	return ai.lastEventMessageContains(buildingTaskDirectoryMessage)
}

func (ai allocationInfo) taskReceived() bool {
	return ai.lastEventMessageContains(taskReceivedMessage)
}

func (ai allocationInfo) String() string {
	events := make([]string, 0, len(ai.events))
	for _, e := range ai.events {
		events = append(events, fmt.Sprintf("%s: %s", e.Type, e.DisplayMessage))
	}

	return fmt.Sprintf("Task: %q, State: %q, Events: %q", ai.taskName, ai.taskState, strings.Join(events, " -> "))
}

func formatPendingReason(task, reason string) string {
	return fmt.Sprintf("task %q is %s", task, reason)
}

type allocations []allocationInfo

func (allocs allocations) pending() (bool, string) {
	for _, a := range allocs {
		if a.downloadingImage() {
			return true, formatPendingReason(a.taskName, downloadingImageMessage)
		}

		if a.downloadingArtifacts() {
			return true, formatPendingReason(a.taskName, downloadingArtifactsMessage)
		}

		if a.taskReceived() {
			return true, formatPendingReason(a.taskName, taskReceivedMessage)
		}

		if a.buildingTaskDirectory() {
			return true, formatPendingReason(a.taskName, buildingTaskDirectoryMessage)
		}
	}

	return false, ""
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
			log.Printf("All tasks for job %q started or finished without an error", jobID)
			return false, nil
		}

		if pending, reason := allocs.pending(); pending {
			log.Printf("Job %q has not timed out because it has pending task: %s", jobID, reason)
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
