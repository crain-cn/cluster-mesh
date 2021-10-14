package client

import (
	"flag"
	"fmt"
	clustermesh "github.com/crain-cn/cluster-mesh/client/clientset/versioned"
	"github.com/crain-cn/cluster-mesh/client/informers/externalversions"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"
	"time"
)

/*
   使用底层client
*/

func TestClusterMeshClient_GetCluster(t *testing.T) {
	// 需要指定--kubeconfig参数
	flag.Parse()
	config := ctrl.GetConfigOrDie()
	client := clustermesh.NewForConfigOrDie(config)
	clusters, err := client.CloudV1beta1().Clusters().List(v1.ListOptions{})
	factory := externalversions.NewSharedInformerFactory(client, time.Minute)
	informer := factory.Cloud().V1beta1().Clusters().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {

			fmt.Println(1)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println(2)
		},
		DeleteFunc: func(obj interface{}) {},
	})
	factory.Start(make(chan struct{}))
	time.Sleep(time.Second * 100)
	if err != nil {
		t.Error(err)
	}
	t.Log(clusters)
	return
}

/*
   使用封装好的client
*/

func TestClusterMeshClient_ListClusterMesh(t *testing.T) {
	// 需要指定--kubeconfig参数
	flag.Parse()
	config := ctrl.GetConfigOrDie()
	err := InitClusterMeshClient(config)
	if err != nil {
		t.Error(err)
	}
	client := GetClusterMeshClient()
	selector, _ := labels.Parse("clusters.cloud.mesh/zhaoyu")
	meshes, err := client.ListClusterMesh(selector)
	if err != nil {
		t.Error(err)
	}
	t.Log(meshes)
	clusterClient, err := client.GetClusterClient("k8s-test")
	if err != nil {
		t.Error(err)
	}
	clusterClient.AppsV1().Deployments("wangxiao-jichujiagou-common").List(v1.ListOptions{})
}
