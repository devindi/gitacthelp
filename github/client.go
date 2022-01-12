package github

import (
	"context"
	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	impl       *github.Client
	owner      string
	repository string
	context    context.Context
}

func NewGithubClient(config GithubConfig) *GithubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return &GithubClient{impl: github.NewClient(tc), owner: config.Owner, repository: config.Repository, context: context.Background()}
}
