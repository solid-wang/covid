package pipelines

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/solid-wang/covid/pkg/gitlab/client"
	"strconv"
)

type Interface interface {
	Create(pipeline *Pipeline) (*Pipeline, error)
}

type PipelineClient struct {
	ProjectID int
	client    *client.Client
}

func NewForClient(client *client.Client) *PipelineClient {
	var pipeline PipelineClient
	pipeline.client = client
	return &pipeline
}

const (
	createUrl       = "/projects/{project_id}/pipeline"
	PipelineSuccess = "success"
	PipelineFailed  = "failed"
)

type Pipeline struct {
	ID         int    `json:"id,omitempty"`
	IID        int    `json:"iid,omitempty"`
	ProjectID  int    `json:"project_id,project_id"`
	SHA        string `json:"sha,omitempty"`
	Tag        bool   `json:"tag,omitempty"`
	Ref        string `json:"ref"`
	StartedAt  string `json:"started_at,omitempty"`
	FinishedAt string `json:"finished_at,omitempty"`
	// Duration second
	Duration  int        `json:"duration,omitempty"`
	Variables []Variable `json:"variables,omitempty"`
}

type Variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (p *PipelineClient) Create(pipeline *Pipeline) (*Pipeline, error) {
	result := &Pipeline{}
	url := p.client.Gitlab + createUrl
	resp, err := resty.New().SetRetryCount(3).R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(p.client.Token).
		SetPathParams(map[string]string{
			"project_id": strconv.Itoa(p.ProjectID),
		}).
		SetBody(pipeline).
		SetResult(result).
		Post(url)
	if err != nil {
		return result, fmt.Errorf("create pipeline err: %s", err)
	}
	if !resp.IsSuccess() {
		return result, fmt.Errorf("create pipeline response code: %d", resp.StatusCode())
	}
	return result, nil
}
