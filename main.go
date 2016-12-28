package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/gotk3/gotk3/gtk"
	"github.com/kagalle/signetie/client_golang/gae"
)

func main() {
	fmt.Printf("A1")
	// Initialize GTK without parsing any command line arguments.
	runtime.LockOSThread()
	gtk.Init(nil)
	fmt.Printf("A2")
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	fmt.Printf("A3")
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Signetie")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	fmt.Printf("A4")

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	if err != nil {
		log.Fatal("Unable to create vertical box:", err)
	}
	win.Add(vbox)
	fmt.Printf("A5")

	auth := gae.NewAuthenticate(win, func(auth *gae.Authenticate) bool {
		continueOn := false
		if auth.GetFound() {
			fmt.Printf("Code obtained %s", auth.GetCode())
			continueOn = true
		} else if auth.GetCancelled() {
			fmt.Printf("User canceled authenticate")
		} else {
			fmt.Printf("Unexpected authenticate condition")
		}
		return continueOn
	})
	fmt.Printf("A6")

	runButton, err := gtk.ButtonNewWithLabel("Run")
	if err != nil {
		log.Fatal("Unable to create run button:", err)
	}
	auth.Setup()
	runButton.Connect("clicked", func() {
		fmt.Printf("M1")
		auth.Run("https://www.googleapis.com/auth/userinfo.profile",
			"192820621204-nrkum19gt8a7hjrrkrdpdhh2qgmi0toq.apps.googleusercontent.com")
		// glib.IdleAdd(auth.Run,
		// 	"https://www.googleapis.com/auth/userinfo.profile",
		// 	"192820621204-nrkum19gt8a7hjrrkrdpdhh2qgmi0toq.apps.googleusercontent.com")
		fmt.Printf("M2")
	})
	vbox.Add(runButton)
	fmt.Printf("A7")

	win.ShowAll()
	gtk.Main()

	fmt.Printf("Z")
}
