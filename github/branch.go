package github

import (
	"fmt"
	"github.com/google/go-github/v28/github"
	log "github.com/sirupsen/logrus"
	"strings"
)

type GithubBranch struct {
	Name        string
	AuthorEmail string
}

func (branch GithubBranch) GetAuthor() string {
	return branch.AuthorEmail[:strings.IndexByte(branch.AuthorEmail, '@')]
}

func (branch GithubBranch) GetUrl(config GithubConfig) string {
	return fmt.Sprintf("<https://github.com/%s/%s/tree/%s|%s>", config.Owner, config.Repository, branch.Name, branch.Name)
}

func (client GithubClient) CreateBranch(base string, name string) (*string, error) {
	originBranchRef, _, err := client.impl.Git.GetRef(client.context, client.owner, client.repository, fmt.Sprintf("heads/%s", base))
	if err != nil {
		log.Errorln("Failed to get base branch")
		return nil, err
	}

	lastCommitSha := originBranchRef.Object.SHA
	log.WithField("commit", *lastCommitSha).Println("Creating release branch on commit")

	_, _, err = client.impl.Git.CreateRef(client.context, client.owner, client.repository, &github.Reference{
		Ref: github.String(fmt.Sprintf("refs/heads/%s", name)),
		Object: &github.GitObject{
			SHA: lastCommitSha,
		},
	})
	if err != nil {
		return nil, err
	}
	branchUrl := fmt.Sprintf("https://github.com/%s/%s/tree/%s", client.owner, client.repository, name)
	return &branchUrl, nil
}

func (client GithubClient) GetBranches() ([]GithubBranch, error) {
	branchOpt := github.ListOptions{PerPage: 150}
	branches, _, err := client.impl.Repositories.ListBranches(client.context, client.owner, client.repository, &branchOpt)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var result []GithubBranch
	for _, branch := range branches {
		branchInfo, _, err := client.impl.Repositories.GetBranch(client.context, client.owner, client.repository, *branch.Name)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if !shouldIgnoreBranch(*branchInfo.Name) {
			log.Debug("Got branch ", *branch.Name)
			result = append(result, GithubBranch{
				Name:        *branchInfo.Name,
				AuthorEmail: *branchInfo.Commit.Commit.Author.Email,
			})
		}
	}
	return result, nil
}

func shouldIgnoreBranch(name string) bool {
	if strings.HasPrefix(name, "release/") {
		return true
	}
	if strings.HasPrefix(name, "test/") {
		return true
	}
	if name == "master" {
		return true
	}
	if name == "develop" {
		return true
	}
	return false
}
