package gaeAuthenticateHttpServer

import (
	"fmt"
	"net/http"
	"time"

	"github.com/braintree/manners"
)

type authServerCallback func(code string)

// authMux implements http/Handler and wraps variables needed when responding to the response.
// https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/03.4.html
type AuthServer struct {
	// the http server
	gracefulServer *manners.GracefulServer
	// state is the value sent into the URL when the request is made;
	// the server will compare the value it gets from the URL to this value
	// to make sure it is handling the correct request.
	state string
	// callback is any function that accepts a string value; to be supplied by
	// the caller so that the server can return the resulting code obtained from
	// the API.
	callback authServerCallback
}

// NewAuthServer is a constructor to receive and process the result of the authentication request.
func NewAuthServer(port int, state string, callback authServerCallback) *AuthServer {
	authServer := new(AuthServer)
	authServer.gracefulServer = manners.NewWithServer(&http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        authServer,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	})
	authServer.state = state
	authServer.callback = callback
	return authServer
}

// Satisfies the Handler interface: ServeHTTP(ResponseWriter, *Request)
func (p *AuthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("output server method=%s  url=%s\n", r.Method, r.URL.Path)
	if r.URL.Path == "/" {
		tempcode := r.URL.Query().Get("code")
		tempstate := r.URL.Query().Get("state")
		if p.state == tempstate {
			if len(tempcode) != 0 {
				p.callback(tempcode)
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

func (p *AuthServer) BlockingClose() {
	p.gracefulServer.BlockingClose()
}

func (p *AuthServer) ListenAndServe() {
	p.gracefulServer.ListenAndServe()
}
