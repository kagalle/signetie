package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/gotk3/gotk3/gtk"
	"github.com/kagalle/signetie/client_golang/gae"
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
		gtk.MainQuit()
	})

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	if err != nil {
		log.Fatal("Unable to create vertical box:", err)
	}
	win.Add(vbox)

	auth := gae.NewAuthenticate(win, func(auth *gae.Authenticate) bool {
		continueOn := false
		if auth.Found() {
			fmt.Printf("Code obtained %s\n", auth.Code())
			continueOn = true
		} else if auth.Cancelled() {
			fmt.Printf("User canceled authenticate\n")
		} else {
			fmt.Printf("Unexpected authenticate condition\n")
		}
		return continueOn
	})

	runButton, err := gtk.ButtonNewWithLabel("Run")
	if err != nil {
		log.Fatal("Unable to create run button:", err)
	}
	runButton.Connect("clicked", func() {
		auth.Run("https://www.googleapis.com/auth/userinfo.profile",
			"192820621204-nrkum19gt8a7hjrrkrdpdhh2qgmi0toq.apps.googleusercontent.com")
	})
	vbox.Add(runButton)
	win.ShowAll()
	gtk.Main()
}
