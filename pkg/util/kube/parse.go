package kube

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/resource"
	"strings"
)

func ParseYamlToObject(kubeConfig string, filename []string) (runtime.Object, error) {
	clusterClientGetter := NewClusterClientGetter(kubeConfig)
	r := resource.NewBuilder(clusterClientGetter).Unstructured().Flatten().FilenameParam(false, &resource.FilenameOptions{
		Filenames: filename,
	}).ContinueOnError().Do()
	return r.Object()
}

func ParseNamespacedName(namespacedName string) (*types.NamespacedName, error) {
	infos := strings.Split(namespacedName, "/")
	if len(infos) != 2 {
		return nil, fmt.Errorf("illgal resource name: %s", namespacedName)
	}
	return &types.NamespacedName{
		Namespace: infos[0],
		Name:      infos[1],
	}, nil
}
