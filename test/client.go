package main

import (
	"fmt"
	clustermesh "github.com/crain-cn/cluster-mesh/client/clientset/versioned"
	"github.com/crain-cn/cluster-mesh/client/informers/externalversions"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

func main() {
	// 需要指定--kubeconfig参数
	var err error
	var config *rest.Config

	kubeconfig := "/Users/edz/.kube/k8s-32-dev"
	if config, err = clientcmd.BuildConfigFromFlags("", kubeconfig); err != nil {
		panic(err.Error())
	}

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
		fmt.Println(err)
	}
	fmt.Println(clusters)
}
