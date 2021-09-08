//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/golang-friends/members/internal/application"
	"github.com/golang-friends/members/internal/config"
	"github.com/golang-friends/members/internal/githubservice"
	"github.com/google/wire"
	"net/http"
)

func ProvideApplication() *application.Application {
	panic(wire.Build(
		application.NewApplication,
		config.FromViper,
		githubservice.New,
		wire.Bind(new(application.GitHubService), new(*githubservice.GitHubService)),
		wire.Value(http.DefaultClient),
		wire.FieldsOf(new(*config.Config), "Orgname"),
	))
}
