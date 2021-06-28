package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func init() {
	viper.SetConfigName("members")
	viper.SetConfigType("yaml")
	// internal/config/../.. => project root.
	viper.AddConfigPath("../..")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed while reading config via viper: %v", err))
	}
}

func TestFromViper(t *testing.T) {
	cfg := FromViper()
	assert.True(t, sort.StringsAreSorted(cfg.Admins))
	assert.True(t, sort.StringsAreSorted(cfg.Members))
}
