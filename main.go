package main

import (
	"errors"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/gotk3/gotk3/gtk"
	"github.com/kagalle/signetie/client_golang/gae/login"
)

func init() {
	// Ref Re: logging: https://dave.cheney.net/2015/11/05/lets-talk-about-logging
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	// http://stackoverflow.com/a/19934989
	defer func() {
		var err error
		r := recover()
		if r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			logrus.WithError(err).Error("Panic")
		}
	}()
	// Initialize GTK without parsing any command line arguments.
	runtime.LockOSThread()
	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		logrus.WithError(err).Error("Unable to create window")
	}
	win.SetDefaultSize(850, 600)
	win.SetTitle("Signetie")
	win.Connect("destroy", func() {
		logrus.Debugln("Quitting")
		gtk.MainQuit()
	})

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	if err != nil {
		logrus.WithError(err).Error("Unable to create vertical box")
	}
	win.Add(vbox)
	runButton, err := gtk.ButtonNewWithLabel("Run")
	if err != nil {
		logrus.WithError(err).Error("Unable to create run button")
	}
	runButton.Connect("clicked", func() {
		gaeLogin := login.NewGaeLogin(win,
			"https://www.googleapis.com/auth/userinfo.profile",                         // scope
			"192820621204-nrkum19gt8a7hjrrkrdpdhh2qgmi0toq.apps.googleusercontent.com", // clientID
			"Tx3wbyqLBjDFOH7l-ZXr7-Ot")                                                 // client secret
		tokenSet := gaeLogin.Login(nil)
		if tokenSet != nil {
			tokenSet.Log()
		}
	})
	vbox.Add(runButton)
	win.ShowAll()
	gtk.Main()
}
