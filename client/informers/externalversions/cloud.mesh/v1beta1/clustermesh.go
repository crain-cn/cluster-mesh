/*
Copyright 2021.

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
// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	time "time"

	cloudmeshv1beta1 "github.com/crain-cn/cluster-mesh/api/cloud.mesh/v1beta1"
	versioned "github.com/crain-cn/cluster-mesh/client/clientset/versioned"
	internalinterfaces "github.com/crain-cn/cluster-mesh/client/informers/externalversions/internalinterfaces"
	v1beta1 "github.com/crain-cn/cluster-mesh/client/listers/cloud.mesh/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ClusterMeshInformer provides access to a shared informer and lister for
// ClusterMeshes.
type ClusterMeshInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.ClusterMeshLister
}

type clusterMeshInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewClusterMeshInformer constructs a new informer for ClusterMesh type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewClusterMeshInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredClusterMeshInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredClusterMeshInformer constructs a new informer for ClusterMesh type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredClusterMeshInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CloudV1beta1().ClusterMeshes(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CloudV1beta1().ClusterMeshes(namespace).Watch(context.TODO(), options)
			},
		},
		&cloudmeshv1beta1.ClusterMesh{},
		resyncPeriod,
		indexers,
	)
}

func (f *clusterMeshInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredClusterMeshInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *clusterMeshInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&cloudmeshv1beta1.ClusterMesh{}, f.defaultInformer)
}

func (f *clusterMeshInformer) Lister() v1beta1.ClusterMeshLister {
	return v1beta1.NewClusterMeshLister(f.Informer().GetIndexer())
}
