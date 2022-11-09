package web

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	larkapproval "github.com/larksuite/oapi-sdk-go/v3/service/approval/v4"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	batchv1 "github.com/solid-wang/covid/pkg/apis/batch/v1"
	"github.com/solid-wang/covid/pkg/ginserver"
	"github.com/solid-wang/covid/pkg/gitlab/projects/pipelines"
	"github.com/solid-wang/covid/pkg/gitlab/projects/webhooks/events"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"net/http"
	"strconv"
	"strings"
)

func GetAPIGroupInfos(c *WebConfiguration) []*ginserver.APIGroupInfo {
	var infos []*ginserver.APIGroupInfo
	infos = append(infos, ginserver.NewAPIGroupInfo(ginserver.PostVerb, "/:gitlab/webhook/merge_request", c.MergeRequestEvent))
	infos = append(infos, ginserver.NewAPIGroupInfo(ginserver.PostVerb, "/:gitlab/webhook/pipeline", c.PipelineEvent))

	return infos
}

func (w *WebConfiguration) MergeRequestEvent(c *gin.Context) {
	mergeReq := &events.MergeRequest{}
	if err := c.ShouldBindJSON(mergeReq); err != nil {
		klog.Errorf("MergeRequest event object bind err: %s", err)
		c.JSON(http.StatusBadRequest, "MergeRequest event object bind err")
		return
	}
	if mergeReq.State != events.MergeRequestStateMerged {
		c.JSON(http.StatusOK, "ok")
		return
	}
	devopsTags := ExtractDevOpsOption(strings.Split(mergeReq.Description, "\n"))
	if !devopsTags.Enable {
		c.JSON(http.StatusOK, "noting todo.")
		return
	}

	gitlabName := c.Param("gitlab")
	gl, err := w.Client.ServiceV1().Gitlabs().Get(context.Background(), gitlabName, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		fmt.Println("1")
		return
	}
	project := gl.Spec.ProjectIndex[strconv.Itoa(mergeReq.MergeRequestProject.ID)]
	for _, name := range devopsTags.Server {
		product, ok := project.ApplicationProductMap[name]
		if !ok {
			continue
		}
		app, err := w.Client.AppV1().Applications(*product).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			c.JSON(http.StatusNotFound, err)
			fmt.Println("2")
			return
		}
		registry, err := w.Client.ServiceV1().DockerRepositories().Get(context.Background(), app.Spec.DockerRepositoryName, metav1.GetOptions{})
		if err != nil {
			c.JSON(http.StatusNotFound, err)
			fmt.Println("3")
			return
		}
		env := app.Spec.BranchEnvMap[mergeReq.MergeRequestObjectAttributes.TargetBranch]
		image := registry.Spec.Host + "/" + *product + "/" + app.Name
		version := mergeReq.LastCommit.ID[:8]
		ciName := app.Name + "-" + version
		cdName := ciName + "-" + env

		// ci
		ci, err := w.Client.BatchV1().ContinuousIntegrations(*product).Get(context.Background(), ciName, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				// create ci
				continuousIntegration := newContinuousIntegrationsFromApp(app,
					gl.Spec.Host, gl.Spec.Token,
					image, version,
					mergeReq.MergeRequestObjectAttributes.TargetBranch,
					registry.Spec.Host, registry.Spec.User, registry.Spec.Password)
				_, err := w.Client.BatchV1().ContinuousIntegrations(*product).Create(context.Background(), continuousIntegration, metav1.CreateOptions{})
				if err != nil {
					c.JSON(http.StatusInternalServerError, err)
					fmt.Println("4")
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, err)
				fmt.Println("5")
				return
			}
		} else {
			klog.Warningf("ContinuousIntegration %s already exists, status is %s. Ignore build of version %s", ciName, ci.Status.Phase, version)
		}

		// cd
		k8s, err := w.Client.ServiceV1().Kuberneteses().Get(context.Background(), app.Spec.ContinuousDeploymentTemplate.EnvMap[env].KubernetesName, metav1.GetOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			fmt.Println("6")
			return
		}
		newcd := newContinuousDeploymentFromApp(app, []byte(k8s.Spec.Config), env, image, version)
		_, err = w.Client.BatchV1().ContinuousDeployments(*product).Get(context.Background(), cdName, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				_, err := w.Client.BatchV1().ContinuousDeployments(*product).Create(context.Background(), newcd, metav1.CreateOptions{})
				if err != nil {
					c.JSON(http.StatusInternalServerError, err)
					fmt.Println("7")
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, err)
				fmt.Println("8")
				return
			}
		}
		ci, err = w.Client.BatchV1().ContinuousIntegrations(*product).Get(context.Background(), ciName, metav1.GetOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			fmt.Println("9")
			return
		}
		continuousIntegration := ci.DeepCopy()
		continuousIntegration.Status.ContinuousDeploymentTrigger = &cdName
		_, err = w.Client.BatchV1().ContinuousIntegrations(continuousIntegration.Namespace).UpdateStatus(context.Background(), continuousIntegration, metav1.UpdateOptions{})
		if err != nil {
			c.JSON(http.StatusNotFound, err)
			fmt.Println("10")
			fmt.Println(err)
			return
		}
	}
	c.JSON(http.StatusOK, "complete.")
}

func (w *WebConfiguration) PipelineEvent(c *gin.Context) {
	pipeline := events.NewPipeline()
	if err := c.ShouldBindJSON(pipeline); err != nil {
		klog.Errorf("Pipeline event object bind err: %s", err)
		c.JSON(http.StatusBadRequest, "Pipeline event object bind err")
		return
	}
	if pipeline.Tag {
		c.JSON(http.StatusOK, "tag pipeline")
		return
	}

	var product, appName, version string
	for _, variable := range pipeline.ObjectAttributes.Variables {
		if variable.Key == "IMAGE" {
			s := strings.Split(variable.Value, "/")
			if len(s) < 3 {
				c.JSON(http.StatusBadRequest, fmt.Sprintf("Pipeline variable image format err %s", variable.Value))
				return
			}
			product = s[1]
			appName = s[2]
		}
		if variable.Key == "VERSION" {
			version = variable.Value
		}
	}
	ci, err := w.Client.BatchV1().ContinuousIntegrations(product).Get(context.Background(), appName+"-"+version, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}
	continuousIntegrations := ci.DeepCopy()

	if pipeline.Status == pipelines.PipelineFailed {
		continuousIntegrations.Status.Phase = batchv1.DevOpsFailedStatus
		//todo send message
	}
	if pipeline.Status == pipelines.PipelineSuccess {
		continuousIntegrations.Status.Phase = batchv1.DevOpsSuccessStatus
	}

	_, err = w.Client.BatchV1().ContinuousIntegrations(product).UpdateStatus(context.Background(), continuousIntegrations, metav1.UpdateOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (w *WebConfiguration) getFeiShuUserOpenID(username string) (string, error) {
	req := larkcontact.NewBatchGetIdUserReqBuilder().
		UserIdType("open_id").
		Body(larkcontact.NewBatchGetIdUserReqBodyBuilder().
			Emails([]string{username + ".ssc-hn.com"}).
			Build()).
		Build()
	resp, err := w.FeiShu.Contact.User.BatchGetId(context.Background(), req)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return "", fmt.Errorf("%d: %s, %s", resp.Code, resp.RequestId(), resp.Msg)
	}
	if len(resp.Data.UserList) != 0 {
		return *resp.Data.UserList[0].UserId, nil
	}
	return "", nil
}

func (w *WebConfiguration) createFeishuApproval(openID, form string) error {
	req := larkapproval.NewCreateInstanceReqBuilder().
		InstanceCreate(larkapproval.NewInstanceCreateBuilder().
			ApprovalCode(w.FeiShuApproval).
			OpenId(openID).
			Form(form).
			Build()).
		Build()

	resp, err := w.FeiShu.Approval.Instance.Create(context.Background(), req)

	if err != nil {
		return err
	}

	if !resp.Success() {
		return fmt.Errorf("%d: %s, %s", resp.Code, resp.RequestId(), resp.Msg)
	}
	return nil
}
