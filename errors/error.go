package errors

import "errors"

// ErrorStatus extends error adding a status method
type ErrorStatus interface {
	error
	Status() int
}

type errStatus struct {
	error
	status int
}

// NewErrorStatus creates an error that implements ErrorStatus based off a
// string
func NewErrorStatus(status int, err string) error {
	return &errStatus{
		error:  errors.New(err),
		status: status,
	}
}

// ConvertErrorStatus creates an error that implements ErrorStatus based off an
// error
func ConvertErrorStatus(status int, err error) error {
	return &errStatus{
		error:  err,
		status: status,
	}
}

func (e errStatus) Status() int {
	return e.status
}
