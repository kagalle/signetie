package login

import (
	"fmt"

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
	tokenState := tokenSet.GetState()
	switch tokenState.GetCurrent() {
	case gaeAccessToken.Active:
		// The current accessToken should still work
		return tokenSet
	case gaeAccessToken.Refresh:
		// Attempt to use the refreshToken to get a new accessToken
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
		var err error
		tokenSet, err = gaeAccessToken.RefreshAccessToken(tokenSet, login.clientID, login.clientSecret)
		if err != nil {
			logrus.WithError(err).Error("Unable to refresh token")
		}
		if tokenSet != nil {
			tokenSet.Log("Refresh tokenSet")
		}
		return tokenSet
	default:
		// There is nothing in the current tokenSet that I can use - get a new one
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
						tokenSet.Log("Create new tokenSet")
					}
				}
			})
		return tokenSet
	}
}
