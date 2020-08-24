package cmd

import (
	"github.com/golang-friends/members/internal/application"
	"github.com/golang-friends/members/internal/client"
	"github.com/golang-friends/members/internal/config"
	"github.com/spf13/cobra"
)

var dryRun bool

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Make REST API calls based on members.yaml",
	Long:    "",
	Example: "members update --dry-run",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.FromViper()
		app := application.NewApplication(cfg, client.NewOAuthClient(gitHubOAuthToken))
		return app.Update(dryRun)
	},
}

func init() {
	updateCmd.Flags().BoolVarP(&dryRun,
		"dry-run",
		"d",
		false,
		"dry run will not make REQUESTS")
}
