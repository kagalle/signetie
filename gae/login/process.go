package login

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
	"github.com/kagalle/signetie/client_golang/gae/login/authenticate"
	"github.com/kagalle/signetie/client_golang/gae/login/codetoken"
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

func NewGaeLogin(parentWindow *gtk.Window, scope string, clientID string, clientSecret string) *GaeLogin {
	login := new(GaeLogin)
	login.parentWindow = parentWindow
	login.scope = scope
	login.clientID = clientID
	login.clientSecret = clientSecret
	return login
}

func (login *GaeLogin) GaeLoginWithoutRefreshToken() {
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
