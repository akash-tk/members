package application

import (
	"context"
	config "github.com/golang-friends/members/internal/config"
	"github.com/google/go-github/v32/github"
	"log"
	"net/http"
)

type Application struct {
	Config *config.Config
	Client *http.Client
}

func (app *Application) GetConfigFromGitHub() config.Config {
	client := github.NewClient(app.Client)

	admins, _, err := client.Organizations.ListMembers(context.Background(), app.Config.Orgname, &github.ListMembersOptions{
		PublicOnly: false,
		Role:       "admin",
	})

	if err != nil {
		panic(err)
	}

	members, _, err := client.Organizations.ListMembers(context.Background(), app.Config.Orgname, &github.ListMembersOptions{
		PublicOnly: false,
		Role:       "member",
	})

	if err != nil {
		panic(err)
	}

	var config config.Config

	for _, admin := range admins {
		config.Admins = append(config.Admins, admin.GetLogin())
	}

	for _, member := range members {
		config.Members = append(config.Members, member.GetLogin())
	}

	return config
}

func (app *Application) Update(dryRun bool) error {
	for _, admin := range app.Config.Admins {
		if dryRun {
			log.Printf("Adding %v as admin", admin)
			continue
		}
	}

	for _, member := range app.Config.Members {
		if dryRun {
			log.Printf("Adding %v as member", member)
			continue
		}
	}
	return nil
}

func NewApplication(cfg *config.Config, client *http.Client) *Application {
	return &Application{Config: cfg, Client: client}
}
