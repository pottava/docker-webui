package engine

import (
	"time"
)

// DockerTimeoutError occers when the request was timeout
type DockerTimeoutError struct {
	duration   time.Duration
	transition string
}

func (err *DockerTimeoutError) Error() string {
	return "Could not transition to " + err.transition + "; timed out after waiting " + err.duration.String()
}

// ErrorName returns the error name
func (err *DockerTimeoutError) ErrorName() string {
	return "DockerTimeoutError"
}

// CannotXContainerError occers wher the request went wrong
type CannotXContainerError struct {
	transition string
	msg        string
}

func (err CannotXContainerError) Error() string {
	return err.msg
}

// ErrorName returns the error name
func (err CannotXContainerError) ErrorName() string {
	return "Cannot" + err.transition + "ContainerError"
}
