package xdb

import (
	"fmt"
	"log"
	"runtime/debug"
)

// WorkerExecutionError represents runtime errors on worker execution
type WorkerExecutionError struct {
	OriginalError error
	StackTrace    string
}

func newWorkerExecutionError(err error, stackTrace string) error {
	return &WorkerExecutionError{
		OriginalError: err,
		StackTrace:    stackTrace,
	}
}

func (i WorkerExecutionError) Error() string {
	return fmt.Sprintf("error message:%v, stacktrace: %v", i.OriginalError, i.StackTrace)
}

// for skipping the logging in testing code
var skipCaptureErrorLogging = false

// MUST be the result from calling recover, which MUST be done in a single level deep
// deferred function. The usual way of calling this is:
// - defer func() { captureStateExecutionError(recover(), logger, &err) }()
func captureStateExecutionError(errPanic interface{}, retError *error) {
	if errPanic != nil || *retError != nil {
		st := string(debug.Stack())

		var err error
		panicError, ok := errPanic.(error)
		if errPanic != nil {
			if ok && panicError != nil {
				err = newWorkerExecutionError(panicError, st)
			} else {
				err = newWorkerExecutionError(fmt.Errorf("errPanic is not an error %v", errPanic), st)
			}
		} else {
			err = newWorkerExecutionError(*retError, st)
		}

		if !skipCaptureErrorLogging && errPanic != nil {
			log.Printf("panic is captured: %v , stacktrace: %v", errPanic, st)
		}
		*retError = err
	}
}
