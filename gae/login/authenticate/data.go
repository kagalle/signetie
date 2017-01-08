package gae

import (
	"github.com/go-errors/errors"
	"github.com/phayes/freeport"
)

type Input struct {
	scope    string
	clientID string
	// clientSecret string
	port int
}

// constructor to create ageLogin object.
func newInput(scope string, clientID string /*, clientSecret string*/) *Input {
	input := new(Input)
	input.scope = scope
	input.clientID = clientID
	// input.clientSecret = clientSecret
	input.port = freeport.GetPort() // 7777
	return input
}

// func newOutput(code string, redirectURI string) *Output {
// 	output := new(Output)
// 	output.code = code
// 	output.redirectURI = redirectURI
// 	return output
// }

// Output is a class to handle the first part of the gae authentication process.
type Output struct {
	// found       bool
	code string
	// cancelled   bool
	redirectURI string
	err         *errors.Error
}

// Code is a getter for the resulting code from the gae authentication process.
func (output *Output) Code() string {
	return output.code
}

func (output *Output) RedirectURI() string {
	return output.redirectURI
}

// Found is a getter; was the process successful.
// func (output *Output) Found() bool {
// 	return output.found
// }

// func (output *Output) setFound() {
// 	output.found = true
// 	output.cancelled = false // insure consistency
// }

// Cancelled is a getter; was the process cancelled by the user, either
// by clicking the cancel button or by closing the authentate window prematurely.
// func (output *Output) Cancelled() bool {
// 	return output.cancelled
// }
// func (output *Output) setCancelled() {
// 	output.found = false
// 	output.cancelled = true // insure consistency
// }

// func (output *Output) Error() error {
// 	return output.err
// }
