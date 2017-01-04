package gae

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/go-errors/errors"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/phayes/freeport"
	"github.com/sourcegraph/go-webkit2/webkit2"
	"github.com/sqs/gojs"
)

// AuthenticateComplete defines a user-supplied function to be called
// at completion of the process.
type AuthenticateComplete func(auth *Authenticate)

// RequestAuthentication is the main method which begins this first part of the authentation process.
func RequestAuthentication(parentWindow *gtk.Window, scope string, clientID string, authCompleteCallback AuthenticateComplete) (err error) {

	port := freeport.GetPort()
	auth := new(Authenticate)
	state := RandomDataBase64url(32)
	var authWindow *gtk.Window
	// create window for browser
	authWindow, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		authWindow.Destroy() // which will trigger win destroy event
		return errors.WrapPrefix(err, "Unable to create authenticate window", 0)
	}

	server := NewAuthServer(port, state, auth, authWindow)

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
		server.BlockingClose()
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
		auth.err = errors.Errorf("Unable to load authentication page")
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

	// start the server to listen for the results
	// http://stackoverflow.com/a/6329459
	go func() {
		server.ListenAndServe()
	}()

	// Once the server is up ready to receive the result,
	// make the web call to open the gae authentation window.
	authWindow.Show()
	// Note that although this blocks until the page is loaded,
	// it doesn't block until the user completes the whole process.
	authURL := new(url.URL)
	authURL.Scheme = "https"
	authURL.Host = "accounts.google.com"
	authURL.Path = "/o/oauth2/v2/auth"
	authURLParams := url.Values{}
	authURLParams.Set("scope", scope)
	authURLParams.Set("redirect_uri", fmt.Sprintf("http://127.0.0.1:%d", port))
	authURLParams.Set("response_type", "code")
	authURLParams.Set("client_id", clientID)
	authURLParams.Set("state", state)
	authURL.RawQuery = authURLParams.Encode()
	webView.LoadURI(authURL.String()) // blocks until it loads - requires UI
	return nil                        // no error
}

// RequestAccessToken takes code obtained from previous step and converts it into a token.
func RequestAccessToken(authCode string, clientID string, clientSecret string) (string, error) /* *AccessToken */ {
	// params := url.Values{}
	port := freeport.GetPort()
	state := RandomDataBase64url(32)
	server := NewTokenServer(port, state, authCode, clientID, clientSecret)
	go func() {
		server.ListenAndServe()
	}()

	webView := webkit2.NewWebView()

	webView.Connect("load-changed", func(_ *glib.Object, i int) {
		loadEvent := webkit2.LoadEvent(i)
		switch loadEvent {
		case webkit2.LoadFinished:
			fmt.Printf("JS\n")
			webView.RunJavaScript("document.getElementById('token_form').submit();", func(g *gojs.Value, err error) {
				// java script callback
				fmt.Printf("Javascript call back: %s\n", err)
				return
			})
		}
	})
	fmt.Printf("load\n")
	webView.LoadURI(fmt.Sprintf("http://localhost:%d/load", port))

	// params.Set("code", authCode)
	// params.Set("client_id", clientID)
	// params.Set("client_secret", clientSecret)
	// // The redirect_uri comes from the "Download JSON" button in the edit client_id screen in the API Manager.
	// // Apparently for type "other" it can't be edited - you are just assigned this a a usable value.
	// params.Set("redirect_uri", fmt.Sprintf("http://localhost:%d", port))
	// params.Set("grant_type", "authorization_code")
	// resp, err := http.PostForm("https://www.googleapis.com/oauth2/v4/token", params)
	// if err != nil {
	// 	return "", errors.WrapPrefix(err, "Unable to convert code into token", 0)
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return "", errors.WrapPrefix(err, "Unable to read get token response", 0)
	// }

	// return string(body[:]), nil
	// Need to parse the JSON response and look for {"error": "something", "error_description": "something has detail"}
	return "", nil
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
