package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func NewClient(ctx context.Context, accesstoken string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accesstoken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
