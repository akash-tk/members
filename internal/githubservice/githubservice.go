package githubservice

import (
	"context"
	"net/http"

	"github.com/golang-friends/members/internal/application"

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
	for _, member := range members {
		log.WithField("member", member).Info("removing a member")
		_, err := g.client.Organizations.RemoveMember(context.Background(), g.orgName, member)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddAdmins adds members as "admin" role.
func (g GitHubService) AddAdmins(admins []string) error {
	for _, admin := range admins {
		log.WithField("admin", admin).Info("adding an admin")
		err := g.editOrgMembership(admin, enums.RoleAdmin)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddMembers adds members as "member" role.
func (g GitHubService) AddMembers(members []string) error {
	for _, member := range members {
		log.WithField("member", member).Info("adding a member")
		err := g.editOrgMembership(member, enums.RoleMember)
		if err != nil {
			return err
		}
	}
	return nil
}

// editOrgMembership is the internal function used inviting members.
func (g GitHubService) editOrgMembership(member string, role enums.Role) error {
	_, _, err := g.client.Organizations.EditOrgMembership(
		context.Background(), member, g.orgName, &github.Membership{
			Role: github.String(role.String()),
		})
	return err
}

// New is the factory for GitHubService.
func New(orgName string, httpClient *http.Client) *GitHubService {
	client := github.NewClient(httpClient)
	return &GitHubService{
		client:  client,
		orgName: orgName,
	}
}
