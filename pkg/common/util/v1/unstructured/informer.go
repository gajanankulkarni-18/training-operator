// Copyright 2021 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

// Package unstructured is the package for unstructured informer,
// which is from https://github.com/argoproj/argo/blob/master/util/unstructured/unstructured.go
// This is a temporary solution for https://github.com/kubeflow/training-operator/issues/561
package unstructured

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"

	informer "github.com/kubeflow/training-operator/pkg/client/informers/externalversions/tensorflow/v1"
	lister "github.com/kubeflow/training-operator/pkg/client/listers/tensorflow/v1"
)

type UnstructuredInformer struct {
	informer cache.SharedIndexInformer
}

func NewTFJobInformer(resource schema.GroupVersionResource, client dynamic.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) informer.TFJobInformer {
	return &UnstructuredInformer{
		informer: newUnstructuredInformer(resource, client, namespace, resyncPeriod, indexers),
	}
}

func (f *UnstructuredInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

func (f *UnstructuredInformer) Lister() lister.TFJobLister {
	return lister.NewTFJobLister(f.Informer().GetIndexer())
}

// newUnstructuredInformer constructs a new informer for Unstructured type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func newUnstructuredInformer(resource schema.GroupVersionResource, client dynamic.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return newFilteredUnstructuredInformer(resource, client, namespace, resyncPeriod, indexers)
}

// newFilteredUnstructuredInformer constructs a new informer for Unstructured type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func newFilteredUnstructuredInformer(resource schema.GroupVersionResource, client dynamic.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return client.Resource(resource).Namespace(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return client.Resource(resource).Namespace(namespace).Watch(context.TODO(), options)
			},
		},
		&unstructured.Unstructured{},
		resyncPeriod,
		indexers,
	)
}
