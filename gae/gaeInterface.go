package gae

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/sourcegraph/go-webkit2/webkit2"
	"github.com/sqs/gojs"
)

// "github.com/sqs/gojs"

type Authenticate struct {
	*gtk.Window // default member
	webView     *webkit2.WebView
	authWindow  *gtk.Window
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

func (auth *Authenticate) Setup() {

}

// Authenicate returns the code if successful, or non-nil error if not (code, err).
func (auth *Authenticate) Run(scope string, clientID string) {
	// create window for browser
	var err error
	auth.authWindow, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create authenticate window:", err)
	}
	auth.authWindow.SetDefaultSize(850, 600)
	auth.authWindow.SetTitle("Authenticate")
	auth.authWindow.SetModal(true)
	auth.authWindow.SetTransientFor(auth.Window)
	auth.authWindow.SetDestroyWithParent(true)
	auth.authWindow.Connect("destroy", func() bool {
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
	vbox.Show()

	auth.webView = webkit2.NewWebView()
	// auth.webView.SetVisible(true)
	auth.webView.Show()
	cancelButton, err := gtk.ButtonNewWithLabel("Cancel")
	if err != nil {
		log.Fatal("Unable to create cancel button:", err)
	}
	cancelButton.Connect("clicked", func() {
		auth.cancelled = true
		auth.found = false        // insure consistency
		auth.authWindow.Destroy() // which will trigger win destroy event
	})
	cancelButton.Show()

	auth.webView.Connect("load-failed", func() {
		fmt.Println("Load failed.")
	})
	auth.webView.Connect("load-changed", func(_ *glib.Object, i int) {
		loadEvent := webkit2.LoadEvent(i)
		switch loadEvent {
		case webkit2.LoadFinished:
			fmt.Println("Load finished.")
			vbox.Add(auth.webView)
			vbox.Add(cancelButton)
			auth.authWindow.Add(vbox)
			fmt.Printf("C4")
			//auth.authWindow.Add(auth.webView)
			fmt.Printf("Title: %q\n", auth.webView.Title())
			fmt.Printf("URI: %s\n", auth.webView.URI())
			auth.webView.RunJavaScript("window.location.hostname", func(val *gojs.Value, err error) {
				if err != nil {
					fmt.Println("JavaScript error.")
				} else {
					fmt.Printf("Hostname (from JavaScript): %q\n", val)
				}
				// gtk.MainQuit()
			})
		}
	})

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
	c1 := make(chan bool)
	go func() {
		fmt.Printf("B3")
		http.HandleFunc("/",
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Printf("B4")
				fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
				// http://stackoverflow.com/a/25606975
				tempcode := r.URL.Query().Get("code")
				if len(tempcode) != 0 {
					code = tempcode
				}
				// io.WriteString(w, "Hello world!")
				fmt.Printf("B5")
			})
		fmt.Printf("B6")
		http.ListenAndServe(":8146", nil) // TODO make the port a parameter
		fmt.Printf("B7")
		c1 <- true
		// return false // ask IdleAdd() to not call this anonymous function again
	}()
	fmt.Printf("C1")

	// do some other stuff here while the blocking function runs
	// make the call
	// c2 := make(chan bool)
	// go func() {
	auth.authWindow.Show()
	glib.IdleAdd(func() bool {
		fmt.Printf("C2")
		auth.webView.LoadURI(authURL) // blocks until it loads - requires UI
		fmt.Printf("C3")
		// c2 <- true
		return false
	})
	// }()
	fmt.Printf("D")

	// wait for the blocking function to finish if it hasn't already
	// <-c
	fmt.Printf("E")

	// return code, err
}
