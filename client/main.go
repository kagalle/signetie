package main

import (
	"errors"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gotk3/gotk3/gtk"
	"github.com/kagalle/signetie/client_golang/config"
	"github.com/kagalle/signetie/client_golang/gae/gaeAccessToken"
	"github.com/kagalle/signetie/client_golang/gae/login"
	"github.com/spf13/viper"
)

func init() {
	// Ref Re: logging: https://dave.cheney.net/2015/11/05/lets-talk-about-logging
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(new(logrus.JSONFormatter))
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
	// initialize configuration
	config.InitConfiguration()
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
		/*        "gae_login_scope":   "https://www.googleapis.com/auth/userinfo.profile",
		"gae_client_id":     "xxxx.apps.googleusercontent.com",
		"gae_client_secret": "xxxx",
		"gae_test_refresh_token":  "xxxx",
		"location_timezone": "America/New_York"
		*/
		gaeLogin := login.NewGaeLogin(win, viper.GetString("gae_login_scope"),
			viper.GetString("gae_client_id"),
			viper.GetString("gae_client_secret"))
		// tokenSet := gaeLogin.Login(nil)
		tokenSet := new(gaeAccessToken.TokenSet)
		tokenSet.RefreshToken = viper.GetString("gae_test_refresh_token")
		location, _ := time.LoadLocation(viper.GetString("location_timezone"))
		tokenSet.ExpiresOn = time.Date(2017, time.March, 12, 12, 20, 9, 75240, location)
		tokenSet.Log("Test tokenset with manual refresh token")
		tokenSet = gaeLogin.Login(tokenSet)
		if tokenSet != nil {
			tokenSet.Log("TokenSet returned from refresh tokenset")
		}
	})
	vbox.Add(runButton)
	win.ShowAll()
	gtk.Main()
}
