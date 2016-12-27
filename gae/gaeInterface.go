package gae

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

//	"cmd/internal/pprof/tempfile"

// returns the code if successful, or non-nil error if not.
// returns (code, err)
func Authenicate(scope string, clientId string) (string, error) {
	// form string
	var code string
	authUrl := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?"+
		"scope=%s&"+
		"redirect_uri=http://127.0.0.1:8146&"+
		"response_type=code&"+
		"client_id=%s", scope, clientId)
	// start the server to listen for the results
	// http://stackoverflow.com/a/6329459
	fmt.Printf("B")
	c := make(chan bool)
	go func() {
		http.HandleFunc("/",
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
				tempcode := r.URL.Query().Get("code")
				if len(tempcode) != 0 {
					code = tempcode
				}
				// io.WriteString(w, "Hello world!")
			})
		log.Fatal(http.ListenAndServe(":8146", nil)) // TODO make the port a parameter
		c <- true
	}()
	fmt.Printf("C")

	// do some other stuff here while the blocking function runs
	// make the call
	// http://stackoverflow.com/a/25344458
	timeout := time.Duration(15 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	_, err := client.Get(authUrl)
	fmt.Printf("D")
	// webView.load_uri(auth_url)

	// wait for the blocking function to finish if it hasn't already
	<-c
	fmt.Printf("E")

	return code, err
}
