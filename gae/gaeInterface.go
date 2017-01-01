package gae

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/go-errors/errors"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/phayes/freeport"
	"github.com/sourcegraph/go-webkit2/webkit2"
)

// AuthenticateComplete defines a user-supplied function to be called
// at completion of the process.
type AuthenticateComplete func(auth *Authenticate)

// RequestAuthentication is the main method which begins this first part of the authentation process.
func RequestAuthentication(parentWindow *gtk.Window, scope string, clientID string, authCompleteCallback AuthenticateComplete) (err error) {
	auth := new(Authenticate)
	var authWindow *gtk.Window
	// create window for browser
	authWindow, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		authWindow.Destroy() // which will trigger win destroy event
		return errors.WrapPrefix(err, "Unable to create authenticate window", 0)
	}
	authWindow.SetDefaultSize(850, 600)
	authWindow.SetTitle("Authenticate")
	// make authwindow a modal child window of the main window
	authWindow.SetModal(true)
	authWindow.SetTransientFor(parentWindow)
	// close the auth window if the parent window is closing
	authWindow.SetDestroyWithParent(true)
	// add event handler for when the auth window is closing
	authWindow.Connect("destroy", func() bool {
		if (!auth.found) && (!auth.cancelled) && (auth.err == nil) {
			auth.setCancelled()
		}
		authCompleteCallback(auth)
		return false // let the window close
	})
	// add box for layout
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	if err != nil {
		authWindow.Destroy() // which will trigger win destroy event
		return errors.WrapPrefix(err, "Unable to create vertical box", 0)
	}
	// create the webview browser
	webView := webkit2.NewWebView()
	cancelButton, err := gtk.ButtonNewWithLabel("Cancel")
	if err != nil {
		authWindow.Destroy() // which will trigger win destroy event
		return errors.WrapPrefix(err, "Unable to create cancel button", 0)
	}
	cancelButton.Connect("clicked", func() {
		auth.setCancelled()
		authWindow.Destroy() // which will trigger win destroy event
	})
	webView.Connect("load-failed", func() {
		auth.err = errors.WrapPrefix(err, "Unable to load authentication page", 0)
		authWindow.Destroy() // which will trigger win destroy event
	})
	// Wait until the browser window is full of data so that it will display
	// with size when added to the window.
	// Only do this for the first load, based on whether window has any content yet
	webView.Connect("load-changed", func(_ *glib.Object, i int) {
		loadEvent := webkit2.LoadEvent(i)
		switch loadEvent {
		case webkit2.LoadFinished:
			childList := authWindow.GetChildren()
			if childList.Length() == 0 {
				vbox.Add(webView)
				vbox.Add(cancelButton)
				authWindow.Add(vbox)
				webView.Show()
				cancelButton.Show()
				vbox.Show()
			}
		}
	})
	port := freeport.GetPort()
	state := RandomDataBase64url(32)
	authURL := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?"+
		"scope=%s&"+
		"redirect_uri=http://127.0.0.1:%d&"+
		"response_type=code&"+
		"client_id=%s&"+
		"state=%s", scope, port, clientID, state)
	// start the server to listen for the results
	// http://stackoverflow.com/a/6329459
	c1 := make(chan bool)
	go func() {
		http.HandleFunc("/",
			func(w http.ResponseWriter, r *http.Request) {
				// http://stackoverflow.com/a/25606975
				tempcode := r.URL.Query().Get("code")
				tempstate := r.URL.Query().Get("state")
				if state == tempstate {
					if len(tempcode) != 0 {
						auth.setFound()
						auth.code = tempcode
					} else {
						auth.err = errors.WrapPrefix(err, "Authentication code not returned from service", 0)
					}
				} else {
					auth.err = errors.WrapPrefix(err, "Authentication received from incorrrect session", 0)
				}
				// ask the main thread to close the auth window
				glib.IdleAdd(func() bool {
					authWindow.Destroy() // which will trigger win destroy event
					return false         // only have IdleAdd() call this once
				})
			})
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		c1 <- true
	}()

	// Once the server is up ready to receive the result,
	// make the web call to open the gae authentation window.
	authWindow.Show()
	// Note that although this blocks until the page is loaded,
	// it doesn't block until the user completes the whole process.
	webView.LoadURI(authURL) // blocks until it loads - requires UI
	return nil               // no error
}

// RandomDataBase64url creates a base64 encoded string of length bytes.
// ref:  https://github.com/googlesamples/oauth-apps-for-windows.git
//     /OAuthDesktopApp/OAuthDesktopApp/MainWindow.xaml.cs
func RandomDataBase64url(length int) string {
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
