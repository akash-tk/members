package application

import (
	"context"
	"log"
	"net/http"
	"sort"
	"sync"

	config "github.com/golang-friends/members/internal/config"
	"github.com/google/go-github/v32/github"
)

// Application is the main entry struct for the program.
type Application struct {
	Config *config.Config
	Client *http.Client
}

// GetConfigFromGitHub fetches members and admins based on the GitHub server information.
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
	config.Orgname = app.Config.Orgname

	for _, admin := range admins {
		config.Admins = append(config.Admins, admin.GetLogin())
	}

	for _, member := range members {
		config.Members = append(config.Members, member.GetLogin())
	}

	sort.Strings(config.Admins)
	sort.Strings(config.Members)

	return config
}

// Update will actually invite users unless dryRun is true.
// In case of dryRun, it will just print without calling the server.
func (app *Application) Update(dryRun bool) error {
	cli := github.NewClient(app.Client)
	var wg sync.WaitGroup

	for _, admin := range app.Config.Admins {
		if dryRun {
			log.Printf("Adding %v as admin", admin)
			continue
		}

		wg.Add(1)
		// needs to capture
		admin := admin
		go func() {
			defer wg.Done()
			editMembership(cli, app.Config.Orgname, admin, "admin")
		}()
	}

	for _, member := range app.Config.Members {
		if dryRun {
			log.Printf("Adding %v as member", member)
			continue
		}

		wg.Add(1)
		// needs to capture
		member := member
		go func() {
			defer wg.Done()
			editMembership(cli, app.Config.Orgname, member, "member")
		}()
	}

	wg.Wait()
	return nil
}

func editMembership(cli *github.Client, orgname, username, role string) {
	log.Printf("Adding %v as %v", username, role)
	_, _, err := cli.Organizations.EditOrgMembership(
		context.Background(), username, orgname, &github.Membership{
			Role: github.String(role),
		})
	if err != nil {
		log.Printf("EditOrgMembership (%v) has failed for %v: %v", role, username, err)
	}
}

// NewApplication is a factory function that returns Application with the dependency.
func NewApplication(cfg *config.Config, client *http.Client) *Application {
	return &Application{Config: cfg, Client: client}
}
