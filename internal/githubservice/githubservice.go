package githubservice

import (
	"context"
	"net/http"

	"github.com/golang-friends/members/internal/application"
	"github.com/golang-friends/members/internal/enums"
	"github.com/google/go-github/v35/github"
	log "github.com/sirupsen/logrus"
)

type GitHubService struct {
	client  *github.Client
	orgName string
}

func (g GitHubService) ListMembersByRole(role enums.Role) ([]*github.User, error) {
	members, _, err := g.client.Organizations.ListMembers(context.Background(), g.orgName, &github.ListMembersOptions{
		PublicOnly: false,
		Role:       role.String(),
	})
	return members, err
}

var _ application.GitHubService = (*GitHubService)(nil)

func (g GitHubService) RemoveMembers(members []string) error {
	for _, member := range members {
		log.WithField("member", member).Info("removing a member")
		_, err := g.client.Organizations.RemoveMember(context.Background(), g.orgName, member)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g GitHubService) AddAdmins(admins []string) error {
	for _, admin := range admins {
		log.WithField("admin", admin).Info("adding an admin")
		err := g.editOrgMembership(admin, enums.RoleAdmin.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func (g GitHubService) AddMembers(members []string) error {
	for _, member := range members {
		log.WithField("member", member).Info("adding a member")
		err := g.editOrgMembership(member, enums.RoleMember.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func (g GitHubService) editOrgMembership(member, role string) error {
	_, _, err := g.client.Organizations.EditOrgMembership(
		context.Background(), member, g.orgName, &github.Membership{
			Role: github.String(role),
		})
	return err
}

func New(orgName string, httpClient *http.Client) *GitHubService {
	client := github.NewClient(httpClient)
	return &GitHubService{
		client:  client,
		orgName: orgName,
	}
}
