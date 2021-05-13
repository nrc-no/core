package informers

import (
  "fmt"
  v1 "github.com/nrc-no/core/apps/api/pkg/apis/apiextensions/v1"
  "github.com/nrc-no/core/apps/api/pkg/runtime/schema"
  "github.com/nrc-no/core/apps/api/pkg/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
  Informer() cache.SharedIndexInformer
  Lister() cache.GenericLister
}

type genericInformer struct {
  informer cache.SharedIndexInformer
  resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
  return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
  return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
  switch resource {
  // Group=apiextensions.k8s.io, Version=v1
  case v1.SchemeGroupVersion.WithResource("customresourcedefinitions"):
    return &genericInformer{resource: resource.GroupResource(), informer: f.Apiextensions().V1().CustomResourceDefinitions().Informer()}, nil
  }

  return nil, fmt.Errorf("no informer found for %v", resource)
}
