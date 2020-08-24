package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Orgname string
	Admins  []string
	Members []string
}

func FromViper() *Config {
	return &Config{
		Orgname: viper.GetString("org"),
		Admins:  viper.GetStringSlice("admins"),
		Members: viper.GetStringSlice("members"),
	}
}
