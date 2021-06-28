package application

import (
	"github.com/golang-friends/members/internal/application/mock_application"
	"github.com/golang-friends/members/internal/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
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
