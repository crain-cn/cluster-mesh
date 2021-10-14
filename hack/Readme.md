# createapi
<pre>
mkdir temp && cd /temp
1. 初始化 
  kubebuilder init --project-version="v1beta1" --domain cloud.mesh 
  kubebuilder edit --multigroup=true

2. createapi
  kubebuilder create api --group cloud --version v1beta1 --kind Cluster 
  kubebuilder create api --group cloud --version v1beta1 --kind ClusterMesh 
  
3. apis/types加上

namespaced
   // +genclient
cluster
   // +genclient
   // +genclient:nonNamespaced
   // +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object 
4.
./hack/update-codegen.sh

5. 去下面目录复制client代码
github.com/crain-cn/cluster-mesh
</pre>
