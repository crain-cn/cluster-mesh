package client

import (
	"fmt"
	"github.com/crain-cn/cluster-mesh/api/cloud.mesh/v1beta1"
	"github.com/crain-cn/cluster-mesh/client/clientset/versioned"
	"github.com/crain-cn/cluster-mesh/client/informers/externalversions"
	informer "github.com/crain-cn/cluster-mesh/client/informers/externalversions/cloud.mesh/v1beta1"
	kubeutil "github.com/crain-cn/cluster-mesh/pkg/util/kube"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"
)

/*
	封装Client用于获取相关CRD
*/

var client *ClusterMeshClient

type ClusterMeshClient struct {
	config *rest.Config
	versioned.Interface
	clusterInformer     informer.ClusterInformer
	clusterMeshInformer informer.ClusterMeshInformer
	stop                chan struct{}
}

func InitClusterMeshClient(config *rest.Config) error {
	return InitClusterMeshClientEnableSync(config, "")
}

func InitClusterMeshClientEnableSync(config *rest.Config, syncCluster string) error {
	var err error
	if config == nil {
		kubeconfigpath := os.Getenv("KUBECONFIG")
		configs, err := kubeutil.LoadClusterConfigs(kubeconfigpath, "")
		if err != nil {
			return err
		}
		c := configs[""]
		config = &c
	}
	meshClient, err := versioned.NewForConfig(config)
	if err != nil {
		return err
	}
	stopChan := make(chan struct{})
	factory := externalversions.NewSharedInformerFactory(meshClient, time.Minute)
	go factory.Cloud().V1beta1().ClusterMeshes().Informer().Run(stopChan)
	go factory.Cloud().V1beta1().Clusters().Informer().Run(stopChan)
	newClient := &ClusterMeshClient{
		config:              config,
		clusterInformer:     factory.Cloud().V1beta1().Clusters(),
		clusterMeshInformer: factory.Cloud().V1beta1().ClusterMeshes(),
		Interface:           meshClient,
		stop:                stopChan,
	}
	if syncCluster != "" {
		cluster, err := meshClient.CloudV1beta1().Clusters().Get(syncCluster, v1.GetOptions{})
		if err != nil {
			return err
		}
		err = newClient.syncClusterMesh(cluster)
		if err != nil {
			return err
		}
	}
	factory.Start(stopChan)
	factory.WaitForCacheSync(stopChan)
	client = newClient
	return nil
}

func GetClusterMeshClient() *ClusterMeshClient {
	return client
}

func (c *ClusterMeshClient) GetCluster(name string) (*v1beta1.Cluster, error) {
	return c.clusterInformer.Lister().Get(name)
}

func (c *ClusterMeshClient) GetClusterMesh(namespace, name string) (*v1beta1.ClusterMesh, error) {
	return c.clusterMeshInformer.Lister().ClusterMeshes(namespace).Get(name)
}

func (c *ClusterMeshClient) GetClusterConfig(name string) (*rest.Config, error) {
	cluster, err := c.GetCluster(name)
	if err != nil {
		return nil, err
	}
	return kubeutil.LoadConfig(cluster.Spec.KubeConfig)
}

func (c *ClusterMeshClient) GetClusterClient(name string) (*kubernetes.Clientset, error) {
	config, err := c.GetClusterConfig(name)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func (c *ClusterMeshClient) ListClusterMesh(selector labels.Selector) ([]*v1beta1.ClusterMesh, error) {
	return c.clusterMeshInformer.Lister().List(selector)
}

func (c *ClusterMeshClient) syncClusterMesh(cluster *v1beta1.Cluster) error {
	config, err := kubeutil.LoadConfig(cluster.Spec.KubeConfig)
	if err != nil {
		return err
	}
	c.syncCluster(config,cluster)
	// c.syncMesh(config,cluster)
	return nil
}


func  (c *ClusterMeshClient) syncMesh(config *rest.Config,cluster *v1beta1.Cluster) {
	log := ctrl.Log.WithName("controllers").WithName("syncMesh")
	c.clusterMeshInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mesh, ok := obj.(*v1beta1.ClusterMesh)
			if !ok {
				return
			}
			meshCopy := mesh.DeepCopy()
			meshCopy.ResourceVersion = ""
			syncClient, err := versioned.NewForConfig(config)
			if err != nil {
				return
			}
			_, err = syncClient.CloudV1beta1().ClusterMeshes(mesh.Namespace).Create(meshCopy)
			if errors.IsAlreadyExists(err) {
				log.Error(err, fmt.Sprintf("新增clustermesh %s 从属集群 %s 已存在 ", mesh.Name, cluster.Name))
			} else if err != nil {
				log.Error(err, fmt.Sprintf("新增clustermesh %s 从属集群 %s 创建失败 ", mesh.Name, cluster.ClusterName))
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			mesh, ok := newObj.(*v1beta1.ClusterMesh)
			if !ok {
				return
			}
			meshCopy := mesh.DeepCopy()
			meshCopy.ResourceVersion = ""
			syncClient, err := versioned.NewForConfig(config)
			if err != nil {
				return
			}
			_, err = syncClient.CloudV1beta1().ClusterMeshes(mesh.Namespace).Update(meshCopy)
			if errors.IsConflict(err) {
				err = syncClient.CloudV1beta1().ClusterMeshes(meshCopy.Namespace).Delete(meshCopy.Name, &v1.DeleteOptions{})
				if err != nil && !errors.IsNotFound(err) {
					log.Error(err, fmt.Sprintf("变更clustermesh %s 从属集群%s 冲突且删除失败 ", meshCopy.Name, cluster.Name))
					return
				}
				_, err = syncClient.CloudV1beta1().ClusterMeshes(meshCopy.Namespace).Create(meshCopy)
				if err != nil {
					log.Error(err, fmt.Sprintf("变更clustermesh %s 从属集群%s 冲突且新建失败 ", meshCopy.Name, cluster.Name))
					return
				}
			} else if errors.IsNotFound(err) {
				_, err := syncClient.CloudV1beta1().ClusterMeshes(mesh.Namespace).Create(meshCopy)
				if err != nil {
					log.Error(err, fmt.Sprintf("变更clustermesh %s 从属集群%s不存在且创建失败 ", meshCopy.Name, cluster.Name))
				}
			} else if err != nil {
				log.Error(err, fmt.Sprintf("变更clustermesh %s 从属集群%s同步失败 ", mesh.Name, cluster.Name))
			}
		},
		DeleteFunc: func(obj interface{}) {
			mesh, ok := obj.(*v1beta1.ClusterMesh)
			if !ok {
				return
			}
			syncClient, err := versioned.NewForConfig(config)
			if err != nil {
				return
			}
			err = syncClient.CloudV1beta1().ClusterMeshes(mesh.Namespace).Delete(mesh.Name, &v1.DeleteOptions{})
			if err != nil && !errors.IsNotFound(err) {
				log.Error(err, fmt.Sprintf("删除clustermesh %s 从属集群%s不存在 ", mesh.Name, cluster.Name))
			}
		},
	})
}


func  (c *ClusterMeshClient) syncCluster(config *rest.Config,cluster *v1beta1.Cluster) {
	log := ctrl.Log.WithName("controllers").WithName("syncCluster")
	c.clusterInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			c, ok := obj.(*v1beta1.Cluster)
			if !ok {
				return
			}
			clusterCopy := c.DeepCopy()
			clusterCopy.ResourceVersion = ""
			syncClient, err := versioned.NewForConfig(config)
			if err != nil {
				return
			}
			_, err = syncClient.CloudV1beta1().Clusters().Create(clusterCopy)
			if errors.IsAlreadyExists(err) {
				log.Error(err, fmt.Sprintf("创建cluster %s 从属集群%s已存在 ", clusterCopy.Name, cluster.Name))
			} else if err != nil {
				log.Error(err, fmt.Sprintf("创建cluster %s 从属集群%s 创建失败 ", clusterCopy.Name, cluster.Name))
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			c, ok := newObj.(*v1beta1.Cluster)
			if !ok {
				return
			}
			clusterCopy := c.DeepCopy()
			clusterCopy.ResourceVersion = ""
			syncClient, err := versioned.NewForConfig(config)
			if err != nil {
				return
			}
			_, err = syncClient.CloudV1beta1().Clusters().Update(clusterCopy)
			if errors.IsConflict(err) {

				err = syncClient.CloudV1beta1().Clusters().Delete(clusterCopy.Name, &v1.DeleteOptions{})
				if err != nil && !errors.IsNotFound(err) {
					log.Error(err, fmt.Sprintf("变更cluster %s 从属集群%s 冲突且删除失败 ", clusterCopy.Name, cluster.Name))
					return
				}
				_, err = syncClient.CloudV1beta1().Clusters().Create(clusterCopy)
				if err != nil {
					log.Error(err, fmt.Sprintf("变更cluster %s 从属集群%s 冲突且创建失败 ", clusterCopy.Name, cluster.Name))
					return
				}
			} else if errors.IsNotFound(err) {
				_, err := syncClient.CloudV1beta1().Clusters().Create(clusterCopy)
				if err != nil {
					log.Error(err, fmt.Sprintf("变更cluster %s 从属集群%s 不存在且创建失败 ", clusterCopy.Name, cluster.Name))
				}
			} else if err != nil {
				log.Error(err, fmt.Sprintf("变更cluster %s 从属集群%s 更新失败 ", c.Name, cluster.Name))
			}
		},
		DeleteFunc: func(obj interface{}) {
			c, ok := obj.(*v1beta1.Cluster)
			if !ok {
				return
			}
			syncClient, err := versioned.NewForConfig(config)
			if err != nil {
				return
			}
			err = syncClient.CloudV1beta1().Clusters().Delete(c.Name, &v1.DeleteOptions{})
			if err != nil && !errors.IsNotFound(err) {
				log.Error(err, fmt.Sprintf("删除cluster %s 从属集群%s 删除失败 ", c.Name, cluster.ClusterName))
			}
		},
	})
}

