package jira

import (
	"errors"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/devindi/gitacthelp/atlassian"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
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

func (clnt JiraClient) CreateReleaseTask(summary string, description string, reporterId string) (*string, error) {
	meta, _, err := clnt.client.Issue.GetCreateMeta(clnt.projectKey)
	if err != nil {
		log.Errorln("Failed to get jira metadata")
		return nil, err
	}

	var releaseIssueType *jira.MetaIssueType
	for _, issueType := range meta.Projects[0].IssueTypes {
		if issueType.Name == "Release" {
			releaseIssueType = issueType
		}
	}

	issue, resp, err := clnt.client.Issue.Create(&jira.Issue{
		Fields: &jira.IssueFields{
			Project: jira.Project{
				Key: clnt.projectKey,
			},
			Type: jira.IssueType{
				ID: releaseIssueType.Id,
			},
			Summary:     summary,
			Description: description,
			Reporter: &jira.User{
				AccountID: reporterId,
			},
		},
	})
	if err != nil {
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Errorln("Failed to create issue. Response: ", bodyString)
		return nil, err
	}

	return &issue.Key, nil
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
