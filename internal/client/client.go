package client

import (
	"context"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

func NewOAuthClient(accessToken string) *http.Client {
	if accessToken == "" {
		log.Fatal("GitHub access token was not provided")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	return oauth2.NewClient(ctx, ts)
}
