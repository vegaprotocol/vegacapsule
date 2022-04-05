package nomad

import (
	"context"
	"io"
	"sync"

	"github.com/hashicorp/nomad/api"
)

var logTypes = [2]string{"stdout", "stderr"}

type framesSet struct {
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
			return nil, err
		}

		for _, as := range jobsAllocs {
			cAlloc := &api.Allocation{ID: as.ID, NodeID: as.NodeID}
			for taskName := range as.TaskStates {
				for _, logType := range logTypes {
					cancelCh := make(chan struct{}, 0)
					framesCh, errsCh := allocsApi.Logs(cAlloc, follow, taskName, logType, origin, offset, cancelCh, queryOpts)

					frameSets = append(frameSets,
						framesSet{
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

func mergeFrameSets(fss []framesSet) (<-chan *api.StreamFrame, <-chan error, chan struct{}) {
	frames := make(chan *api.StreamFrame)
	errs := make(chan error, 1)
	cancel := make(chan struct{}, 0)

	var wg sync.WaitGroup
	wg.Add(len(fss))
	for _, c := range fss {
		go func(fs framesSet) {
			for frame := range fs.frames {
				frames <- frame
			}
			wg.Done()
		}(c)
		go func(fs framesSet) {
			for err := range fs.errs {
				errs <- err
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
