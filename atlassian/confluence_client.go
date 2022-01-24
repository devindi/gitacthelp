package atlassian

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/virtomize/confluence-go-api"
	"io"
	"net/http"
	"strings"
)

type ConfluenceClient struct {
	api     *goconfluence.API
	baseUrl string
}

func NewConfluenceClient(config AtlassianConfig) (*ConfluenceClient, error) {
	api, err := goconfluence.NewAPI(fmt.Sprintf("%s/wiki/rest/api", config.BaseUrl), config.Username, config.Token)
	goconfluence.SetDebug(false)
	if err != nil {
		return nil, err
	}
	return &ConfluenceClient{
		api: api,
	}, nil
}

func (clnt ConfluenceClient) InsertIntoBeggining(spaceId string, pageId string, text string) error {
	var page, err = clnt.api.GetContentByID(pageId, goconfluence.ContentQuery{
		SpaceKey: spaceId,
		Expand:   []string{"body.storage", "version"},
	})
	if err != nil {
		log.Warn("Failed to load page")
		return err
	}

	var updatedBody = text + page.Body.Storage.Value
	_, err = clnt.sendContentRequest(pageId, UpdateBodyRequest{
		Type:  page.Type,
		Title: page.Title,
		Number: Version{
			Number: page.Version.Number + 1,
		},
		Body: Body{
			Storage: Storage{
				Value:          updatedBody,
				Representation: "storage",
			},
		},
	})
	return err

}

//TODO fix clnt.api.UpdateContent
func (clnt ConfluenceClient) sendContentRequest(pageId string, request interface{}) (interface{}, error) {
	var body io.Reader
	if request != nil {
		js, err := json.Marshal(request)
		if err != nil {
			return nil, err
		}
		body = strings.NewReader(string(js))
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/wiki/rest/api/content/%s", clnt.baseUrl, pageId), body)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := clnt.api.Request(req)
	if err != nil {
		return nil, err
	}

	var content interface{}

	err = json.Unmarshal(res, &content)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

type Version struct {
	Number int `json:"number"`
}

type Storage struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}

type Body struct {
	Storage Storage `json:"storage"`
}

type UpdateBodyRequest struct {
	Type   string  `json:"type"`
	Title  string  `json:"title"`
	Number Version `json:"version"`
	Body   Body    `json:"body"`
}
