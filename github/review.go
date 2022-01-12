package github

import (
	"fmt"
	"github.com/google/go-github/v28/github"
	log "github.com/sirupsen/logrus"
	"time"
)

type GithubReview struct {
	ID                 int
	Title              string
	Author             string
	RequestedReviewers []string
	CreatedAt          *time.Time
}

type IssueTimelineItem struct {
	CreatedAt *string `json:"created_at"`
	Event     *string `json:"event"`
}

func (client GithubClient) GetReviews() ([]GithubReview, error) {
	pullRequestOpt := github.PullRequestListOptions{ListOptions: github.ListOptions{PerPage: 100}}
	pullRequests, _, err := client.impl.PullRequests.List(client.context, client.owner, client.repository, &pullRequestOpt)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var result []GithubReview
	for _, pullRequest := range pullRequests {
		var reviewers []string
		for _, reviewer := range pullRequest.RequestedReviewers {
			reviewers = append(reviewers, reviewer.GetLogin())
		}
		var timestamp = client.tryToFetchRequestDate(*pullRequest.Number)
		if timestamp == nil {
			timestamp = pullRequest.CreatedAt
		}
		result = append(result, GithubReview{
			ID:                 *pullRequest.Number,
			Title:              pullRequest.GetTitle(),
			Author:             pullRequest.GetUser().GetLogin(),
			RequestedReviewers: reviewers,
			CreatedAt:          timestamp,
		})
	}
	return result, nil
}

//Beware. We are using experimental API here
func (client GithubClient) tryToFetchRequestDate(issueNumber int) *time.Time {
	request, _ := client.impl.NewRequest("GET", fmt.Sprintf("repos/%s/%s/issues/%d/timeline", client.owner, client.repository, issueNumber), nil)
	request.Header.Add("Accept", "application/vnd.github.mockingbird-preview")
	var issueTimeline []*IssueTimelineItem
	_, err := client.impl.Do(client.context, request, &issueTimeline)
	if err != nil {
		log.Error(err)
	}
	var timestamp *string
	for _, item := range issueTimeline {
		if *item.Event == "review_requested" {
			timestamp = item.CreatedAt
		}
	}
	if timestamp != nil {
		time, err := time.Parse(time.RFC3339, *timestamp)
		if err != nil {
			return nil
		}
		return &time
	}
	return nil
}
