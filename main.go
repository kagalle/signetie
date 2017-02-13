package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/gotk3/gotk3/gtk"
	"github.com/kagalle/signetie/client_golang/gae/login"
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
		gaeLogin := login.NewGaeLogin(win,
			"https://www.googleapis.com/auth/userinfo.profile",                         // scope
			"192820621204-nrkum19gt8a7hjrrkrdpdhh2qgmi0toq.apps.googleusercontent.com", // clientID
			"Tx3wbyqLBjDFOH7l-ZXr7-Ot")                                                 // client secret
		tokenSet := gaeLogin.Login(nil)
		if tokenSet != nil {
			tokenSet.Print()
		}
	})
	vbox.Add(runButton)
	win.ShowAll()
	gtk.Main()
}
