package projects

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/solid-wang/covid/pkg/gitlab/client"
	"github.com/solid-wang/covid/pkg/gitlab/projects/pipelines"
	"github.com/solid-wang/covid/pkg/gitlab/projects/repository"
	"github.com/solid-wang/covid/pkg/gitlab/projects/webhooks"
	"strconv"
)

type Interface interface {
	Create(projectName string, namespaceID int) (*Project, error)
	Get(id int) (*Project, error)
	EditCIConfigPath(id int, ciConfigPath string) (*Project, error)
	Webhook(projectID int) webhooks.Interface
	Pipeline(projectID int) pipelines.Interface
	Tag(projectID int) repository.Interface
}

type ProjectClient struct {
	client   *client.Client
	webhook  *webhooks.WebhookClient
	pipeline *pipelines.PipelineClient
	tag      *repository.TagClient
}

func NewForClient(client *client.Client) *ProjectClient {
	var project ProjectClient
	project.client = client
	project.webhook = webhooks.NewForClient(client)
	project.pipeline = pipelines.NewForClient(client)
	project.tag = repository.NewForClient(client)
	return &project
}

func (p *ProjectClient) Webhook(projectID int) webhooks.Interface {
	p.webhook.ProjectID = projectID
	return p.webhook
}

func (p *ProjectClient) Pipeline(projectID int) pipelines.Interface {
	p.pipeline.ProjectID = projectID
	return p.pipeline
}

func (p *ProjectClient) Tag(projectID int) repository.Interface {
	p.tag.ProjectID = projectID
	return p.tag
}

const (
	createUrl = "/projects"
	getUrl    = "/projects/{id}"
	editUrl   = "/projects/{id}"
)

type Project struct {
	ID           int `json:"id"`
	Namespace    `json:"namespace,omitempty"`
	CIConfigPath string `json:"ci_config_path"`
	RunnersToken string `json:"runners_token,omitempty"`
}

type Namespace struct {
	ID   int    `json:"id,omitempty"`
	Kind string `json:"kind,omitempty"`
	Name string `json:"name,omitempty"`
}

func (p *ProjectClient) Create(projectName string, namespaceID int) (*Project, error) {
	result := &Project{}
	url := p.client.Gitlab + createUrl
	resp, err := resty.New().SetRetryCount(3).R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(p.client.Token).
		SetBody(&map[string]interface{}{
			"name":         projectName,
			"namespace_id": namespaceID,
		}).
		SetResult(result).
		Post(url)
	if err != nil {
		return result, fmt.Errorf("create project err: %s", err)
	}
	if !resp.IsSuccess() {
		return result, fmt.Errorf("create project response code: %d", resp.StatusCode())
	}
	return result, nil
}

func (p *ProjectClient) Get(id int) (*Project, error) {
	result := &Project{}
	url := p.client.Gitlab + getUrl
	resp, err := resty.New().SetRetryCount(3).R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(p.client.Token).
		SetPathParams(map[string]string{
			"id": strconv.Itoa(id),
		}).
		SetResult(result).
		Get(url)
	if err != nil {
		return result, fmt.Errorf("get project err: %s", err)
	}
	if !resp.IsSuccess() {
		return result, fmt.Errorf("get project response code: %d", resp.StatusCode())
	}
	return result, nil
}

func (p *ProjectClient) EditCIConfigPath(id int, ciConfigPath string) (*Project, error) {
	result := &Project{}
	url := p.client.Gitlab + editUrl
	resp, err := resty.New().SetRetryCount(3).R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(p.client.Token).
		SetPathParams(map[string]string{
			"id": strconv.Itoa(id),
		}).
		SetBody(&Project{CIConfigPath: ciConfigPath}).
		SetResult(result).
		Put(url)
	if err != nil {
		return result, fmt.Errorf("edit project err: %s", err)
	}
	if !resp.IsSuccess() {
		return result, fmt.Errorf("edit project response code: %d", resp.StatusCode())
	}
	return result, nil
}
