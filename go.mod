module github.com/crain-cn/cluster-mesh

go 1.15

require (
	github.com/go-logr/logr v0.2.0
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/sirupsen/logrus v1.6.0
	k8s.io/api v0.20.0
	k8s.io/apimachinery v0.20.0
	k8s.io/cli-runtime v0.22.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/code-generator v0.22.2
	sigs.k8s.io/controller-runtime v0.4.0
	sigs.k8s.io/yaml v1.2.0
)

replace (
	k8s.io/api => k8s.io/api v0.20.0
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.20.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.0
	k8s.io/apiserver => k8s.io/apiserver v0.20.0
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.20.0
	k8s.io/client-go => k8s.io/client-go v0.20.0
	k8s.io/code-generator => k8s.io/code-generator v0.20.0
)
