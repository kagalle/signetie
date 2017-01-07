package gae

import (
	"fmt"
	"net/http"
	"time"

	"github.com/braintree/manners"
	"github.com/go-errors/errors"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// AuthServer is a constructor to receive and process the result of the authentication request.
func NewAuthServer(port int, state string, auth *Authenticate, authWindow *gtk.Window) *manners.GracefulServer {
	mux := newAuthMux(state, auth, authWindow)
	server := manners.NewWithServer(&http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	})
	return server
}

// authMux implements http/Handler and wraps variables needed when responding to the response.
// https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/03.4.html
type authMux struct {
	state      string
	auth       *Authenticate
	authWindow *gtk.Window
}

// NewAuthMux is a constructor to create MyMux.
func newAuthMux(state string, auth *Authenticate, authWindow *gtk.Window) *authMux {
	mux := new(authMux)
	mux.state = state
	mux.auth = auth
	mux.authWindow = authWindow
	return mux
}

func (p *authMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("auth server method=%s  url=%s\n", r.Method, r.URL.Path)
	if r.URL.Path == "/" {
		tempcode := r.URL.Query().Get("code")
		tempstate := r.URL.Query().Get("state")
		if p.state == tempstate {
			if len(tempcode) != 0 {
				p.auth.setFound()
				p.auth.code = tempcode
			} else {
				p.auth.err = errors.Errorf("Authentication code not returned from service")
			}
		} else {
			p.auth.err = errors.Errorf("Authentication received from incorrrect session: original=%s  returned=%s", p.state, tempstate)
		}
		// ask the main thread to close the auth window
		glib.IdleAdd(func() bool {
			p.authWindow.Destroy() // which will trigger win destroy event
			return false           // only have IdleAdd() call this once
		})
		return
	}
	http.NotFound(w, r)
	return
}
