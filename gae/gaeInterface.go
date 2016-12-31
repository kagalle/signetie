package gae

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/go-errors/errors"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/phayes/freeport"
	"github.com/sourcegraph/go-webkit2/webkit2"
)

// "github.com/sqs/gojs"

// Authenticate is a class to handle the first part of the gae authentication process.
type Authenticate struct {
	*gtk.Window  // default member
	webView      *webkit2.WebView
	parentWindow *gtk.Window
	authWindow   *gtk.Window
	callback     AuthenticateComplete
	found        bool
	code         string
	cancelled    bool
	port         int
}

// AuthenticateComplete defines a user-supplied function to be called
// at completion of the process.
type AuthenticateComplete func(auth *Authenticate) bool

// NewAuthenticate is a constructor for Authenicate.
// Provide the gtk main window and the callback function defintion.
func NewAuthenticate(parentWindow *gtk.Window, authComplete AuthenticateComplete) *Authenticate {
	auth := new(Authenticate)
	auth.parentWindow = parentWindow
	auth.callback = authComplete
	auth.port = freeport.GetPort()
	return auth
}

// Code is a getter for the resulting code from the gae authentication process.
func (auth *Authenticate) Code() string {
	return auth.code
}

// Found is a getter; was the process successful.
func (auth *Authenticate) Found() bool {
	return auth.found
}

func (auth *Authenticate) setFound() {
	auth.found = true
	auth.cancelled = false // insure consistency
}

// Cancelled is a getter; was the process cancelled by the user, either
// by clicking the cancel button or by closing the authentate window prematurely.
func (auth *Authenticate) Cancelled() bool {
	return auth.cancelled
}
func (auth *Authenticate) setCancelled() {
	auth.found = false
	auth.cancelled = true // insure consistency
}

// Run is the main method which begins this first part of the authentation process.
func (auth *Authenticate) Run(scope string, clientID string) (err error) {
	// create window for browser
	auth.authWindow, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return errors.WrapPrefix(err, "Unable to create authenticate window", 0)
	}
	auth.authWindow.SetDefaultSize(850, 600)
	auth.authWindow.SetTitle("Authenticate")
	// make authwindow a modal child window of the main window
	auth.authWindow.SetModal(true)
	auth.authWindow.SetTransientFor(auth.parentWindow)
	// close the auth window if the parent window is closing
	auth.authWindow.SetDestroyWithParent(true)
	// add event handler for when the auth window is closing
	auth.authWindow.Connect("destroy", func() bool {
		if (!auth.found) && (!auth.cancelled) {
			auth.setCancelled()
		}
		auth.callback(auth)
		return false // let the window close
	})

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	if err != nil {
		log.Fatal("Unable to create vertical box:", err)
	}

	auth.webView = webkit2.NewWebView()
	cancelButton, err := gtk.ButtonNewWithLabel("Cancel")
	if err != nil {
		log.Fatal("Unable to create cancel button:", err)
	}
	cancelButton.Connect("clicked", func() {
		auth.setCancelled()
		auth.authWindow.Destroy() // which will trigger win destroy event
	})
	auth.webView.Connect("load-failed", func() {
	})
	// Wait until the browser window is full of data so that it will display
	// with size when added to the window.
	// Only do this for the first load, based on whether window has any content yet
	auth.webView.Connect("load-changed", func(_ *glib.Object, i int) {
		loadEvent := webkit2.LoadEvent(i)
		switch loadEvent {
		case webkit2.LoadFinished:
			childList := auth.authWindow.GetChildren()
			if childList.Length() == 0 {
				vbox.Add(auth.webView)
				vbox.Add(cancelButton)
				auth.authWindow.Add(vbox)
				auth.webView.Show()
				cancelButton.Show()
				vbox.Show()
			}
		}
	})
	state := randomDataBase64url(32)
	authURL := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?"+
		"scope=%s&"+
		"redirect_uri=http://127.0.0.1:%d&"+
		"response_type=code&"+
		"client_id=%s&"+
		"state=%s", scope, auth.port, clientID, state)
	// start the server to listen for the results
	// http://stackoverflow.com/a/6329459
	c1 := make(chan bool)
	go func() {
		http.HandleFunc("/",
			func(w http.ResponseWriter, r *http.Request) {
				// http://stackoverflow.com/a/25606975
				tempcode := r.URL.Query().Get("code")
				tempstate := r.URL.Query().Get("state")
				if (len(tempcode) != 0) && (state == tempstate) {
					auth.setFound()
					auth.code = tempcode
					// ask the main thread to close the auth window
					glib.IdleAdd(func() bool {
						auth.authWindow.Destroy() // which will trigger win destroy event
						return false
					})
				}
			})
		http.ListenAndServe(fmt.Sprintf(":%d", auth.port), nil)
		c1 <- true
	}()

	// Once the server is up ready to receive the result,
	// make the web call to open the gae authentation window.
	auth.authWindow.Show()
	auth.webView.LoadURI(authURL) // blocks until it loads - requires UI
	return nil                    // no error
}

// ref:  https://github.com/googlesamples/oauth-apps-for-windows.git
//     /OAuthDesktopApp/OAuthDesktopApp/MainWindow.xaml.cs
func randomDataBase64url(length int) string {
	var randomString string
	// create byte array of length <length> full of random data
	data := make([]byte, length)
	_, err := rand.Read(data)
	if err == nil {
		// base 64 encode the byte array, creating a string
		randomString = base64.StdEncoding.EncodeToString(data)
	}
	return randomString
}
