package config

import (
	"github.com/spf13/viper"
)

// Config matches the schema of members.yaml.
type Config struct {
	Orgname string
	Admins  []string
	Members []string
}

// FromViper parses members.yaml via Viper.
func FromViper() *Config {
	return &Config{
		Orgname: viper.GetString("org"),
		Admins:  viper.GetStringSlice("admins"),
		Members: viper.GetStringSlice("members"),
	}
}
