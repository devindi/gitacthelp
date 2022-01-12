package jira

type version struct {
	Name string `json:"name"`
}

type add struct {
	Add version `json:"add"`
}

type updateBody struct {
	FixVersions []add `json:"fixVersions"`
}

func (jira JiraClient) SetFixVersion(issueKey string, versionName string) error {
	payload := make(map[string]interface{})

	payload["update"] = updateBody{
		FixVersions: []add{
			{
				Add: version{
					Name: versionName,
				},
			},
		},
	}

	_, err := jira.client.Issue.UpdateIssue(issueKey, payload)
	return err
}
