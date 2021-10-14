package v1beta1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GatewayConfig struct {
	GatewaySign string `json:"gatewaySign"`
	BaseURL     string `json:"baseURL"`
}

type GatewayCluster struct {
	Id          int                `json:"id"`
	ClusterType GatewayClusterType `json:"type"`
	Description string             `json:"description"`
}

type GateWayClusterId struct {
	RootId       int                `json:"rootId"`
	LeafId       int                `json:"leafId"`
	ClusterType  GatewayClusterType `json:"type"`
	DoubleLeafId int                `json:"doubleLeafId"`
}

const (
	HAGateway  GatewayClusterType = "HA"
	APIGateway GatewayClusterType = "API"
)

type GatewayClusterType string

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	//集群所在网段
	NetworkSegment string `json:"network_segment"`
	//集群出口ip
	ExportIp string `json:"export_ip"`
	// 集群名
	Name string `json:"name"`
	// kubectl alisa
	Alias string `json:"alias"`
	// 集群描述
	Description string `json:"description"`
	// 集群apiserver
	ApiServer string `json:"apiServer"`
	// 集群kubeconfig
	KubeConfig string `json:"kubeConfig"`
	// 集群环境(test/dev/gray/online/aliyun/aliyun-gray)
	Environment string `json:"environment"`
	// 集群允许项目发布的目标环境
	SupportEnv map[AppEnv]ClusterEnv `json:"supportEnv"`
	// 是否为CI集群(构建任务执行集群)
	CiCluster bool `json:"ciCluster"`
	// 集群类型
	ClusterType ClusterType `json:"clusterType"`
	// 集群默认镜像仓库域名
	RegistryDomain string `json:"registryDomain"`
	// 集群前端页面提示
	Notice string `json:"notice"`
	// 集群初始化完成
	InitComplete bool `json:"initComplete"`
	// 该集群是否支持SLB
	SupportSLB bool `json:"supportSLB"`
	// Service网段
	ServiceSubnet string `json:"serviceSubnet"`
	// 是否存在监控
	SupportMonitor bool `json:"supportMonitor"`
	// 是否开启从主集群同步数据
	DataSyncEnable bool `json:"dataSyncEnable"`
	// 集群保密信息
	Secret map[string]string `json:"secret"`
	// 是否为主集群(元数据集群)
	MainCluster bool `json:"mainCluster"`
}

// K8S集群类型
type ClusterType string

const (
	IDC        ClusterType = "std-idc"
	AliyunACK  ClusterType = "ali-ack"
	AliyunASK  ClusterType = "ali-ask"
	TencentTke ClusterType = "ten-tke"
	TencentEks ClusterType = "ten-eks"
)

const (
	EciLogProjectKey      = "aliyun_logs_project"
	EciLogMachineGroupKey = "aliyun_logs_machinegroup"
	EciLogPrefix          = "aliyun_logs_"
	TopicPrefix           = "log_type"
	EciLogValuePrefix     = "wxy-k8s-eci-"
)

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes        []string `json:"nodes"`
	BackupStatus []string `json:"backup_status`
	// 集群状态信息
	Stage       clusterInitStage `json:"stage"`
	Message     string           `json:"message"`
	Connectable bool             `json:"connectable"`
}

type clusterInitStage string

const (
	UNINIT       clusterInitStage = "集群尚未开始初始化"
	INITING      clusterInitStage = "集群正在初始化"
	INITCOMPLETE clusterInitStage = "集群初始化完成"
	INITFAILED   clusterInitStage = "集群初始化失败"
)

type ClusterEnv struct {
	// 集群应用发布默认节点组
	DefaultNodeGroup string `json:"defaultNodeGroup"`
	// 上报网关ip类型
	ReportGateway ReportGatewayType `json:"reportGateway"`
	// 默认资源限制
	DefaultResourceRequirements v1.ResourceRequirements `json:"defaultResourceRequirements"`
	// 默认Ingress
	DefaultIngressController string `json:"defaultIngressController"`
	// 切流所需权限
	MinPermission ClusterPermission `json:"minPermission"`
	// 该环境是否支持操作网关
	SupportGateway bool `json:"supportGateway"`
	// 网关配置
	GatewayConfig GatewayConfig `json:"gatewayConfig"`
	// 支持的网关集群
	GatewayClusters []GatewayCluster `json:"gatewayClusters"`
	// DNS策略
	DNSPolicy v1.DNSPolicy `json:"dnsPolicy"`
	// DNS配置
	DNSConfig *v1.PodDNSConfig `json:"dnsConfig"`
	// ExtendDNS配置
	ExtendDnsConfig *v1.PodDNSConfig `json:"extendDnsConfig"`
	//双活 集群配置
	DoubleLiveGateway []GateWayClusterId `json:"doubleLiveGateway"`
}

type ClusterPermission string

const (
	ClusterAdmin       ClusterPermission = "kube_admin"
	ClusterDeveloper   ClusterPermission = "kube_developer"
	ClusterViewer      ClusterPermission = "kube_viewer"
	NamespaceOwner     ClusterPermission = "namespace_owner"
	NamespaceAdmin     ClusterPermission = "namespace_admin"
	NamespaceDeveloper ClusterPermission = "namespace_developer"
	NamespaceViewer    ClusterPermission = "namespace_viewer"
)

// 上报给网关的入口ip类型
type ReportGatewayType string

const (
	ReportGatewayIngressEndpoint     ReportGatewayType = "ingress-endpoint"
	ReportGatewayIngressLoadBalancer ReportGatewayType = "ingress-lb"
	ReportGatewayServiceEndpoint     ReportGatewayType = "service-endpoint"
	ReportGatewayServiceLoadBalancer ReportGatewayType = "service-lb"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// Cluster is the Schema for the Clusters API
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterList contains a list of Cluster
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}
