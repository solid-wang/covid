package webhooks

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/solid-wang/covid/pkg/gitlab/client"
	"strconv"
)

type Interface interface {
	Get(id int) (*Webhook, error)
	List() (*WebhookList, error)
	Add(hook *Webhook) (*Webhook, error)
	Delete(id int) error
}

type WebhookClient struct {
	ProjectID int
	client    *client.Client
}

func NewForClient(client *client.Client) *WebhookClient {
	var webhook WebhookClient
	webhook.client = client
	return &webhook
}

const (
	getUrl    = "/projects/{project_id}/hooks/{id}"
	listUrl   = "/projects/{project_id}/hooks"
	addUrl    = "/projects/{project_id}/hooks"
	deleteUrl = "/projects/{project_id}/hooks/{id}"
)

type WebhookList []Webhook

type Webhook struct {
	HookID                   int    `json:"hook_id"`
	ID                       int    `json:"id"`
	URL                      string `json:"url"`
	ProjectID                int    `json:"project_id"`
	PushEvents               bool   `json:"push_events"`
	PushEventsBranchFilter   string `json:"push_events_branch_filter"`
	IssuesEvents             bool   `json:"issues_events"`
	ConfidentialIssuesEvents bool   `json:"confidential_issues_events"`
	MergeRequestsEvents      bool   `json:"merge_requests_events"`
	TagPushEvents            bool   `json:"tag_push_events"`
	NoteEvents               bool   `json:"note_events"`
	ConfidentialNoteEvents   bool   `json:"confidential_note_events"`
	JobEvents                bool   `json:"job_events"`
	PipelineEvents           bool   `json:"pipeline_events"`
	WikiPageEvents           bool   `json:"wiki_page_events"`
	DeploymentEvents         bool   `json:"deployment_events"`
	ReleasesEvents           bool   `json:"releases_events"`
	EnableSSLVerification    bool   `json:"enable_ssl_verification"`
	Token                    string `json:"token"`
	CreatedAt                string `json:"created_at"`
}

func (w *WebhookClient) Get(id int) (*Webhook, error) {
	result := &Webhook{}
	url := w.client.Gitlab + getUrl
	resp, err := resty.New().SetRetryCount(3).R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(w.client.Token).
		SetPathParams(map[string]string{
			"project_id": strconv.Itoa(w.ProjectID),
			"id":         strconv.Itoa(id),
		}).
		SetResult(result).
		Get(url)
	if err != nil {
		return result, fmt.Errorf("get hook err: %s", err)
	}
	if !resp.IsSuccess() {
		return result, fmt.Errorf("get hook response code: %d", resp.StatusCode())
	}
	return result, nil
}

func (w *WebhookClient) List() (*WebhookList, error) {
	result := &WebhookList{}
	url := w.client.Gitlab + listUrl
	resp, err := resty.New().SetRetryCount(3).R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(w.client.Token).
		SetPathParams(map[string]string{
			"project_id": strconv.Itoa(w.ProjectID),
		}).
		SetResult(result).
		Get(url)
	if err != nil {
		return result, fmt.Errorf("list hooks err: %s", err)
	}
	if !resp.IsSuccess() {
		return result, fmt.Errorf("list hooks response code: %d", resp.StatusCode())
	}
	return result, nil
}

func (w *WebhookClient) Add(hook *Webhook) (*Webhook, error) {
	result := &Webhook{}
	url := w.client.Gitlab + addUrl
	resp, err := resty.New().SetRetryCount(3).R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(w.client.Token).
		SetPathParams(map[string]string{
			"project_id": strconv.Itoa(w.ProjectID),
		}).
		SetBody(hook).
		SetResult(result).
		Post(url)
	if err != nil {
		return result, fmt.Errorf("add hook err: %s", err)
	}
	if !resp.IsSuccess() {
		return result, fmt.Errorf("add hook response code: %d", resp.StatusCode())
	}
	return result, nil
}

func (w *WebhookClient) Delete(id int) error {
	url := w.client.Gitlab + deleteUrl
	resp, err := resty.New().SetRetryCount(3).R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(w.client.Token).
		SetPathParams(map[string]string{
			"project_id": strconv.Itoa(w.ProjectID),
			"id":         strconv.Itoa(id),
		}).
		Delete(url)
	if err != nil {
		return fmt.Errorf("delete hook err: %s", err)
	}
	if !resp.IsSuccess() {
		return fmt.Errorf("delete hook response code: %d", resp.StatusCode())
	}
	return nil
}
