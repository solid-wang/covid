package options

import (
	cmdbcontrollerconfig "github.com/solid-wang/covid/cmd/cmdb-controller-manager/app/config"
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
	// CMDBControllerManagerUserAgent is the userAgent name when starting devops-controller managers.
	CMDBControllerManagerUserAgent = "cmdb-controller-manager"
)

type CMDBControllerOptions struct {
	Covid            string
	ServerController *ServerControllerOptions
}

func NewCMDBControllerOptions() *CMDBControllerOptions {
	s := CMDBControllerOptions{
		ServerController: &ServerControllerOptions{},
	}
	return &s
}

// Flags returns flags for a specific APIServer by section name
func (s *CMDBControllerOptions) Flags(fs *pflag.FlagSet) {
	fs.StringVar(&s.Covid, "covid", "", "Covid server http url.")
	s.ServerController.AddFlags(fs)
}

// ApplyTo fills up controller manager config with options.
func (s *CMDBControllerOptions) ApplyTo(c *cmdbcontrollerconfig.Config) error {
	if err := s.ServerController.ApplyTo(&c.ComponentConfig.Server); err != nil {
		return err
	}
	return nil
}

// Validate is used to validate the options and config before launching the controller manager
func (s *CMDBControllerOptions) Validate() error {
	var errs []error

	errs = append(errs, s.ServerController.Validate()...)

	return utilerrors.NewAggregate(errs)
}

// Config return a controller manager config objective
func (s *CMDBControllerOptions) Config() (*cmdbcontrollerconfig.Config, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}

	config := &rest.Config{Host: s.Covid}
	clientSet, err := clientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	informerFactory := informers.NewSharedInformerFactory(clientSet, time.Second*30)

	eventRecorder := createRecorder(clientSet, CMDBControllerManagerUserAgent)

	c := &cmdbcontrollerconfig.Config{
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
