package cmd

import (
	"github.com/golang-friends/members/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "it will write `members.yaml` by fetching the server information",
	Long:  "It fetches members from the server and overwrite `members.yaml`",
	RunE: func(cmd *cobra.Command, args []string) error {
		newConfig := internal.ProvideApplication().GetConfigFromGitHub()
		viper.Set("admins", newConfig.Admins)
		viper.Set("members", newConfig.Members)
		return viper.WriteConfig()
	},
}
