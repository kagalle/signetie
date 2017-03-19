package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Settings is a struct to represent what is stored as settings.
// https://github.com/spf13/viper
// http://stackoverflow.com/a/13004756
// http://localhost:6060/pkg/os/user/#Current

/*
   "gae_login_scope":   "https://www.googleapis.com/auth/userinfo.profile",
   "gae_client_id":     "xxxx.apps.googleusercontent.com",
   "gae_client_secret": "xxxx",
   "gae_test_refresh_token":  "xxxx",
   "location_timezone": "America/New_York"
*/

func InitConfiguration() {
	viper.SetConfigName("config")          // name of config file (without extension)
	viper.AddConfigPath("/etc/signetie/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.signetie") // call multiple times to add many search paths
	err := viper.ReadInConfig()            // Find and read the config file
	if err != nil {                        // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
