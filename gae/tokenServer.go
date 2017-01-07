package gae

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/braintree/manners"
)

// NewServer is a constructor to receive and process the result of the authentication request.
func NewTokenServer(port int, state string, authCode string, clientID string, clientSecret string) *manners.GracefulServer {
	mux := newTokenMux(state, authCode, clientID, clientSecret, port)
	server := manners.NewWithServer(&http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	})
	return server
}

// tokenMux implements http/Handler and wraps variables needed when responding to the response.
// https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/03.4.html
type tokenMux struct {
	state        string
	authCode     string
	clientID     string
	clientSecret string
	port         int
}

// newTokenMux is a constructor to create tokenMux.
func newTokenMux(state string, authCode string, clientID string, clientSecret string, port int) *tokenMux {
	mux := new(tokenMux)
	mux.state = state
	mux.authCode = authCode
	mux.clientID = clientID
	mux.clientSecret = clientSecret
	mux.port = port
	return mux
}

func (p *tokenMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// document.getElementById('token_form').submit();
	fmt.Printf("token server method=%s  url=%s\n", r.Method, r.URL.Path)
	if r.URL.Path == "/load" {
		fmt.Printf("server load\n")
		w.Write([]byte(fmt.Sprintf("<html><body><form id=\"token_form\" action=\"https://www.googleapis.com/oauth2/v4/token\" method=\"post\">"+
			"<input type=\"text\" name=\"code\" value=\"%s\">"+
			"<input type=\"text\" name=\"client_id\" value=\"%s\">"+
			"<input type=\"text\" name=\"client_secret\" value=\"%s\">"+
			"<input type=\"text\" name=\"redirect_uri\" value=\"%s\">"+
			"<input type=\"text\" name=\"grant_type\" value=\"authorization_code\">"+
			// button name cannot be 'submit'
			// http://stackoverflow.com/a/39149592
			"<input type=\"submit\" name=\"submitButton\" value=\"submit\">"+
			"</form></body></html>",
			// TODO: I'm POSTing to the wrong place...
			p.authCode, p.clientID, p.clientSecret, fmt.Sprintf("http://localhost:%d", p.port))))
		//"http://localhost")))
		return
	} else if r.URL.Path == "/" {
		fmt.Printf("server post\n")
		// defer resp.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Unable to read get token response", err)
		}
		fmt.Printf("Token obtained %s\n", string(body[:]))
		// return string(body[:]), nil

		// r.Body.Read()
		// a := r.URL.Query().Get("access_token")
		// a := r.URL.Query().Get("id_token")
		// a := r.URL.Query().Get("refresh_token")
		// a := r.URL.Query().Get("expires_in")
		// a := r.URL.Query().Get("")
		// a := r.URL.Query().Get("")
		// a := r.URL.Query().Get("")
		// a := r.URL.Query().Get("")
		// tempstate := r.URL.Query().Get("state")
		// if p.state == tempstate {
		// 	if len(tempcode) != 0 {
		// 		p.auth.setFound()
		// 		p.auth.code = tempcode
		// 	} else {
		// 		p.auth.err = errors.Errorf("token not returned from service")
		// 	}
		// } else {
		// 	p.auth.err = errors.Errorf("token received from incorrrect session: original=%s  returned=%s", p.state, tempstate)
		// }
		return
	}
	fmt.Printf("server - url not found\n")
	http.NotFound(w, r)
	return
}
