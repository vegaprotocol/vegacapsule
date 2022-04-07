package nomad

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/nomad/api"
)

var logTypes = [2]string{"stdout", "stderr"}

type framesSet struct {
	name   string
	frames <-chan *api.StreamFrame
	errs   <-chan error
	cancel chan struct{}
}

func (n *Client) LogJobs(ctx context.Context, follow bool, origin string, offset int64, jobIDs []string) (io.ReadCloser, error) {
	jobsApi := n.API.Jobs()
	allocsApi := n.API.AllocFS()

	var frameSets []framesSet
	queryOpts := new(api.QueryOptions).WithContext(ctx)
	for _, jobID := range jobIDs {
		jobsAllocs, _, err := jobsApi.Allocations(jobID, true, queryOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to get nomad job %q: %w", jobID, err)
		}

		for _, as := range jobsAllocs {
			cAlloc := &api.Allocation{ID: as.ID, NodeID: as.NodeID}
			for taskName := range as.TaskStates {
				for _, logType := range logTypes {
					cancelCh := make(chan struct{}, 0)
					framesCh, errsCh := allocsApi.Logs(cAlloc, follow, taskName, logType, origin, offset, cancelCh, queryOpts)

					frameSets = append(frameSets,
						framesSet{
							name:   fmt.Sprintf("Job: %s, Task: %s", jobID, taskName),
							frames: framesCh,
							errs:   errsCh,
							cancel: cancelCh,
						},
					)
				}
			}
		}
	}

	return NewFrameReader(mergeFrameSets(frameSets)), nil
}

func mergeFrameSets(fss []framesSet) (<-chan *StreamFrame, <-chan error, chan struct{}) {
	frames := make(chan *StreamFrame)
	errs := make(chan error, 1)
	cancel := make(chan struct{}, 0)

	var wg sync.WaitGroup
	wg.Add(len(fss))
	for _, c := range fss {
		go func(fs framesSet) {
			for frame := range fs.frames {
				frames <- &StreamFrame{
					Name:        fs.name,
					StreamFrame: frame,
				}
			}
			wg.Done()
		}(c)
		go func(fs framesSet) {
			for err := range fs.errs {
				errs <- fmt.Errorf("received error from %q: %w", fs.name, err)
			}
		}(c)
	}

	go func(fss []framesSet) {
		select {
		case c := <-cancel:
			for _, fs := range fss {
				fs.cancel <- c
			}
		}
	}(fss)

	go func() {
		wg.Wait()
		close(frames)
		close(errs)
		close(cancel)
	}()

	return frames, errs, cancel
}
