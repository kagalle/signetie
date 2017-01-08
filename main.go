package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/gotk3/gotk3/gtk"
	"github.com/kagalle/signetie/client_golang/gae"
	"github.com/phayes/freeport"
)

func main() {
	// Initialize GTK without parsing any command line arguments.
	runtime.LockOSThread()
	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetDefaultSize(850, 600)
	win.SetTitle("Signetie")
	win.Connect("destroy", func() {
		fmt.Printf("Quitting\n")
		gtk.MainQuit()
	})

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	if err != nil {
		log.Fatal("Unable to create vertical box:", err)
	}
	win.Add(vbox)
	runButton, err := gtk.ButtonNewWithLabel("Run")
	if err != nil {
		log.Fatal("Unable to create run button:", err)
	}
	runButton.Connect("clicked", func() {
		port := freeport.GetPort() // 7777

		// RequestAuthentication(parentWindow, scope, client_id, callback)
		err := gae.RequestAuthentication(win, "https://www.googleapis.com/auth/userinfo.profile",
			"192820621204-nrkum19gt8a7hjrrkrdpdhh2qgmi0toq.apps.googleusercontent.com", port,
			afterRequestAuthentication)
		if err != nil {
			log.Fatal(err)
		}
	})
	vbox.Add(runButton)
	win.ShowAll()
	gtk.Main()
}

func afterRequestAuthentication(auth *gae.Authenticate, port int) {
	continueOn := false
	if auth.Found() {
		fmt.Printf("Code obtained %s\n", auth.Code())
		continueOn = true

		// RequestAccessToken(authCode string, clientID string, clientSecret string) (string, error)
		accessToken, err := gae.RequestAccessToken(auth.Code(),
			"192820621204-nrkum19gt8a7hjrrkrdpdhh2qgmi0toq.apps.googleusercontent.com",
			"Tx3wbyqLBjDFOH7l-ZXr7-Ot", port)
		if err != nil {
			log.Fatal("Unable to exchange code for token:", err)
		}
		fmt.Printf("token_response:%s", accessToken)

	} else if auth.Cancelled() {
		fmt.Printf("User canceled authenticate\n")
	} else {
		log.Fatal(auth.Error())
	}
	if continueOn {
		fmt.Printf("Continuing\n")
	} else {
		fmt.Printf("Not continuing\n")
	}
}
