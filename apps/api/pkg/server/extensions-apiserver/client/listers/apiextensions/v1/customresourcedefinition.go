package v1

import (
  v1 "github.com/nrc-no/core/apps/api/pkg/apis/apiextensions/v1"
  "github.com/nrc-no/core/apps/api/pkg/labels"
  "github.com/nrc-no/core/apps/api/pkg/tools/cache"
  "github.com/nrc-no/core/apps/api/pkg/util/exceptions"
)

// CustomResourceDefinitionLister helps list CustomResourceDefinitions.
// All objects returned here must be treated as read-only.
type CustomResourceDefinitionLister interface {
  // List lists all CustomResourceDefinitions in the indexer.
  // Objects returned here must be treated as read-only.
  List(selector labels.Selector) (ret []*v1.CustomResourceDefinition, err error)
  // Get retrieves the CustomResourceDefinition from the index for a given name.
  // Objects returned here must be treated as read-only.
  Get(name string) (*v1.CustomResourceDefinition, error)
}

// customResourceDefinitionLister implements the CustomResourceDefinitionLister interface.
type customResourceDefinitionLister struct {
  indexer cache.Indexer
}

// NewCustomResourceDefinitionLister returns a new CustomResourceDefinitionLister.
func NewCustomResourceDefinitionLister(indexer cache.Indexer) CustomResourceDefinitionLister {
  return &customResourceDefinitionLister{indexer: indexer}
}

// List lists all CustomResourceDefinitions in the indexer.
func (s *customResourceDefinitionLister) List(selector labels.Selector) (ret []*v1.CustomResourceDefinition, err error) {
  err = cache.ListAll(s.indexer, selector, func(m interface{}) {
    ret = append(ret, m.(*v1.CustomResourceDefinition))
  })
  return ret, err
}

// Get retrieves the CustomResourceDefinition from the index for a given name.
func (s *customResourceDefinitionLister) Get(name string) (*v1.CustomResourceDefinition, error) {
  obj, exists, err := s.indexer.GetByKey(name)
  if err != nil {
    return nil, err
  }
  if !exists {
    return nil, exceptions.NewNotFound(v1.Resource("customresourcedefinition"), name)
  }
  return obj.(*v1.CustomResourceDefinition), nil
}
