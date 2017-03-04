package login

import (
	"fmt"
	"log"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/kagalle/signetie/client_golang/gae/authenticate"
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

// I have a parent window, scope, clientID, and clientSecret
//  regardless of the code path needed
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
func (login *GaeLogin) Login(tokenSet *TokenSet) *TokenSet {
	// Determine if the accessToken in the tokenSet is stil valid.
	if (tokenSet != nil) &&
		(tokenSet.AccessToken != "") &&
		(tokenSet.ExpiresOn.After(time.Now())) {

		// The current accessToken should still work
		return tokenSet
	} else if (tokenSet != nil) && (tokenSet.RefreshToken != "") {
		// Attempt to use the refreshToken to get a new accessToken
		// TODO
		/*
					https://developers.google.com/identity/protocols/OAuth2InstalledApp#offline
					POST /oauth2/v4/token HTTP/1.1
			Host: www.googleapis.com
			Content-Type: application/x-www-form-urlencoded

			client_id=<your_client_id>&
			client_secret=<your_client_secret>&
			refresh_token=<refresh_token>&
			grant_type=refresh_token
		*/
		return nil
	} else {
		// Nothing in the current tokenSet that I can use - get a new one
		var tokenSet *TokenSet
		var err error
		port := freeport.GetPort() // 7777
		redirectURI := fmt.Sprintf("http://localhost:%d", port)
		auth := new(authenticate.Authenticate)
		auth.RequestAuthentication(login.parentWindow, login.scope, login.clientID,
			port, redirectURI, func(code string) {
				fmt.Printf("Code obtained %s\n", code)
				if code != "" {
					tokenSet, err = RequestAccessToken(code, login.clientID,
						login.clientSecret, redirectURI)

					if err != nil {
						log.Fatal("Unable to exchange code for token:", err)
					}
					if tokenSet != nil {
						tokenSet.Print()
					}
				}
			})
		return tokenSet
	}

}
