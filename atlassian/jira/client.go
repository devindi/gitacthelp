package jira

import (
	"github.com/andygrunwald/go-jira"
	"github.com/devindi/gitacthelp/atlassian"
)

type JiraClient struct {
	client     *jira.Client
	projectKey string
}

func NewJiraClient(config atlassian.AtlassianConfig) (*JiraClient, error) {
	tp := jira.BasicAuthTransport{
		Username: config.Username,
		Password: config.Token,
	}
	client, err := jira.NewClient(tp.Client(), config.BaseUrl)
	if err != nil {
		return nil, err
	}
	return &JiraClient{client: client, projectKey: config.ProjectKey}, nil
}
