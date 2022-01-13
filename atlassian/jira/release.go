package jira

import (
	"github.com/andygrunwald/go-jira"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type JiraRelease struct {
	Name string
	Id   int
}

func (clnt JiraClient) CreateRelease(name string) (*JiraRelease, error) {
	var project, _, err = clnt.client.Project.Get(clnt.projectKey)
	if err != nil {
		log.Errorln("Failed to find project: ", clnt.projectKey)
		return nil, err
	}

	projectId, _ := strconv.Atoi(project.ID)

	version, _, err := clnt.client.Version.Create(&jira.Version{
		Name:            name,
		Description:     "",
		Archived:        jira.Bool(false),
		Released:        jira.Bool(false),
		ReleaseDate:     "",
		UserReleaseDate: "",
		ProjectID:       projectId,
		StartDate:       "",
	})
	if err != nil {
		return nil, err
	}
	versionId, _ := strconv.Atoi(version.ID)
	return &JiraRelease{
		Name: version.Name,
		Id:   versionId,
	}, nil
}

func (clnt JiraClient) RenameRelease(oldName string, newName string) (*JiraRelease, error) {

	var project, _, err = clnt.client.Project.Get(clnt.projectKey)
	if err != nil {
		log.Errorln("Failed to find project: ", clnt.projectKey)
		return nil, err
	}

	var versionToRename jira.Version
	for _, version := range project.Versions {
		if version.Name == oldName {
			versionToRename = version
		}
	}

	versionId, _ := strconv.Atoi(versionToRename.ID)

	versionToRename.Name = newName
	_, _, err = clnt.client.Version.Update(&versionToRename)
	if err != nil {
		return nil, err
	}
	return &JiraRelease{
		Name: newName,
		Id:   versionId,
	}, nil
}
