// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/solid-wang/covid/pkg/apis/service/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// GitlabLister helps list Gitlabs.
// All objects returned here must be treated as read-only.
type GitlabLister interface {
	// List lists all Gitlabs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Gitlab, err error)
	// Get retrieves the Gitlab from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Gitlab, error)
	GitlabListerExpansion
}

// gitlabLister implements the GitlabLister interface.
type gitlabLister struct {
	indexer cache.Indexer
}

// NewGitlabLister returns a new GitlabLister.
func NewGitlabLister(indexer cache.Indexer) GitlabLister {
	return &gitlabLister{indexer: indexer}
}

// List lists all Gitlabs in the indexer.
func (s *gitlabLister) List(selector labels.Selector) (ret []*v1.Gitlab, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Gitlab))
	})
	return ret, err
}

// Get retrieves the Gitlab from the index for a given name.
func (s *gitlabLister) Get(name string) (*v1.Gitlab, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("gitlab"), name)
	}
	return obj.(*v1.Gitlab), nil
}
