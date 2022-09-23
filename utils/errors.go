package utils

import "strings"

type MultiError struct {
	errors []error
}

func NewMultiError() *MultiError {
	return &MultiError{}
}

func (e *MultiError) Add(err error) {
	e.errors = append(e.errors, err)
}

func (e *MultiError) HasAny() bool {
	return len(e.errors) > 0
}

func (e *MultiError) Error() string {
	fmtErrors := make([]string, 0, len(e.errors))
	for _, err := range e.errors {
		fmtErrors = append(fmtErrors, err.Error())
	}

	return strings.Join(fmtErrors, ", also ")
}
