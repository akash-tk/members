package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "members",
	Short: "Manages GitHub org members",
}
var gitHubOAuthToken string

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	const flagName = "github-oauth-token"
	rootCmd.PersistentFlags().StringVarP(
		&gitHubOAuthToken,
		flagName,
		"t",
		"",
		"GitHub OAuth Token (required)")

	viper.SetConfigName("members")
	viper.SetConfigType("yaml")
	// PWD
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed while reading config via viper: %v", err))
	}

	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(writeCmd)
}
