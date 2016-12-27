package main

import "fmt"
import "github.com/kagalle/signetie/client_golang/gae"

func main() {
	fmt.Printf("A")
	gae.Authenicate("https://www.googleapis.com/auth/userinfo.profile",
		"192820621204-nrkum19gt8a7hjrrkrdpdhh2qgmi0toq.apps.googleusercontent.com")
	fmt.Printf("Z")
}
