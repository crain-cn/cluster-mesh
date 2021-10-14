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

package v1beta1

import (
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ClusterMeshSpec defines the desired state of ClusterMesh
type ClusterMeshSpec struct {
	Env         AppEnv       `json:"env"`
	Deployment  string       `json:"deployment"`
	Service     string       `json:"service"`
	Ingress     string       `json:"ingress"`
	ClusterMode ClusterMode  `json:"clusterMode"`
	ConfigMap   []string     `json:"configMap"`
	MainCluster AppCluster   `json:"mainCluster"`
	Clusters    []AppCluster `json:"clusters"`
	FlowMode    flowMode     `json:"flowMode"`
}

type ClusterMode string

const (
	MainClusterMode   ClusterMode = "main"
	ExtendClusterMode ClusterMode = "extend"
	BackupClusterMode ClusterMode = "backup"
)

type flowMode string

const (
	GatewayFlowMode flowMode = "gateway"
	IngressSLB      flowMode = "ingress-slb"
	ServiceSLB      flowMode = "service-slb"
)

type clusterMode string

const (
	// 应用部署的主集群 仅有一个
	MainCluster clusterMode = "main"
	// 应用部署的扩展集群，任意多个
	ExtendCluster clusterMode = "extend"
	// 应用部署的自动切换型冷备集群 至多一个
	AutoBackupCluster clusterMode = "autoBackup"
	// 应用部署的冷备集群 任意多个
	BackupCluster clusterMode = "backup"
	// double live
	DoubleAliveCluster clusterMode = "doubleAlive"
)

const (
	ECILogVolumeName = "log-clear"
	ECILogClearName  = "eci-log-clear"
	ECILogVolumePath = "/data/logclear"
)

type AppEnv string

const (
	DevEnv    AppEnv = "dev"
	TestEnv   AppEnv = "test"
	GrayEnv   AppEnv = "gray"
	OnlineEnv AppEnv = "online"

	AliyunEnv  AppEnv = "aliyun"
	IdcEnv     AppEnv = "idc"
	TencentEnv AppEnv = "tencent"
)

type AppCluster struct {
	ClusterName   string      `json:"clusterName"`
	ClusterMode   clusterMode `json:"clusterMode"`
	ClusterAppEnv AppEnv      `json:"appEnv"`
	// 是否开启与主集群同步升级
	UpgradeSyncEnable bool `json:"upgradeSyncEnable"`
	// 是否开启与主集群配置同步
	ConfigSyncEnable bool `json:"configSyncEnable"`
	// 流量比例(对灾备集群无效)
	RatioFlow int `json:"ratioFlow"`
	// 集群删除标记位(对主集群无效)
	Deleted bool `json:"deleted"`
}

type DeletedStatus struct {
	OfflineStatus bool   `json:"offlineStatus"`
	Message       string `json:"message"`
}

// ClusterMeshStatus defines the observed state of ClusterMesh
type ClusterMeshStatus struct {
	// 流量分布 只记录非主集群集群流量占比 主集群通过计算得出
	FlowDistribution map[string]int `json:"flowDistribution"`
	// 集群下应用状态
	DeploymentStatus map[string]v1.DeploymentStatus `json:"deploymentStatus"`
	// 集群下应用版本信息
	ImageStatus map[string][]string `json:"imageStatus"`
	// 应用下线状态
	DeletedStatus map[string]*DeletedStatus `json:"deleteStatus"`

	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes        []string `json:"nodes"`
	BackupStatus []string `json:"backupStatus`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +genclient
// ClusterMesh is the Schema for the ClusterMeshs API
type ClusterMesh struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterMeshSpec   `json:"spec,omitempty"`
	Status ClusterMeshStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterMeshList contains a list of ClusterMesh
type ClusterMeshList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterMesh `json:"items"`
}
