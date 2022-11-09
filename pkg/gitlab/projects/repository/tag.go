package repository

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/solid-wang/covid/pkg/gitlab/client"
	"strconv"
)

type Interface interface {
	GetTag(tagName string) (*Tag, error)
}

type TagClient struct {
	ProjectID int
	client    *client.Client
}

func NewForClient(client *client.Client) *TagClient {
	var tag TagClient
	tag.client = client
	return &tag
}

type Tag struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Target  string `json:"target"`
}

const (
	getUrl = "/projects/{project_id}/repository/tags/{tag_name}"
)

func (t *TagClient) GetTag(tagName string) (*Tag, error) {
	result := &Tag{}
	url := t.client.Gitlab + getUrl
	resp, err := resty.New().SetRetryCount(3).R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(t.client.Token).
		SetPathParams(map[string]string{
			"project_id": strconv.Itoa(t.ProjectID),
			"tag_name":   tagName,
		}).
		SetResult(result).
		Get(url)
	if err != nil {
		return result, fmt.Errorf("get tag err: %s", err)
	}
	if !resp.IsSuccess() {
		return result, fmt.Errorf("get tag response code: %d", resp.StatusCode())
	}
	return result, nil
}
