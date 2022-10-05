package nomad

import (
	"errors"
	"fmt"
	"strings"
)

type ConnectionError struct {
	Err error
}

func (ce *ConnectionError) Error() string {
	return fmt.Sprintf("failed to connect to nomad: %s", ce.Err.Error())
}

func newConnectionErr(err error) *ConnectionError {
	return &ConnectionError{
		Err: err,
	}
}

func IsConnectionErr(err error) bool {
	var cerr *ConnectionError
	return errors.As(err, &cerr)
}

type JobTimeoutError struct {
	JobID string
}

func (ce *JobTimeoutError) Error() string {
	return fmt.Sprintf("failed to run %s job: starting deadline has been exceeded", ce.JobID)
}

func newJobTimeoutErr(jobID string) *JobTimeoutError {
	return &JobTimeoutError{
		JobID: jobID,
	}
}

func IsJobTimeoutErr(err error) bool {
	var cerr *JobTimeoutError
	return errors.As(err, &cerr)
}

func isCancelledError(err error) bool {
	return strings.Contains(err.Error(), "Cancelled")
}
