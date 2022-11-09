package options

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	devopscontrollerconfig "github.com/solid-wang/covid/cmd/devops-controller-manager/app/config"
	corev1 "github.com/solid-wang/covid/pkg/apis/core/v1"
	clientset "github.com/solid-wang/covid/pkg/generated/clientset/versioned"
	"github.com/solid-wang/covid/pkg/generated/clientset/versioned/scheme"
	typedcorev1 "github.com/solid-wang/covid/pkg/generated/clientset/versioned/typed/core/v1"
	informers "github.com/solid-wang/covid/pkg/generated/informers/externalversions"
	"github.com/solid-wang/covid/pkg/tools/record"
	"github.com/spf13/pflag"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/rest"
	"time"
)

const (
	// DevOpsControllerManagerUserAgent is the userAgent name when starting devops-controller managers.
	DevOpsControllerManagerUserAgent = "devops-controller-manager"
)

type DevOpsControllerOptions struct {
	Covid string
	FeiShuOpt
	WebServer                       *WebServerOptions
	ContinuousIntegrationController *ContinuousIntegrationControllerOptions
	ContinuousDeploymentController  *ContinuousDeploymentControllerOptions
	ApplicationController           *ApplicationControllerOptions
}

type FeiShuOpt struct {
	ID       string
	Secret   string
	Approval string
}

func NewDevOpsControllerOptions() *DevOpsControllerOptions {
	s := DevOpsControllerOptions{
		WebServer:                       &WebServerOptions{},
		ContinuousIntegrationController: &ContinuousIntegrationControllerOptions{},
		ContinuousDeploymentController:  &ContinuousDeploymentControllerOptions{},
		ApplicationController:           &ApplicationControllerOptions{},
	}
	return &s
}

// Flags returns flags for a specific APIServer by section name
func (s *DevOpsControllerOptions) Flags(fs *pflag.FlagSet) {
	fs.StringVar(&s.Covid, "covid", "", "Covid server http url.")
	fs.StringVar(&s.FeiShuOpt.ID, "feishu-app-id", "", "FeiShu app id.")
	fs.StringVar(&s.FeiShuOpt.Secret, "feishu-app-secret", "", "FeiShu app secret.")
	fs.StringVar(&s.FeiShuOpt.Approval, "feishu-approval", "", "FeiShu devops approval code.")
	s.WebServer.AddFlags(fs)
	s.ContinuousIntegrationController.AddFlags(fs)
	s.ContinuousDeploymentController.AddFlags(fs)
	s.ApplicationController.AddFlags(fs)
}

// ApplyTo fills up controller manager config with options.
func (s *DevOpsControllerOptions) ApplyTo(c *devopscontrollerconfig.Config) error {
	if err := s.WebServer.ApplyTo(&c.WebServer, c.FeiShu, c.Client); err != nil {
		return err
	}
	if err := s.ContinuousIntegrationController.ApplyTo(&c.ComponentConfig.ContinuousIntegration); err != nil {
		return err
	}
	if err := s.ContinuousDeploymentController.ApplyTo(&c.ComponentConfig.ContinuousDeployment); err != nil {
		return err
	}
	if err := s.ApplicationController.ApplyTo(&c.ComponentConfig.Application); err != nil {
		return err
	}
	return nil
}

// Validate is used to validate the options and config before launching the controller manager
func (s *DevOpsControllerOptions) Validate() error {
	var errs []error

	errs = append(errs, s.WebServer.Validate()...)
	errs = append(errs, s.ContinuousIntegrationController.Validate()...)
	errs = append(errs, s.ContinuousDeploymentController.Validate()...)
	errs = append(errs, s.ApplicationController.Validate()...)

	return utilerrors.NewAggregate(errs)
}

// Config return a controller manager config objective
func (s *DevOpsControllerOptions) Config() (*devopscontrollerconfig.Config, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}

	feishuClient := lark.NewClient(s.FeiShuOpt.ID, s.FeiShuOpt.Secret)

	config := &rest.Config{Host: s.Covid}
	clientSet, err := clientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	informerFactory := informers.NewSharedInformerFactory(clientSet, time.Minute*30)

	eventRecorder := createRecorder(clientSet, DevOpsControllerManagerUserAgent)

	c := &devopscontrollerconfig.Config{
		FeiShu:          feishuClient,
		Client:          clientSet,
		InformerFactory: informerFactory,
		EventRecorder:   eventRecorder,
	}
	if err := s.ApplyTo(c); err != nil {
		return nil, err
	}

	return c, nil
}

func createRecorder(Client clientset.Interface, userAgent string) record.EventRecorder {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: Client.CoreV1().Events("")})
	return eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: userAgent})
}
