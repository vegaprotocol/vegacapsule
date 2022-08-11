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
	downloadingImageMessage = "Downloading image"
	eventStartedType        = "Started"
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

func (ai allocationInfo) String() string {
	events := make([]string, 0, len(ai.events))
	for _, e := range ai.events {
		events = append(events, fmt.Sprintf("%s: %s", e.Type, e.DisplayMessage))
	}

	return fmt.Sprintf("Task: %q, State: %q, Events: %q", ai.taskName, ai.taskState, strings.Join(events, " -> "))
}

type allocations []allocationInfo

func (allocs allocations) downloadingImage() bool {
	for _, a := range allocs {
		if a.downloadingImage() {
			return true
		}
	}

	return false
}

func (allocs allocations) started() bool {
	for _, a := range allocs {
		if !a.started() {
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

		if allocs.started() {
			return false, nil
		}

		if allocs.downloadingImage() {
			fmt.Println("has not timed out because downloading image")
			return false, nil
		}

		for _, alloc := range allocs {
			log.Printf("Job  %q has timed out output: %s", jobID, alloc)
		}

		return true, fmt.Errorf("failed to run %s job: starting deadline has been exceeded", jobID)
	default:
		return false, nil
	}
}
