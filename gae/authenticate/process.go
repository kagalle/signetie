package authenticate

import (
	"fmt"
	"net/url"

	"github.com/go-errors/errors"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/kagalle/signetie/client_golang/gae/util"
	"github.com/sourcegraph/go-webkit2/webkit2"
)

// AuthCompleteCallback defines the user-supplied method to be called with the auth code.
type AuthCompleteCallback func(code string)

// Authenticate is a class that does the authenicate with GAE, the result of which is a string code.
type Authenticate string

// RequestAuthentication is the main method which begins this first part of the authentation process.
func (p Authenticate) RequestAuthentication(parentWindow *gtk.Window, scope string, clientID string,
	port int, redirectURI string, callback AuthCompleteCallback) error {

	var err error // return value
	var authWindow *gtk.Window
	// create window for browser
	authWindow, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		authWindow.Destroy() // which will trigger win destroy event
		return errors.WrapPrefix(err, "Unable to create authenticate window", 0)
	}

	state := util.RandomDataBase64url(32)
	authServer := NewAuthServer(port, state, func(newCode string) {
		p = Authenticate(newCode) // type assert that newCode, a string, is compatible with Authenicate
		// this is the server's thread, so ask the main thread for this
		glib.IdleAdd(func() bool {
			authWindow.Destroy() // which will trigger win destroy event
			return false         // only have IdleAdd() call this once
		})
	})

	authWindow.SetDefaultSize(850, 600)
	authWindow.SetTitle("Authenticate")
	// make authwindow a modal child window of the main window
	authWindow.SetModal(true)
	authWindow.SetTransientFor(parentWindow)
	// close the data window if the parent window is closing
	authWindow.SetDestroyWithParent(true)
	// add event handler for when the data window is closing
	authWindow.Connect("destroy", func() bool {
		// if (!data.found) && (!data.cancelled) && (data.err == nil) {
		// 	data.setCancelled()
		// }
		authServer.srv.BlockingClose()
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
		//data.setCancelled()
		authWindow.Destroy() // which will trigger win destroy event
	})
	webView.Connect("load-failed", func() {
		fmt.Printf("Unable to load authentication page")
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

	// Once the server is up ready to receive the result,
	// make the web call to open the gae authentation window.
	authWindow.Show()
	// Note that although this blocks until the page is loaded,
	// it doesn't block until the user completes the whole process.
	url := formAuthURL(scope, clientID, redirectURI, state)
	fmt.Printf("Auth url:%s", url)
	webView.LoadURI(url)

	// start the server to listen for the results
	// http://stackoverflow.com/a/6329459
	go func() {
		authServer.srv.ListenAndServe()
	}()
	return nil // no error
}

func formAuthURL(scope string, clientID string, redirectURI string, state string) string {
	authURL := new(url.URL)
	authURL.Scheme = "https"
	authURL.Host = "accounts.google.com"
	authURL.Path = "/o/oauth2/v2/auth"
	authURLParams := url.Values{}
	authURLParams.Set("scope", scope)
	authURLParams.Set("redirect_uri", redirectURI)
	authURLParams.Set("response_type", "code")
	authURLParams.Set("client_id", clientID)
	authURLParams.Set("state", state)
	authURL.RawQuery = authURLParams.Encode()
	return authURL.String()
}
