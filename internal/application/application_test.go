package application

import (
	"testing"

	"github.com/golang-friends/members/internal/application/mock_application"
	"github.com/golang-friends/members/internal/config"
	"github.com/golang-friends/members/internal/enums"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v35/github"
	"github.com/stretchr/testify/assert"
)

func TestApplication_UpdateV2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	{
		mockGitHubService := mock_application.NewMockGitHubService(ctrl)

		// GIVEN the old config from GitHub shows one admin.
		oldConfig := config.Config{
			Orgname: "orgname",
			Admins:  []string{"kkweon"},
			Members: nil,
		}

		// WHEN the new config without any members.
		app := Application{
			cfg: &config.Config{
				Orgname: "orgname",
				Admins:  nil,
				Members: nil,
			},
			gitHubService: mockGitHubService,
		}

		// THEN it should remove members in the old config.
		mockGitHubService.
			EXPECT().
			RemoveMembers(gomock.Eq([]string{"kkweon"})).
			Return(nil)

		err := app.UpdateV2(oldConfig, false)
		assert.NoError(t, err)
	}
}

func TestApplication_GetConfigFromGitHub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGitHubService := mock_application.NewMockGitHubService(ctrl)

	// GIVEN some members in remote.
	mockGitHubService.
		EXPECT().
		ListMembersByRole(gomock.Eq(enums.RoleAdmin)).
		Return([]*github.User{
			getGitHubUser("admin1"),
			getGitHubUser("admin2"),
		}, nil)
	mockGitHubService.
		EXPECT().
		ListMembersByRole(gomock.Eq(enums.RoleMember)).
		Return([]*github.User{getGitHubUser("member1"), getGitHubUser("member2")}, nil)

	app := NewApplication(&config.Config{
		Orgname: "orgname",
		Admins:  nil,
		Members: nil,
	}, mockGitHubService)

	// WHEN cfg is requested from GitHub.
	cfg := app.GetConfigFromGitHub()

	// THEN it should return the config with members from GitHub.
	assert.Equal(t, config.Config{
		Orgname: "orgname",
		Admins:  []string{"admin1", "admin2"},
		Members: []string{"member1", "member2"},
	}, cfg)
}

func getGitHubUser(s string) *github.User {
	return &github.User{
		Login: &s,
	}
}
