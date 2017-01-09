package authenticate

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/braintree/manners"
	"github.com/kagalle/signetie/client_golang/gae/login/util"
)

type AuthServer struct {
	srv    *manners.GracefulServer
	input  *Input
	output *AuthOutput
	// authMux implements http/Handler and wraps variables needed when responding to the response.
	// https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/03.4.html
	// type authMux struct {
	state                string
	authCompleteCallback ProcessCallback
	// redirectURI          string
	// }
}

// NewAuthServer is a constructor to receive and process the result of the authentication request.
func NewAuthServer(input *Input, authCompleteCallback ProcessCallback) *AuthServer {
	output := new(AuthOutput)
	authServer := new(AuthServer)
	// mux := newAuthMux(input, authCompleteCallback)
	server := manners.NewWithServer(&http.Server{
		Addr:           fmt.Sprintf(":%d", input.port),
		Handler:        authServer,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	})
	state := util.RandomDataBase64url(32)
	authServer.srv = server
	authServer.input = input
	authServer.output = output
	authServer.state = state
	authServer.authCompleteCallback = authCompleteCallback
	return authServer
}

func (p *AuthServer) FormAuthURL() string {
	authURL := new(url.URL)
	authURL.Scheme = "https"
	authURL.Host = "accounts.google.com"
	authURL.Path = "/o/oauth2/v2/auth"
	authURLParams := url.Values{}
	authURLParams.Set("scope", p.input.scope)
	authURLParams.Set("redirect_uri", p.input.RedirectURI())
	authURLParams.Set("response_type", "code")
	authURLParams.Set("client_id", p.input.clientID)
	authURLParams.Set("state", p.state)
	authURL.RawQuery = authURLParams.Encode()
	return authURL.String()
}

func (p *AuthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("output server method=%s  url=%s\n", r.Method, r.URL.Path)
	if r.URL.Path == "/" {
		tempcode := r.URL.Query().Get("code")
		tempstate := r.URL.Query().Get("state")
		if p.state == tempstate {
			if len(tempcode) != 0 {
				p.output.code = tempcode
				p.output.SetRedirectURI(p.input.RedirectURI()) // pass this through
				p.authCompleteCallback(p.output)
			} else {
				fmt.Printf("Authentication code not returned from service")
			}
		} else {
			fmt.Printf("Authentication received from incorrrect session: original=%s  returned=%s", p.state, tempstate)
		}
	}
	//	else {
	//		http.NotFound(w, r)
	//	}
	return
}
