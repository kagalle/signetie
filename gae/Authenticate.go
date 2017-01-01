package gae

import "github.com/go-errors/errors"

// Authenticate is a class to handle the first part of the gae authentication process.
type Authenticate struct {
	found     bool
	code      string
	cancelled bool
	err       *errors.Error
}

// Code is a getter for the resulting code from the gae authentication process.
func (auth *Authenticate) Code() string {
	return auth.code
}

// Found is a getter; was the process successful.
func (auth *Authenticate) Found() bool {
	return auth.found
}

func (auth *Authenticate) setFound() {
	auth.found = true
	auth.cancelled = false // insure consistency
}

// Cancelled is a getter; was the process cancelled by the user, either
// by clicking the cancel button or by closing the authentate window prematurely.
func (auth *Authenticate) Cancelled() bool {
	return auth.cancelled
}
func (auth *Authenticate) setCancelled() {
	auth.found = false
	auth.cancelled = true // insure consistency
}

func (auth *Authenticate) Error() error {
	return auth.err
}
