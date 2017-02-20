package login

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gotk3/gotk3/gtk"
	"github.com/kagalle/signetie/client_golang/gae/gaeAccessToken"
	"github.com/kagalle/signetie/client_golang/gae/gaeAuthenticate"
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
func (login *GaeLogin) Login(tokenSet *gaeAccessToken.TokenSet) *gaeAccessToken.TokenSet {
	// Determine if the accessToken in the tokenSet is stil valid.
	if (tokenSet != nil) &&
		(tokenSet.AccessToken != "") &&
		(tokenSet.ExpiresOn.After(time.Now())) {

		// The current accessToken should still work
		return tokenSet
	} else if (tokenSet != nil) && (tokenSet.RefreshToken != "") {
		// Attempt to use the refreshToken to get a new accessToken
		// TODO
		return nil
	} else {
		// Nothing in the current tokenSet that I can use - get a new one
		var tokenSet *gaeAccessToken.TokenSet
		var err error
		port := freeport.GetPort() // 7777
		redirectURI := fmt.Sprintf("http://localhost:%d", port)
		gaeAuthenticate.RequestAuthentication(login.parentWindow, login.scope, login.clientID,
			port, redirectURI, func(code string) {
				logrus.WithField("Code obtained", code).Debug("gaeAuthenticate.RequestAuthentication() result")
				if code != "" {
					tokenSet, err = gaeAccessToken.RequestAccessToken(code, login.clientID,
						login.clientSecret, redirectURI)

					if err != nil {
						logrus.WithError(err).Error("Unable to exchange code for token")
					}
					if tokenSet != nil {
						tokenSet.Log()
					}
				}
			})
		return tokenSet
	}

}
