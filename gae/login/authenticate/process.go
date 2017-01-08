package gae

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-errors/errors"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/sourcegraph/go-webkit2/webkit2"
)

// TODO: the authenticate part of this might as well move into Data, and Data be changed to authenciateResult.
// that leaves this class to call out to do the authenticate if needed, followed by the code/token exchange.

// authenticateComplete defines a user-supplied function to be called
// at completion of the process.
// type loginComplete func(data *Data, port int)

// TODO: needs to return a structure that contains
//  access_token, id_token (decode the JWT into user fields needed), refresh_token, expires_in (or convert to expires_on:date/time).
func GaeLoginWithoutRefreshToken(scope string, clientID string, clientSecret string) {

}

func GaeLoginWithRefreshToken(scope string, clientID string, clientSecret string) {

}

// AuthenticateComplete is a callback function to communicate both from the server to here, and from here to the caller.
type AuthenticateComplete func(output *Output)

// Data shared by both the authentication process and code/token exchange process.
// type gaeLogin struct {
// 	scope        string
// 	clientID     string
// 	clientSecret string
// 	port         int
// }

// // constructor to create ageLogin object.
// func newGaeLogin(scope string, clientID string, clientSecret string) *gaeLogin {
// 	login := new(gaeLogin)
// 	login.scope = scope
// 	login.clientID = clientID
// 	login.clientSecret = clientSecret
// 	return login
// }

// RequestAuthentication is the main method which begins this first part of the authentation process.
func RequestAuthentication(parentWindow *gtk.Window, scope string, clientID string, port int, processCallback AuthenticateComplete) error {

	var err error
	input := newInput(scope, clientID)
	var authWindow *gtk.Window
	// create window for browser
	authWindow, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		authWindow.Destroy() // which will trigger win destroy event
		return errors.WrapPrefix(err, "Unable to create authenticate window", 0)
	}

	server := NewAuthServer(input, func(output *Output) {
		// this is the server's thread, so ask the main thread for this
		glib.IdleAdd(func() bool {
			authWindow.Destroy() // which will trigger win destroy event
			// authCompleteCallback calls the processCallback
			processCallback(output)
			return false // only have IdleAdd() call this once
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
		server.BlockingClose()
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
	// authURL := new(url.URL)
	// authURL.Scheme = "https"
	// authURL.Host = "accounts.google.com"
	// authURL.Path = "/o/oauth2/v2/data"
	// authURLParams := url.Values{}
	// authURLParams.Set("scope", scope)
	// authURLParams.Set("redirect_uri", fmt.Sprintf("http://localhost:%d", port))
	// //"http://localhost")
	// authURLParams.Set("response_type", "code")
	// authURLParams.Set("client_id", clientID)
	// authURLParams.Set("state", state)
	// authURL.RawQuery = authURLParams.Encode()
	url := server.FormAuthURL()
	webView.LoadURI(authURL.String()) // blocks until it loads - requires UI
	return nil                        // no error
}

// RequestAccessToken takes code obtained from previous step and converts it into a token.
func RequestAccessToken(authCode string, clientID string, clientSecret string, port int) (string, error) /* *AccessToken */ {
	/*
		server := NewTokenServer(port, state, authCode, clientID, clientSecret)
		go func() {
			server.ListenAndServe()
		}()

		authWindow, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
		if err != nil {
			authWindow.Destroy() // which will trigger win destroy event
			return "", errors.WrapPrefix(err, "Unable to create authenticate window", 0)
		}
		authWindow.SetDefaultSize(850, 600)
		authWindow.SetTitle("Authenticate")
		// make authwindow a modal child window of the main window
		authWindow.SetModal(true)
		// authWindow.SetTransientFor(parentWindow)
		// close the data window if the parent window is closing
		authWindow.SetDestroyWithParent(true)
		// add event handler for when the data window is closing
		authWindow.Connect("destroy", func() bool {
			server.BlockingClose()
			return false // let the window close
		})
		// add box for layout
		vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
		if err != nil {
			authWindow.Destroy() // which will trigger win destroy event
			return "", errors.WrapPrefix(err, "Unable to create vertical box", 0)
		}
		// create the webview browser
		webView := webkit2.NewWebView()
		cancelButton, err := gtk.ButtonNewWithLabel("Cancel")
		if err != nil {
			authWindow.Destroy() // which will trigger win destroy event
			return "", errors.WrapPrefix(err, "Unable to create cancel button", 0)
		}
		cancelButton.Connect("clicked", func() {
			authWindow.Destroy() // which will trigger win destroy event
		})
		runButton, err := gtk.ButtonNewWithLabel("Run")
		if err != nil {
			authWindow.Destroy() // which will trigger win destroy event
			return "", errors.WrapPrefix(err, "Unable to create run button", 0)
		}
		runButton.Connect("clicked", func() {
			fmt.Printf("JS\n")
			webView.RunJavaScript("document.getElementById(\"token_form\").submit()", func(val *gojs.Value, err error) {
				if err != nil {
					fmt.Printf("JavaScript error  %s\n", err)
				} else {
					fmt.Printf("Hostname (from JavaScript): %q\n", val)
				}
				// java script callback
				// return
			})
		})
		webView.Connect("load-failed", func() {
			authWindow.Destroy() // which will trigger win destroy event
			fmt.Printf("Unable to load authentication page\n")
		})

		webView.Connect("load-changed", func(_ *glib.Object, i int) {
			loadEvent := webkit2.LoadEvent(i)
			switch loadEvent {
			case webkit2.LoadFinished:
				childList := authWindow.GetChildren()
				if childList.Length() == 0 {
					vbox.Add(webView)
					vbox.Add(cancelButton)
					vbox.Add(runButton)
					authWindow.Add(vbox)
					webView.Show()
					cancelButton.Show()
					runButton.Show()
					vbox.Show()
					authWindow.Show()
				} else {

				}

			}
		})
		fmt.Printf("load\n")
		webView.LoadURI(fmt.Sprintf("http://localhost:%d/load", port))
	*/
	params := url.Values{}
	// port := 7777 // freeport.GetPort()
	state := RandomDataBase64url(32)
	params.Set("code", authCode)
	params.Set("client_id", clientID)
	params.Set("client_secret", clientSecret)
	// The redirect_uri comes from the "Download JSON" button in the edit client_id screen in the API Manager.
	// Apparently for type "other" it can't be edited - you are just assigned this a a usable value.
	params.Set("redirect_uri", fmt.Sprintf("http://localhost:%d", port))
	// "http://localhost")

	params.Set("grant_type", "authorization_code")
	params.Set("state", state)
	resp, err := http.PostForm("https://www.googleapis.com/oauth2/v4/token", params)
	if err != nil {
		return "", errors.WrapPrefix(err, "Unable to convert code into token", 0)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.WrapPrefix(err, "Unable to read get token response", 0)
	}
	fmt.Printf(string(body[:]))
	return string(body[:]), nil
	// Need to parse the JSON response and look for {"error": "something", "error_description": "something has detail"}
	// return "", nil
}
