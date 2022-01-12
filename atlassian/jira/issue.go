package jira

import (
	"errors"
	"fmt"
	"gitacthelp/atlassian"
	"github.com/andygrunwald/go-jira"
	"time"
)

type Issue struct {
	Key        string
	Title      string
	Type       string
	Status     string
	FixVersion *string
	ResolvedAt *time.Time
}

func (issue Issue) GetUrl(config atlassian.AtlassianConfig) string {
	return fmt.Sprintf("<%s/browse/%s|%s>", config.BaseUrl, issue.Key, issue.Key)
}

func (issue Issue) GetFixVersion() string {
	if issue.FixVersion == nil {
		return "nil"
	}
	return *issue.FixVersion
}

func (clnt JiraClient) GetIssue(issueId string) (task *Issue, err error) {
	opt := jira.GetQueryOptions{Expand: "changelog"}
	rawIssue, _, err := clnt.client.Issue.Get(issueId, &opt)
	if err != nil {
		return nil, err
	}
	var fixVer *string = nil
	versions := rawIssue.Fields.FixVersions
	if len(versions) > 0 {
		fixVer = &versions[0].Name
	}
	at, err := findResolvedAt(*rawIssue)
	if err != nil {
		return nil, err
	}
	return &Issue{
		Key:        rawIssue.Key,
		Status:     rawIssue.Fields.Status.Name,
		FixVersion: fixVer,
		ResolvedAt: at,
	}, nil
}

func findResolvedAt(issue jira.Issue) (*time.Time, error) {
	if issue.Fields.Status.Name != "Resolved" {
		return nil, nil
	}
	changeLog := issue.Changelog.Histories
	for _, changeLogItem := range changeLog {
		items := changeLogItem.Items
		for _, item := range items {
			if item.ToString == "Resolved" {
				createdTime, e := changeLogItem.CreatedTime()
				return &createdTime, e
			}
		}
	}
	return nil, errors.New("failed to find resolved at")
}
