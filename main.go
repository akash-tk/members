//go:generate mockgen -source internal/application/application.go -destination internal/application/mock_application/mock_application.go GitHubService
package main

import (
	"github.com/golang-friends/members/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.WithError(err).Error("cmd.Execute() has failed")
	}
}
