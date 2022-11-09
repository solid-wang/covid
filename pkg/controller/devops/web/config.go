package web

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	clientset "github.com/solid-wang/covid/pkg/generated/clientset/versioned"
)

// WebConfiguration holds configuration for a web controller-manager
type WebConfiguration struct {
	Port int32

	FeiShu         *lark.Client
	FeiShuApproval string
	// the general covid client
	Client *clientset.Clientset
}
