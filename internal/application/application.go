package application

import (
	"context"
	config "github.com/golang-friends/members/internal/config"
	"github.com/google/go-github/v32/github"
	"log"
	"net/http"
	"sync"
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
	cli := github.NewClient(app.Client)
	var wg sync.WaitGroup

	for _, admin := range app.Config.Admins {
		if dryRun {
			log.Printf("Adding %v as admin", admin)
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("Adding %v as admin", admin)
			_, _, err := cli.Organizations.EditOrgMembership(
				context.Background(), admin, app.Config.Orgname, &github.Membership{
					Role: github.String("admin"),
				})
			if err != nil {
				log.Printf("EditOrgMembership (admin) has failed for %v: %v", admin, err)
			}
		}()
	}

	for _, member := range app.Config.Members {
		if dryRun {
			log.Printf("Adding %v as member", member)
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("Adding %v as member", member)
			_, _, err := cli.Organizations.EditOrgMembership(
				context.Background(), member, app.Config.Orgname, &github.Membership{
					Role: github.String("member"),
				})
			if err != nil {
				log.Printf("EditOrgMembership (member) has failed for %v: %v", member, err)
			}
		}()
	}

	return nil
}

func NewApplication(cfg *config.Config, client *http.Client) *Application {
	return &Application{Config: cfg, Client: client}
}
