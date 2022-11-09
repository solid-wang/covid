package continuousintegration

import (
	batchv1 "github.com/solid-wang/covid/pkg/apis/batch/v1"
	"github.com/solid-wang/covid/pkg/gitlab"
	"github.com/solid-wang/covid/pkg/gitlab/projects/pipelines"
)

func createPipeline(ci *batchv1.ContinuousIntegration) error {
	gitlabClient := gitlab.NewGitlabClient(ci.Spec.GitlabHost, ci.Spec.GitlabToken)
	pipeline := &pipelines.Pipeline{
		Ref: ci.Spec.Ref,
		Variables: []pipelines.Variable{
			{
				Key:   "BUILD_IMAGE",
				Value: ci.Spec.BuildImage,
			},
			{
				Key:   "BUILD_DIR",
				Value: ci.Spec.BuildDir,
			},
			{
				Key:   "BUILD_COMMAND",
				Value: ci.Spec.BuildCommand,
			},
			{
				Key:   "ARTIFACT_PATH",
				Value: ci.Spec.ArtifactPath,
			},
			{
				Key:   "FROM_IMAGE",
				Value: ci.Spec.FromImage,
			},
			{
				Key:   "IMAGE",
				Value: ci.Spec.Image,
			},
			{
				Key:   "VERSION",
				Value: ci.Spec.Version,
			},
			{
				Key:   "CI_REGISTRY",
				Value: ci.Spec.Registry,
			},
			{
				Key:   "CI_REGISTRY_USER",
				Value: ci.Spec.RegistryUser,
			},
			{
				Key:   "CI_REGISTRY_PASSWORD",
				Value: ci.Spec.RegistryPassword,
			},
		},
	}
	_, err := gitlabClient.Project().Pipeline(ci.Spec.ProjectID).Create(pipeline)
	return err
}
