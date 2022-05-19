package config

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	GithubToken string
}

var Config config

func MustLoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalln("Fatal not found config file:", err)
		} else {
			log.Fatalln("Fatal error config file:", err)
		}
	}

	viper.Unmarshal(&Config)
}
