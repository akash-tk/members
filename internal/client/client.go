package client

import (
	"context"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
)

// NewOAuthClient returns httpClient with the GitHub OAuth access token.
// Without the token, it's not able to view private members nor to invite users.
// Hence, this is required. Otherwise, it will panic.
func NewOAuthClient(accessToken string) *http.Client {
	if accessToken == "" {
		panic(errors.New("GitHub access token was not provided"))
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	return oauth2.NewClient(ctx, ts)
}
