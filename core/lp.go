package core

import "strconv"

// LookPathError is returned by LookPath when it fails to classify a file as an
// executable.
type LookPathError struct {
	// Name is the file name for which the error occurred.
	Name string
	// Err is the underlying error.
	Err error
}

func (e *LookPathError) Error() string {
	return "LookPath: " + strconv.Quote(e.Name) + ": " + e.Err.Error()
}

func (e *LookPathError) Unwrap() error { return e.Err }
