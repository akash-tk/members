package cmd

import (
	"github.com/golang-friends/members/internal"
	"github.com/golang-friends/members/internal/githubservice"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "it will write `members.yaml` by fetching the server information",
	Long:  "It fetches members from the server and overwrite `members.yaml`",
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := internal.ProvideApplication(githubservice.GitHubOAuthToken(gitHubOAuthToken))
		if err != nil {
			return err
		}
		newConfig := app.GetConfigFromGitHub()
		viper.Set("admins", newConfig.Admins)
		viper.Set("members", newConfig.Members)
		return viper.WriteConfig()
	},
}
