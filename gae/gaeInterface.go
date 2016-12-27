package gae

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gotk3/gotk3/gtk"
	"github.com/sourcegraph/go-webkit2/webkit2"
)

type Authenticate struct {
	*gtk.Window // default member
	callback    AuthenticateComplete
	found       bool
	code        string
	cancelled   bool
}

type AuthenticateComplete func(auth *Authenticate) bool

func NewAuthenticate(parentWindow *gtk.Window, authComplete AuthenticateComplete) *Authenticate {
	auth := new(Authenticate)
	auth.Window = parentWindow
	auth.callback = authComplete

	// create new signal
	//http://zetcode.com/gui/pygtk/signals/
	//	GObject.type_register(GaeAuthenicate)
	//	GObject.signal_new("authenticate_complete", GaeAuthenicate, GObject.SIGNAL_RUN_FIRST, GObject.TYPE_NONE, ())

	return auth
}

func (auth *Authenticate) GetCode() string {
	return auth.code
}

func (auth *Authenticate) GetFound() bool {
	return auth.found
}

func (auth *Authenticate) GetCancelled() bool {
	return auth.cancelled
}

//	"cmd/internal/pprof/tempfile"

// Authenicate returns the code if successful, or non-nil error if not (code, err).
func (auth *Authenticate) Run(scope string, clientID string) (string, error) {
	// create window for browser
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create authenticate window:", err)
	}
	win.SetTitle("Authenticate")
	win.SetModal(true)
	win.SetTransientFor(auth.Window)
	win.SetDestroyWithParent(true)
	win.Connect("destroy", func() bool {
		if (!auth.found) && (!auth.cancelled) {
			auth.cancelled = true
		}
		auth.callback(auth)
		fmt.Printf("B1")
		return false // let the window close
	})

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	if err != nil {
		log.Fatal("Unable to create vertical box:", err)
	}
	win.Add(vbox)

	webView := webkit2.NewWebView()
	vbox.Add(webView)

	cancelButton, err := gtk.ButtonNewWithLabel("Cancel")
	if err != nil {
		log.Fatal("Unable to create cancel button:", err)
	}
	cancelButton.Connect("clicked", func() {
		auth.cancelled = true
		auth.found = false // insure consistency
		win.Destroy()      // which will trigger win destroy event
	})
	vbox.Add(cancelButton)
	win.ShowAll()

	// form string
	var code string
	authURL := "http://www.google.com"

	// authURL := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?"+
	// 	"scope=%s&"+
	// 	"redirect_uri=http://127.0.0.1:8146&"+
	// 	"response_type=code&"+
	// 	"client_id=%s", scope, clientID)

	// start the server to listen for the results
	// http://stackoverflow.com/a/6329459
	fmt.Printf("B2")
	c := make(chan bool)
	go func() {
		http.HandleFunc("/",
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
				// http://stackoverflow.com/a/25606975
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
	webView.LoadURI(authURL)
	fmt.Printf("D")

	// wait for the blocking function to finish if it hasn't already
	<-c
	fmt.Printf("E")

	return code, err
}
