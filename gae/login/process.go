package login

import (
	"fmt"
	"log"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/kagalle/signetie/client_golang/gae/login/authenticate"
	"github.com/kagalle/signetie/client_golang/gae/login/codetoken"
	"github.com/phayes/freeport"
)

// authenticateComplete defines a user-supplied function to be called
// at completion of the process.
// type loginComplete func(data *Data, port int)

type GaeLogin struct {
	parentWindow *gtk.Window
	scope        string
	clientID     string
	clientSecret string
}

// I have a parent window regardless of the path
// I have the scope regardless of the path needed
// I have a clientID regardless of the path needed
// I have a clientSecret regardless of the path needed
func NewGaeLogin(parentWindow *gtk.Window, scope string, clientID string, clientSecret string) *GaeLogin {
	login := new(GaeLogin)
	login.parentWindow = parentWindow
	login.scope = scope
	login.clientID = clientID
	login.clientSecret = clientSecret
	return login
}

// Login calls appropriate code in order to get a valid tokenSet.
// The input tokenSet may be nil.
func (login *GaeLogin) Login(tokenSet TokenSet) TokenSet {
	// Determine if the accessToken in the tokenSet is stil valid.
	if (tokenSet != nil) &&
		(tokenSet.accessToken != "") &&
		(tokenSet.expiresOn != nil) &&
		(tokenSet.expiresOn.After(time.Now())) {

		// The current accessToken should still work
		return tokenSet
	} else if (tokenSet != nil) && (tokenSet.refreshToken != "") {
		// Attempt to use the refreshToken to get a new accessToken
		// TODO
	} else {
		// Nothing in the current tokenSet that I can use - get a new one
		var authcode string
		port := freeport.GetPort() // 7777
		redirectURI := 
		auth := new authenticate.Authenticate()
		auth.RequestAuthentication(login.parentWindow, login.scope, login.clientID,
			port, redirectURI, func(code string) {
				authcode = code
			})
	}

	// Determine if we need to authenicate again.

	// If I was given a refresh token, then attempt to use that.
	if (tokenSet.refreshToken != "") && (token) {
		// TODO
	} else {
		//
	}
	input := authenticate.NewInput(login.scope, login.clientID)

	
	authenticate.RequestAuthentication(login.parentWindow, input, login.afterRequestAuthentication)

}

func (login *GaeLogin) GaeLoginWithRefreshToken() {

}

func (login *GaeLogin) afterRequestAuthentication(output *authenticate.AuthOutput) {
	if output.Code() != "" {
		fmt.Printf("Code obtained %s\n", output.Code())
		// authCode string, clientID string, clientSecret string, redirectURI string
		accessToken, err := codetoken.RequestAccessToken(output.Code(), login.clientID,
			login.clientSecret, output.RedirectURI())

		if err != nil {
			log.Fatal("Unable to exchange code for token:", err)
		}
		fmt.Printf("token_response:%s", accessToken)
	}
}
