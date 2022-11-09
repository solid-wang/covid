package ginserver

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
	"log"
	"net/http"
)

const (
	PostVerb   = "POST"
	GetVerb    = "GET"
	DeleteVerb = "DELETE"
	PutVerb    = "PUT"
	PatchVerb  = "PATCH"
)

type Server struct {
	Engine *gin.Engine
}

func NewServer() *Server {
	return &Server{Engine: gin.Default()}
}

type APIGroupInfo struct {
	Verb    string
	Path    string
	Handler gin.HandlerFunc
}

func NewAPIGroupInfo(verb, path string, handler gin.HandlerFunc) *APIGroupInfo {
	return &APIGroupInfo{
		Verb:    verb,
		Path:    path,
		Handler: handler,
	}
}

func (s *Server) InstallAPIs(apiGroupInfos []*APIGroupInfo) error {
	for _, info := range apiGroupInfos {
		switch info.Verb {
		case GetVerb:
			s.Engine.GET(info.Path, info.Handler)
		case PostVerb:
			s.Engine.POST(info.Path, info.Handler)
		case PutVerb:
			s.Engine.POST(info.Path, info.Handler)
		case DeleteVerb:
			s.Engine.DELETE(info.Path, info.Handler)
		case PatchVerb:
			s.Engine.PATCH(info.Path, info.Handler)
		default:
			return fmt.Errorf("verb %s not support", info.Verb)
		}
	}
	return nil
}

func (s *Server) Run(Port int32, stopCh <-chan struct{}) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", Port),
		Handler: s.Engine,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("listen: %s\n", err))
		}
	}()

	<-stopCh
	klog.Warning("Shutdown Server ...")

	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	klog.Warning("Server exiting")
}
