package cmd

import (
	"github.com/golang-friends/members/internal/application"
	"github.com/golang-friends/members/internal/client"
	"github.com/golang-friends/members/internal/config"
	"github.com/golang-friends/members/internal/githubservice"
	"github.com/spf13/cobra"
)

var dryRun bool

// updateCmd will remove/invite members based on members.yaml
var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Make REST API calls based on members.yaml",
	Long:    "",
	Example: "members update --dry-run",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.FromViper()
		githubService := githubservice.New(cfg.Orgname, client.NewOAuthClient(gitHubOAuthToken))
		app := application.NewApplication(cfg, githubService)
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
