package authenticate

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/phayes/freeport"
)

// ProcessCallback is a callback function to communicate both from the server to here, and from here to the caller.
type ProcessCallback func(output *AuthOutput)

type Input struct {
	scope    string
	clientID string
	// clientSecret string
	port        int
	redirectURI string
}

// constructor to create ageLogin object.
func NewInput(scope string, clientID string /*, clientSecret string*/) *Input {
	input := new(Input)
	input.scope = scope
	input.clientID = clientID
	// input.clientSecret = clientSecret
	input.port = freeport.GetPort() // 7777
	return input
}

func (input *Input) RedirectURI() string {
	if input.redirectURI == "" {
		input.redirectURI = fmt.Sprintf("http://localhost:%d", input.port)
	}
	return input.redirectURI
}

// func newOutput(code string, redirectURI string) *AuthOutput {
// 	output := new(AuthOutput)
// 	output.code = code
// 	output.redirectURI = redirectURI
// 	return output
// }

// AuthOutput is a class to handle the first part of the gae authentication process.
type AuthOutput struct {
	// found       bool
	code string
	// cancelled   bool
	redirectURI string
	err         *errors.Error
}

// Code is a getter for the resulting code from the gae authentication process.
func (output *AuthOutput) Code() string {
	return output.code
}

func (output *AuthOutput) RedirectURI() string {
	return output.redirectURI
}

func (output *AuthOutput) SetRedirectURI(redirectURI string) {
	output.redirectURI = redirectURI
}

// Found is a getter; was the process successful.
// func (output *AuthOutput) Found() bool {
// 	return output.found
// }

// func (output *AuthOutput) setFound() {
// 	output.found = true
// 	output.cancelled = false // insure consistency
// }

// Cancelled is a getter; was the process cancelled by the user, either
// by clicking the cancel button or by closing the authentate window prematurely.
// func (output *AuthOutput) Cancelled() bool {
// 	return output.cancelled
// }
// func (output *AuthOutput) setCancelled() {
// 	output.found = false
// 	output.cancelled = true // insure consistency
// }

// func (output *AuthOutput) Error() error {
// 	return output.err
// }
