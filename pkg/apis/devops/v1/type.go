package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	FinalizerContinuousIntegration = "continuousIntegration"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContinuousIntegrationList is a list of ContinuousIntegration objects.
type ContinuousIntegrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ContinuousIntegration `json:"items"`
}

// ContinuousIntegrationSpec is the specification of a ContinuousIntegration.
type ContinuousIntegrationSpec struct {
	//GitlabName      string          `json:"gitlabName"`
	//GitlabProjectID int             `json:"gitlabProjectID"`
	CIConfigPath       string `json:"ciConfigPath"`
	BuildImage         string `json:"buildImage"`
	FromImage          string `json:"fromImage"`
	BuildDir           string `json:"buildDir"`
	BuildCommand       string `json:"buildCommand"`
	ArtifactPath       string `json:"artifactPath"`
	Image              string `json:"image"`
	BuiltCommitHistory []string
	//BranchHashMap   map[string]Hash `json:"branchHashMap"`
}

// ContinuousIntegrationStatus is the status of a ContinuousIntegration.
type ContinuousIntegrationStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContinuousIntegration is a ContinuousIntegration type with a spec and a status.
type ContinuousIntegration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContinuousIntegrationSpec   `json:"spec,omitempty"`
	Status ContinuousIntegrationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContinuousDeploymentList is a list of ContinuousDeployment objects.
type ContinuousDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ContinuousDeployment `json:"items"`
}

// ContinuousDeploymentSpec is the specification of a ContinuousDeployment.
type ContinuousDeploymentSpec struct {
	ManifestRepo string            `json:"manifestRepo"`
	EnvDeployMap map[string]Deploy `json:"envDeployMap"`
}

type Deploy struct {
	Approval               bool   `json:"approval"`
	KubernetesName         string `json:"kubernetesName"`
	KubernetesNamespace    string `json:"kubernetesNamespace"`
	DeployedManifestCommit string `json:"deployedManifestCommit"`
	Deploying              `json:"deploying"`
}

type Deploying struct {
	ManifestCommit       string `json:"manifestCommit"`
	ApprovalInstanceCode string `json:"approvalInstanceCode"`
}

// ContinuousDeploymentStatus is the status of a ContinuousDeployment.
type ContinuousDeploymentStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContinuousDeployment is a ContinuousDeployment type with a spec and a status.
type ContinuousDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContinuousDeploymentSpec   `json:"spec,omitempty"`
	Status ContinuousDeploymentStatus `json:"status,omitempty"`
}
