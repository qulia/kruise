/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/openkruise/kruise/pkg/apis/apps/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// NodeImageLister helps list NodeImages.
type NodeImageLister interface {
	// List lists all NodeImages in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.NodeImage, err error)
	// Get retrieves the NodeImage from the index for a given name.
	Get(name string) (*v1alpha1.NodeImage, error)
	NodeImageListerExpansion
}

// nodeImageLister implements the NodeImageLister interface.
type nodeImageLister struct {
	indexer cache.Indexer
}

// NewNodeImageLister returns a new NodeImageLister.
func NewNodeImageLister(indexer cache.Indexer) NodeImageLister {
	return &nodeImageLister{indexer: indexer}
}

// List lists all NodeImages in the indexer.
func (s *nodeImageLister) List(selector labels.Selector) (ret []*v1alpha1.NodeImage, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.NodeImage))
	})
	return ret, err
}

// Get retrieves the NodeImage from the index for a given name.
func (s *nodeImageLister) Get(name string) (*v1alpha1.NodeImage, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("nodeimage"), name)
	}
	return obj.(*v1alpha1.NodeImage), nil
}
