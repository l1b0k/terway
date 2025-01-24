/*
Copyright 2021 Terway Authors.

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

package v1beta1

import (
	v1beta1 "github.com/AliyunContainerService/terway/pkg/apis/network.alibabacloud.com/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// NodeRuntimeLister helps list NodeRuntimes.
// All objects returned here must be treated as read-only.
type NodeRuntimeLister interface {
	// List lists all NodeRuntimes in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.NodeRuntime, err error)
	// Get retrieves the NodeRuntime from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta1.NodeRuntime, error)
	NodeRuntimeListerExpansion
}

// nodeRuntimeLister implements the NodeRuntimeLister interface.
type nodeRuntimeLister struct {
	indexer cache.Indexer
}

// NewNodeRuntimeLister returns a new NodeRuntimeLister.
func NewNodeRuntimeLister(indexer cache.Indexer) NodeRuntimeLister {
	return &nodeRuntimeLister{indexer: indexer}
}

// List lists all NodeRuntimes in the indexer.
func (s *nodeRuntimeLister) List(selector labels.Selector) (ret []*v1beta1.NodeRuntime, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.NodeRuntime))
	})
	return ret, err
}

// Get retrieves the NodeRuntime from the index for a given name.
func (s *nodeRuntimeLister) Get(name string) (*v1beta1.NodeRuntime, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("noderuntime"), name)
	}
	return obj.(*v1beta1.NodeRuntime), nil
}