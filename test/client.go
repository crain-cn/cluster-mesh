package main

import (
	"flag"
	"fmt"
	clientv1 "github.com/crain-cn/cluster-mesh/client"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 需要指定--kubeconfig参数
	var err error
	var config *rest.Config

	kubeconfig := flag.String("kubeconfig", "/Users/edz/.kube/k8s-32-dev", "absolute path to the kubeconfig file")


	// 使用 ServiceAccount 创建集群配置（InCluster模式）
	if config, err = rest.InClusterConfig(); err != nil {
		// 使用 KubeConfig 文件创建集群配置
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig); err != nil {
			panic(err.Error())
		}
	}

	err = clientv1.InitClusterMeshClient(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := clientv1.GetClusterMeshClient()
	selector, _ := labels.Parse("clusters.cloud.mesh/zhaoyu")
	meshes, err := client.ListClusterMesh(selector)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(meshes)
	clusterClient, err := client.GetClusterClient("k8s-test")
	if err != nil {
		fmt.Println(err)
	}
	clusterClient.AppsV1().Deployments("wangxiao-jichujiagou-common").List(v1.ListOptions{})
}
