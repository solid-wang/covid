package gitlab

import (
	servicev1 "github.com/solid-wang/covid/pkg/apis/service/v1"
	"github.com/solid-wang/covid/pkg/gitlab/client"
	"github.com/solid-wang/covid/pkg/gitlab/projects"
)

const (
	apiV4 = "/api/v4"
)

type Interface interface {
	Project() projects.Interface
}

func NewForGitlab(gitlab *servicev1.Gitlab) *Client {
	c := &client.Client{
		Gitlab: gitlab.Spec.Host + apiV4,
		Token:  gitlab.Spec.Token,
	}
	var g Client
	g.project = projects.NewForClient(c)
	return &g
}

func NewGitlabClient(host, token string) *Client {
	c := &client.Client{
		Gitlab: host + apiV4,
		Token:  token,
	}
	var g Client
	g.project = projects.NewForClient(c)
	return &g
}

type Client struct {
	project *projects.ProjectClient
}

func (c *Client) Project() projects.Interface {
	return c.project
}
