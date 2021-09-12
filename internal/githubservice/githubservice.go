package githubservice

import (
	"context"
	"github.com/golang-friends/members/internal/application"
	"github.com/golang-friends/members/internal/config"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/golang-friends/members/internal/enums"
	"github.com/google/go-github/v35/github"
	log "github.com/sirupsen/logrus"
)

// GitHubService is the actual implementation of application.GitHubService interface.
type GitHubService struct {
	client  *github.Client
	orgName string
}

// This ensures the GitHubService implements application.GitHubService.
var _ application.GitHubService = (*GitHubService)(nil)

// ListMembersByRole returns GitHub users based on role (admin/member).
func (g GitHubService) ListMembersByRole(role enums.Role) ([]*github.User, error) {
	members, _, err := g.client.Organizations.ListMembers(context.Background(), g.orgName, &github.ListMembersOptions{
		PublicOnly: false,
		Role:       role.String(),
	})
	return members, err
}

// RemoveMembers removes members from orgName.
func (g GitHubService) RemoveMembers(members []string) error {
	var errMerged error
	for _, member := range members {
		log.WithField("member", member).Info("removing a member")
		_, err := g.client.Organizations.RemoveMember(context.Background(), g.orgName, member)
		if err != nil {
			if errMerged == nil {
				errMerged = err
			} else {
				errMerged = errors.Wrap(err, errMerged.Error())
			}
		}
	}
	return errMerged
}

// AddAdmins adds members as "admin" role.
func (g GitHubService) AddAdmins(admins []string) error {
	var errMerged error
	for _, admin := range admins {
		log.WithField("admin", admin).Info("adding an admin")
		err := g.editOrgMembership(admin, enums.RoleAdmin)
		if err != nil {
			if errMerged == nil {
				errMerged = err
			} else {
				errMerged = errors.Wrap(err, errMerged.Error())
			}
		}
	}
	return errMerged
}

// AddMembers adds members as "member" role.
func (g GitHubService) AddMembers(members []string) error {
	var errMerged error
	for _, member := range members {
		log.WithField("member", member).Info("adding a member")
		err := g.editOrgMembership(member, enums.RoleMember)
		if err != nil {
			if errMerged == nil {
				errMerged = err
			} else {
				errMerged = errors.Wrap(err, errMerged.Error())
			}
		}
	}
	return errMerged
}

// editOrgMembership is the internal function used inviting members.
func (g GitHubService) editOrgMembership(member string, role enums.Role) error {
	_, _, err := g.client.Organizations.EditOrgMembership(
		context.Background(), member, g.orgName, &github.Membership{
			Role: github.String(role.String()),
		})
	return err
}

type GitHubOAuthToken string

// New is the factory for GitHubService.
func New(cfg *config.Config, githubOAuthToken GitHubOAuthToken) (*GitHubService, error) {
	if githubOAuthToken == "" {
		return nil, errors.New("GitHub OAuth Token (passed as -t) is needed")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(githubOAuthToken)},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &GitHubService{
		client:  client,
		orgName: cfg.Orgname,
	}, nil
}
