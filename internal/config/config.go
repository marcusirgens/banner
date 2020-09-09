package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	GithubAccessToken string
}

var DefaultConfig Config

func init() {
	viper.SetConfigName(".banner")
	viper.AddConfigPath("$HOME")
	viper.SetConfigType("env")
	viper.SetEnvPrefix("BANNER")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		nfe := &viper.ConfigFileNotFoundError{}
		if errors.As(err, nfe) {
			// nop
		} else {
			fmt.Println("Error reading config file")
		}

	}

	setDefaults()
}

func setDefaults() {
	DefaultConfig.GithubAccessToken = viper.GetString("github_access_token")
}
