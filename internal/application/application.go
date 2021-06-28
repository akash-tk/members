package application

import (
	"sort"

	"github.com/golang-friends/members/internal/config"
	"github.com/golang-friends/members/internal/enums"
	"github.com/google/go-github/v35/github"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
)

// GitHubService is a dependency used by Application.
// It wraps the implementation so that it's easier to test.
type GitHubService interface {
	RemoveMembers(members []string) error
	AddAdmins(admins []string) error
	AddMembers(members []string) error
	ListMembersByRole(role enums.Role) ([]*github.User, error)
}

// Application is the main entry struct for the program.
type Application struct {
	cfg           *config.Config
	gitHubService GitHubService
}

// GetConfigFromGitHub fetches members and admins based on the GitHub server information.
func (app *Application) GetConfigFromGitHub() config.Config {
	admins, err := app.gitHubService.ListMembersByRole(enums.RoleAdmin)
	if err != nil {
		panic(err)
	}

	members, err := app.gitHubService.ListMembersByRole(enums.RoleMember)
	if err != nil {
		panic(err)
	}

	var cfg config.Config
	cfg.Orgname = app.cfg.Orgname

	for _, admin := range admins {
		cfg.Admins = append(cfg.Admins, admin.GetLogin())
	}

	for _, member := range members {
		cfg.Members = append(cfg.Members, member.GetLogin())
	}

	sort.Strings(cfg.Admins)
	sort.Strings(cfg.Members)

	return cfg
}

// Update will actually invite users unless dryRun is true.
// In case of dryRun, it will just print without calling the server.
func (app *Application) Update(dryRun bool) error {
	oldConfig := app.GetConfigFromGitHub()
	return app.UpdateV2(oldConfig, dryRun)
}

// UpdateV2 is the new version that supports remove feature.
func (app Application) UpdateV2(oldConfig config.Config, dryRun bool) error {
	oldAllMembers :=
		sets.NewString(append(oldConfig.Members, oldConfig.Admins...)...)
	newAllMembers :=
		sets.NewString(append(app.cfg.Members, app.cfg.Admins...)...)

	shouldRemoveMembers := oldAllMembers.Difference(newAllMembers)

	if len(shouldRemoveMembers) > 0 {
		if dryRun {
			log.WithFields(log.Fields{
				"TO_BE_REMOVED_MEMBERS": shouldRemoveMembers,
			}).Info("removing members")
		} else {
			err := app.gitHubService.RemoveMembers(shouldRemoveMembers.UnsortedList())
			if err != nil {
				return err
			}
		}
	}

	if len(app.cfg.Admins) > 0 {
		if dryRun {
			log.WithFields(log.Fields{"TO_BE_ADMINS": app.cfg.Admins}).Info("adding admins")
		} else {
			err := app.gitHubService.AddAdmins(app.cfg.Admins)
			if err != nil {
				return err
			}
		}
	}

	if len(app.cfg.Members) > 0 {
		if dryRun {
			log.WithFields(log.Fields{
				"TO_BE_MEMBERS": app.cfg.Members,
			}).Info("adding members")
		} else {
			err := app.gitHubService.AddMembers(app.cfg.Members)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// NewApplication is a factory function that returns Application with the dependency.
func NewApplication(cfg *config.Config, githubService GitHubService) *Application {
	return &Application{cfg, githubService}
}
