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
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1beta1 "github.com/crain-cn/cluster-mesh/api/cloud.mesh/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeClusterMeshes implements ClusterMeshInterface
type FakeClusterMeshes struct {
	Fake *FakeCloudV1beta1
	ns   string
}

var clustermeshesResource = schema.GroupVersionResource{Group: "cloud.mesh", Version: "v1beta1", Resource: "clustermeshes"}

var clustermeshesKind = schema.GroupVersionKind{Group: "cloud.mesh", Version: "v1beta1", Kind: "ClusterMesh"}

// Get takes name of the clusterMesh, and returns the corresponding clusterMesh object, and an error if there is any.
func (c *FakeClusterMeshes) Get(name string, options v1.GetOptions) (result *v1beta1.ClusterMesh, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(clustermeshesResource, c.ns, name), &v1beta1.ClusterMesh{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ClusterMesh), err
}

// List takes label and field selectors, and returns the list of ClusterMeshes that match those selectors.
func (c *FakeClusterMeshes) List(opts v1.ListOptions) (result *v1beta1.ClusterMeshList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(clustermeshesResource, clustermeshesKind, c.ns, opts), &v1beta1.ClusterMeshList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.ClusterMeshList{ListMeta: obj.(*v1beta1.ClusterMeshList).ListMeta}
	for _, item := range obj.(*v1beta1.ClusterMeshList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterMeshes.
func (c *FakeClusterMeshes) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(clustermeshesResource, c.ns, opts))

}

// Create takes the representation of a clusterMesh and creates it.  Returns the server's representation of the clusterMesh, and an error, if there is any.
func (c *FakeClusterMeshes) Create(clusterMesh *v1beta1.ClusterMesh) (result *v1beta1.ClusterMesh, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(clustermeshesResource, c.ns, clusterMesh), &v1beta1.ClusterMesh{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ClusterMesh), err
}

// Update takes the representation of a clusterMesh and updates it. Returns the server's representation of the clusterMesh, and an error, if there is any.
func (c *FakeClusterMeshes) Update(clusterMesh *v1beta1.ClusterMesh) (result *v1beta1.ClusterMesh, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(clustermeshesResource, c.ns, clusterMesh), &v1beta1.ClusterMesh{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ClusterMesh), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeClusterMeshes) UpdateStatus(clusterMesh *v1beta1.ClusterMesh) (*v1beta1.ClusterMesh, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(clustermeshesResource, "status", c.ns, clusterMesh), &v1beta1.ClusterMesh{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ClusterMesh), err
}

// Delete takes name of the clusterMesh and deletes it. Returns an error if one occurs.
func (c *FakeClusterMeshes) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(clustermeshesResource, c.ns, name), &v1beta1.ClusterMesh{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterMeshes) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(clustermeshesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1beta1.ClusterMeshList{})
	return err
}

// Patch applies the patch and returns the patched clusterMesh.
func (c *FakeClusterMeshes) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.ClusterMesh, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(clustermeshesResource, c.ns, name, pt, data, subresources...), &v1beta1.ClusterMesh{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ClusterMesh), err
}