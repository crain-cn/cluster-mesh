/*


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

package cluster

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cloudmeshv1beta1 "github.com/crain-cn/cluster-mesh/api/cloud.mesh/v1beta1"
	"github.com/crain-cn/cluster-mesh/pkg/util/kube"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// ClusterReconciler reconciles a ClusterEnv object
type ClusterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cloud.mesh,resources=Cluster,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cloud.mesh,resources=Cluster/status,verbs=get;update;patch

func (r *ClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	cluster := &cloudmeshv1beta1.Cluster{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, cluster)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			r.Log.Info(fmt.Sprintf("Cluster resource %s not found. Ignoring since object must be deleted.", req.NamespacedName))
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		r.Log.Error(err, "Failed to getClusterMesh.")
		return reconcile.Result{}, err
	}
	// TODO 如果初始化过程中程序重启不会继续初始化，需要完善这个代码
	if !cluster.Spec.InitComplete {
		err := r.initCluster(cluster)
		if err != nil {
			cluster.Status.Message = err.Error()
			cluster.Status.Stage = cloudmeshv1beta1.INITFAILED
			cluster.Spec.InitComplete = true
		} else {
			cluster.Status.Stage = cloudmeshv1beta1.INITCOMPLETE
			cluster.Spec.InitComplete = true
		}
		return reconcile.Result{}, r.Update(context.Background(), cluster)
	}
	return reconcile.Result{}, nil
}

func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudmeshv1beta1.Cluster{}).
		Complete(r)
}

func (r *ClusterReconciler) initCluster(cluster *cloudmeshv1beta1.Cluster) error {
	objects, err := kube.ParseYamlToObject(cluster.Spec.KubeConfig, []string{"/etc/cluster"})
	if err != nil {
		r.Log.Error(err, "Failed to parse yaml ")
		return err
	}
	clusterConfig, err := kube.LoadConfig(cluster.Spec.KubeConfig)
	if err != nil {
		r.Log.Error(err, "Failed to loadConfig  Cluster:%v", cluster.Name)
		return err
	}
	oClient, err := kube.NewObjectClient(clusterConfig)
	if err != nil {
		r.Log.Error(err, "Failed to get Client")
		return err
	}
	err = oClient.Create(context.TODO(), objects)
	if err != nil {
		r.Log.Error(err, "Failed to create resource")
		return err
	}
	return nil
}

func Add(mgr manager.Manager) error {
	return add(mgr, &ClusterReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log.WithName("controllers").WithName("Cluster"),
	})
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("Cluster-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource
	err = c.Watch(&source.Kind{Type: &cloudmeshv1beta1.Cluster{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	return nil
}
