package probes

import (
	"errors"
	"fmt"
	"strings"
)

type RetryableError struct {
	e error
}

func (re *RetryableError) Error() string {
	return fmt.Sprintf("failed to probe: %s", re.e)
}

func newRetryableError(err error) *RetryableError {
	return &RetryableError{e: err}
}

func IsRetryableErr(err error) bool {
	var cerr *RetryableError
	return errors.As(err, &cerr)
}

func isConnectionErr(err error) bool {
	return strings.Contains(err.Error(), "connection refused")
}
