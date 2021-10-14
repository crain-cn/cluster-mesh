/*
Copyright 2019 The Kubernetes Authors.

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

package kube

import (
	"errors"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"regexp"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	diskcached "k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	defaultCacheDir = filepath.Join(homedir.HomeDir(), ".kube", "cache")
)

type ClusterClientGetter struct {
	config string
}

func NewClusterClientGetter(config string) *ClusterClientGetter {
	return &ClusterClientGetter{
		config: config,
	}
}

func (c *ClusterClientGetter) ToRESTConfig() (*rest.Config, error) {
	return LoadConfig(c.config)
}

func (c *ClusterClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	config, err := c.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	config.Burst = 100
	cacheDir := defaultCacheDir
	httpCacheDir := filepath.Join(cacheDir, "http")
	discoveryCacheDir := computeDiscoverCacheDir(filepath.Join(cacheDir, "discovery"), config.Host)
	return diskcached.NewCachedDiscoveryClientForConfig(config, discoveryCacheDir, httpCacheDir, time.Duration(10*time.Minute))
}

var overlyCautiousIllegalFileCharacters = regexp.MustCompile(`[^(\w/\.)]`)

func computeDiscoverCacheDir(parentDir, host string) string {
	// strip the optional scheme from host if its there:
	schemelessHost := strings.Replace(strings.Replace(host, "https://", "", 1), "http://", "", 1)
	// now do a simple collapse of non-AZ09 characters.  Collisions are possible but unlikely.  Even if we do collide the problem is short lived
	safeHost := overlyCautiousIllegalFileCharacters.ReplaceAllString(schemelessHost, "_")
	return filepath.Join(parentDir, safeHost)
}

func (c *ClusterClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	discoveryClient, err := c.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, discoveryClient)
	return expander, nil
}

func localConfig() (*rest.Config, error) {
	return rest.InClusterConfig()
}

func kubeConfigs(kubeconfig string) (map[string]rest.Config, string, error) {
	// Attempt to load external clusters too
	var loader clientcmd.ClientConfigLoader
	if kubeconfig != "" { // load from --kubeconfig
		loader = &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
	} else {
		loader = clientcmd.NewDefaultClientConfigLoadingRules()
	}

	cfg, err := loader.Load()
	if err != nil && kubeconfig != "" {
		return nil, "", fmt.Errorf("load: %v", err)
	}
	if err != nil {
		logrus.WithError(err).Warn("Cannot load kubecfg")
		return nil, "", nil
	}
	configs := map[string]rest.Config{}
	for context := range cfg.Contexts {
		contextCfg, err := clientcmd.NewNonInteractiveClientConfig(*cfg, context, &clientcmd.ConfigOverrides{}, loader).ClientConfig()
		if err != nil {
			return nil, "", fmt.Errorf("create %s client: %v", context, err)
		}
		configs[context] = *contextCfg
		logrus.Infof("Parsed kubeconfig context: %s", context)
	}
	return configs, cfg.CurrentContext, nil
}

func buildConfigs(buildCluster string) (map[string]rest.Config, error) {
	if buildCluster == "" { // load from --build-cluster
		return nil, nil
	}
	data, err := ioutil.ReadFile(buildCluster)
	if err != nil {
		return nil, fmt.Errorf("read: %v", err)
	}
	raw, err := UnmarshalClusterMap(data)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}
	cfg := &clientcmdapi.Config{
		Clusters:  map[string]*clientcmdapi.Cluster{},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{},
		Contexts:  map[string]*clientcmdapi.Context{},
	}
	for alias, config := range raw {
		cfg.Clusters[alias] = &clientcmdapi.Cluster{
			Server:                   config.Endpoint,
			CertificateAuthorityData: config.ClusterCACertificate,
		}
		cfg.AuthInfos[alias] = &clientcmdapi.AuthInfo{
			ClientCertificateData: config.ClientCertificate,
			ClientKeyData:         config.ClientKey,
		}
		cfg.Contexts[alias] = &clientcmdapi.Context{
			Cluster:  alias,
			AuthInfo: alias,
			// TODO(fejta): Namespace?
		}
	}
	configs := map[string]rest.Config{}
	for context := range cfg.Contexts {
		logrus.Infof("* %s", context)
		contextCfg, err := clientcmd.NewNonInteractiveClientConfig(*cfg, context, &clientcmd.ConfigOverrides{}, nil).ClientConfig()
		if err != nil {
			return nil, fmt.Errorf("create %s client: %v", context, err)
		}
		// An arbitrary high number we expect to not exceed. There are various components that need more than the default 5 QPS/10 Burst, e.G.
		// hook for creating ProwJobs and Plank for creating Pods.
		contextCfg.QPS = 100
		contextCfg.Burst = 1000
		configs[context] = *contextCfg
	}
	return configs, nil
}

func mergeConfigs(local *rest.Config, foreign map[string]rest.Config, currentContext string, buildClusters map[string]rest.Config) (map[string]rest.Config, error) {
	if buildClusters != nil {
		if _, ok := buildClusters[DefaultClusterAlias]; !ok {
			return nil, fmt.Errorf("build-cluster must have a %s context", DefaultClusterAlias)
		}
	}
	ret := map[string]rest.Config{}
	for ctx, cfg := range foreign {
		ret[ctx] = cfg
	}
	for ctx, cfg := range buildClusters {
		ret[ctx] = cfg
	}
	if local != nil {
		ret[InClusterContext] = *local
	} else if currentContext != "" {
		ret[InClusterContext] = ret[currentContext]
	} else {
		return nil, errors.New("no prow cluster access: in-cluster current kubecfg context required")
	}
	if len(ret) == 0 {
		return nil, errors.New("no client contexts found")
	}
	if _, ok := ret[DefaultClusterAlias]; !ok {
		ret[DefaultClusterAlias] = ret[InClusterContext]
	}
	return ret, nil
}

// LoadClusterConfigs loads rest.Configs for creation of clients, by using either a normal
// .kube/config file, a custom `Cluster` file, or both. The configs are returned in a mapping
// of context --> config. The default context is included in this mapping and specified as a
// return vaule. Errors are returned if .kube/config is specified and invalid or if no valid
// contexts are found.
func LoadClusterConfigs(kubeconfig, buildCluster string) (map[string]rest.Config, error) {

	logrus.Infof("Loading cluster contexts...")
	// This will work if we are running inside kubernetes
	localCfg, err := localConfig()
	if err != nil {
		logrus.WithError(err).Warn("Could not create in-cluster config (expected when running outside the cluster).")
	}

	kubeCfgs, currentContext, err := kubeConfigs(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("kubecfg: %v", err)
	}

	// TODO(fejta): drop build-cluster support
	buildCfgs, err := buildConfigs(buildCluster)
	if err != nil {
		return nil, fmt.Errorf("build-cluster: %v", err)
	}

	return mergeConfigs(localCfg, kubeCfgs, currentContext, buildCfgs)
}

func LoadConfig(kubeconfig string) (*rest.Config, error) {
	cfg, err := clientcmd.Load([]byte(kubeconfig))
	if err != nil {
		logrus.WithError(err).Fatalf("unmarshal: %v", err)
	}
	for context := range cfg.Contexts {
		contextCfg, err := clientcmd.NewNonInteractiveClientConfig(*cfg, context, &clientcmd.ConfigOverrides{}, nil).ClientConfig()
		if err != nil {
			logrus.WithError(err).Fatalf("create %s client: %v", context, err)
		}
		// An arbitrary high number we expect to not exceed. There are various components that need more than the default 5 QPS/10 Burst, e.G.
		// hook for creating ProwJobs and Plank for creating Pods.
		contextCfg.QPS = 100
		contextCfg.Burst = 1000
		return contextCfg, nil
	}
	return nil, errors.New("invalid kubeconfig")
}

func NewObjectClient(config *rest.Config) (client.Client, error) {
	c, err := client.New(config, client.Options{})
	if err != nil {
		return nil, err
	}
	ca, err := cache.New(config, cache.Options{})
	if err != nil {
		return nil, err
	}
	return &client.DelegatingClient{
		Reader: &client.DelegatingReader{
			CacheReader:  ca,
			ClientReader: c,
		},
		Writer:       c,
		StatusClient: c,
	}, nil
}
