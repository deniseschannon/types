package v3

import (
	"fmt"
	"strings"

	"github.com/rancher/types/image"
)

const (
	DefaultK8s = "v1.11.1-rancher1-1"
)

var (
	m = image.Mirror

	// k8sVersionsCurrent are the latest versions available for installation
	k8sVersionsCurrent = []string{
		"v1.9.7-rancher2-2",
		"v1.10.5-rancher1-2",
		"v1.11.1-rancher1-1",
	}

	// K8sVersionToRKESystemImages is dynamically populated on init() with the latest versions
	K8sVersionToRKESystemImages map[string]RKESystemImages

	// K8sVersionServiceOptions - service options per k8s version
	K8sVersionServiceOptions = map[string]KubernetesServicesOptions{
		"v1.10": {
			KubeAPI: map[string]string{
				"tls-cipher-suites":        "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
				"endpoint-reconciler-type": "lease",
			},
			Kubelet: map[string]string{
				"tls-cipher-suites": "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
			},
		},
		"v1.9": {
			KubeAPI: map[string]string{
				"endpoint-reconciler-type": "lease",
			},
		},
	}

	// ToolsSystemImages default images for alert, pipeline, logging
	ToolsSystemImages = struct {
		AlertSystemImages    AlertSystemImages
		PipelineSystemImages PipelineSystemImages
		LoggingSystemImages  LoggingSystemImages
	}{
		AlertSystemImages: AlertSystemImages{
			AlertManager:       m("prom/alertmanager:v0.11.0"),
			AlertManagerHelper: m("rancher/alertmanager-helper:v0.0.2"),
		},
		PipelineSystemImages: PipelineSystemImages{
			Jenkins:       m("jenkins/jenkins:2.107-slim"),
			JenkinsJnlp:   m("jenkins/jnlp-slave:3.10-1-alpine"),
			AlpineGit:     m("alpine/git:1.0.4"),
			PluginsDocker: m("plugins/docker:17.12"),
		},
		LoggingSystemImages: LoggingSystemImages{
			Fluentd:                       m("rancher/fluentd:v0.1.10"),
			FluentdHelper:                 m("rancher/fluentd-helper:v0.1.2"),
			LogAggregatorFlexVolumeDriver: m("rancher/log-aggregator:v0.1.3"),
			Elaticsearch:                  m("quay.io/pires/docker-elasticsearch-kubernetes:5.6.2"),
			Kibana:                        m("kibana:5.6.4"),
			Busybox:                       m("busybox"),
		},
	}

	AllK8sVersions = map[string]RKESystemImages{
		"v1.8.10-rancher1-1": {
			Etcd:                      m("quay.io/coreos/etcd:v3.0.17"),
			Kubernetes:                m("rancher/hyperkube:v1.8.10-rancher2"),
			Alpine:                    m("rancher/rke-tools:v0.1.4"),
			NginxProxy:                m("rancher/rke-tools:v0.1.4"),
			CertDownloader:            m("rancher/rke-tools:v0.1.4"),
			KubernetesServicesSidecar: m("rancher/rke-tools:v0.1.4"),
			KubeDNS:                   m("gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.5"),
			DNSmasq:                   m("gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.5"),
			KubeDNSSidecar:            m("gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.5"),
			KubeDNSAutoscaler:         m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
			Flannel:                   m("quay.io/coreos/flannel:v0.9.1"),
			FlannelCNI:                m("quay.io/coreos/flannel-cni:v0.2.0"),
			CalicoNode:                m("quay.io/calico/node:v3.1.1"),
			CalicoCNI:                 m("quay.io/calico/cni:v3.1.1"),
			CalicoCtl:                 m("quay.io/calico/ctl:v2.0.0"),
			CanalNode:                 m("quay.io/calico/node:v3.1.1"),
			CanalCNI:                  m("quay.io/calico/cni:v3.1.1"),
			CanalFlannel:              m("quay.io/coreos/flannel:v0.9.1"),
			WeaveNode:                 m("weaveworks/weave-kube:2.1.2"),
			WeaveCNI:                  m("weaveworks/weave-npc:2.1.2"),
			PodInfraContainer:         m("gcr.io/google_containers/pause-amd64:3.0"),
			Ingress:                   m("rancher/nginx-ingress-controller:0.10.2-rancher3"),
			IngressBackend:            m("k8s.gcr.io/defaultbackend:1.4"),
		},
		"v1.8.11-rancher1": {
			Etcd:                      m("quay.io/coreos/etcd:v3.0.17"),
			Kubernetes:                m("rancher/hyperkube:v1.8.11-rancher2"),
			Alpine:                    m("rancher/rke-tools:v0.1.4"),
			NginxProxy:                m("rancher/rke-tools:v0.1.4"),
			CertDownloader:            m("rancher/rke-tools:v0.1.4"),
			KubernetesServicesSidecar: m("rancher/rke-tools:v0.1.4"),
			KubeDNS:                   m("gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.5"),
			DNSmasq:                   m("gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.5"),
			KubeDNSSidecar:            m("gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.5"),
			KubeDNSAutoscaler:         m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
			Flannel:                   m("quay.io/coreos/flannel:v0.9.1"),
			FlannelCNI:                m("quay.io/coreos/flannel-cni:v0.2.0"),
			CalicoNode:                m("quay.io/calico/node:v3.1.1"),
			CalicoCNI:                 m("quay.io/calico/cni:v3.1.1"),
			CalicoCtl:                 m("quay.io/calico/ctl:v2.0.0"),
			CanalNode:                 m("quay.io/calico/node:v3.1.1"),
			CanalCNI:                  m("quay.io/calico/cni:v3.1.1"),
			CanalFlannel:              m("quay.io/coreos/flannel:v0.9.1"),
			WeaveNode:                 m("weaveworks/weave-kube:2.1.2"),
			WeaveCNI:                  m("weaveworks/weave-npc:2.1.2"),
			PodInfraContainer:         m("gcr.io/google_containers/pause-amd64:3.0"),
			Ingress:                   m("rancher/nginx-ingress-controller:0.10.2-rancher3"),
			IngressBackend:            m("k8s.gcr.io/defaultbackend:1.4"),
		},
		"v1.8.11-rancher2-1": {
			Etcd:                      m("quay.io/coreos/etcd:v3.0.17"),
			Kubernetes:                m("rancher/hyperkube:v1.8.11-rancher2"),
			Alpine:                    m("rancher/rke-tools:v0.1.8"),
			NginxProxy:                m("rancher/rke-tools:v0.1.8"),
			CertDownloader:            m("rancher/rke-tools:v0.1.8"),
			KubernetesServicesSidecar: m("rancher/rke-tools:v0.1.8"),
			KubeDNS:                   m("gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.5"),
			DNSmasq:                   m("gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.5"),
			KubeDNSSidecar:            m("gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.5"),
			KubeDNSAutoscaler:         m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
			Flannel:                   m("quay.io/coreos/flannel:v0.9.1"),
			FlannelCNI:                m("quay.io/coreos/flannel-cni:v0.2.0"),
			CalicoNode:                m("quay.io/calico/node:v3.1.1"),
			CalicoCNI:                 m("quay.io/calico/cni:v3.1.1"),
			CalicoCtl:                 m("quay.io/calico/ctl:v2.0.0"),
			CanalNode:                 m("quay.io/calico/node:v3.1.1"),
			CanalCNI:                  m("quay.io/calico/cni:v3.1.1"),
			CanalFlannel:              m("quay.io/coreos/flannel:v0.9.1"),
			WeaveNode:                 m("weaveworks/weave-kube:2.1.2"),
			WeaveCNI:                  m("weaveworks/weave-npc:2.1.2"),
			PodInfraContainer:         m("gcr.io/google_containers/pause-amd64:3.0"),
			Ingress:                   m("rancher/nginx-ingress-controller:0.10.2-rancher3"),
			IngressBackend:            m("k8s.gcr.io/defaultbackend:1.4"),
		},
		"v1.9.5-rancher1-1": {
			Etcd:                      m("quay.io/coreos/etcd:v3.1.12"),
			Kubernetes:                m("rancher/hyperkube:v1.9.5-rancher1"),
			Alpine:                    m("rancher/rke-tools:v0.1.4"),
			NginxProxy:                m("rancher/rke-tools:v0.1.4"),
			CertDownloader:            m("rancher/rke-tools:v0.1.4"),
			KubernetesServicesSidecar: m("rancher/rke-tools:v0.1.4"),
			KubeDNS:                   m("gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.7"),
			DNSmasq:                   m("gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.7"),
			KubeDNSSidecar:            m("gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.7"),
			KubeDNSAutoscaler:         m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
			Flannel:                   m("quay.io/coreos/flannel:v0.9.1"),
			FlannelCNI:                m("quay.io/coreos/flannel-cni:v0.2.0"),
			CalicoNode:                m("quay.io/calico/node:v3.1.1"),
			CalicoCNI:                 m("quay.io/calico/cni:v3.1.1"),
			CalicoCtl:                 m("quay.io/calico/ctl:v2.0.0"),
			CanalNode:                 m("quay.io/calico/node:v3.1.1"),
			CanalCNI:                  m("quay.io/calico/cni:v3.1.1"),
			CanalFlannel:              m("quay.io/coreos/flannel:v0.9.1"),
			WeaveNode:                 m("weaveworks/weave-kube:2.1.2"),
			WeaveCNI:                  m("weaveworks/weave-npc:2.1.2"),
			PodInfraContainer:         m("gcr.io/google_containers/pause-amd64:3.0"),
			Ingress:                   m("rancher/nginx-ingress-controller:0.10.2-rancher3"),
			IngressBackend:            m("k8s.gcr.io/defaultbackend:1.4"),
		},
		"v1.9.7-rancher1": {
			Etcd:                      m("quay.io/coreos/etcd:v3.1.12"),
			Kubernetes:                m("rancher/hyperkube:v1.9.7-rancher2"),
			Alpine:                    m("rancher/rke-tools:v0.1.4"),
			NginxProxy:                m("rancher/rke-tools:v0.1.4"),
			CertDownloader:            m("rancher/rke-tools:v0.1.4"),
			KubernetesServicesSidecar: m("rancher/rke-tools:v0.1.4"),
			KubeDNS:                   m("gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.7"),
			DNSmasq:                   m("gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.7"),
			KubeDNSSidecar:            m("gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.7"),
			KubeDNSAutoscaler:         m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
			Flannel:                   m("quay.io/coreos/flannel:v0.9.1"),
			FlannelCNI:                m("quay.io/coreos/flannel-cni:v0.2.0"),
			CalicoNode:                m("quay.io/calico/node:v3.1.1"),
			CalicoCNI:                 m("quay.io/calico/cni:v3.1.1"),
			CalicoCtl:                 m("quay.io/calico/ctl:v2.0.0"),
			CanalNode:                 m("quay.io/calico/node:v3.1.1"),
			CanalCNI:                  m("quay.io/calico/cni:v3.1.1"),
			CanalFlannel:              m("quay.io/coreos/flannel:v0.9.1"),
			WeaveNode:                 m("weaveworks/weave-kube:2.1.2"),
			WeaveCNI:                  m("weaveworks/weave-npc:2.1.2"),
			PodInfraContainer:         m("gcr.io/google_containers/pause-amd64:3.0"),
			Ingress:                   m("rancher/nginx-ingress-controller:0.10.2-rancher3"),
			IngressBackend:            m("k8s.gcr.io/defaultbackend:1.4"),
		},
		"v1.9.7-rancher2-1": {
			Etcd:                      m("quay.io/coreos/etcd:v3.1.12"),
			Kubernetes:                m("rancher/hyperkube:v1.9.7-rancher2"),
			Alpine:                    m("rancher/rke-tools:v0.1.8"),
			NginxProxy:                m("rancher/rke-tools:v0.1.8"),
			CertDownloader:            m("rancher/rke-tools:v0.1.8"),
			KubernetesServicesSidecar: m("rancher/rke-tools:v0.1.8"),
			KubeDNS:                   m("gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.7"),
			DNSmasq:                   m("gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.7"),
			KubeDNSSidecar:            m("gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.7"),
			KubeDNSAutoscaler:         m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
			Flannel:                   m("quay.io/coreos/flannel:v0.9.1"),
			FlannelCNI:                m("quay.io/coreos/flannel-cni:v0.2.0"),
			CalicoNode:                m("quay.io/calico/node:v3.1.1"),
			CalicoCNI:                 m("quay.io/calico/cni:v3.1.1"),
			CalicoCtl:                 m("quay.io/calico/ctl:v2.0.0"),
			CanalNode:                 m("quay.io/calico/node:v3.1.1"),
			CanalCNI:                  m("quay.io/calico/cni:v3.1.1"),
			CanalFlannel:              m("quay.io/coreos/flannel:v0.9.1"),
			WeaveNode:                 m("weaveworks/weave-kube:2.1.2"),
			WeaveCNI:                  m("weaveworks/weave-npc:2.1.2"),
			PodInfraContainer:         m("gcr.io/google_containers/pause-amd64:3.0"),
			Ingress:                   m("rancher/nginx-ingress-controller:0.10.2-rancher3"),
			IngressBackend:            m("k8s.gcr.io/defaultbackend:1.4"),
		},
		"v1.9.7-rancher2-2": {
			Etcd:                      m("quay.io/coreos/etcd:v3.1.12"),
			Kubernetes:                m("rancher/hyperkube:v1.9.7-rancher2"),
			Alpine:                    m("rancher/rke-tools:v0.1.13"),
			NginxProxy:                m("rancher/rke-tools:v0.1.13"),
			CertDownloader:            m("rancher/rke-tools:v0.1.13"),
			KubernetesServicesSidecar: m("rancher/rke-tools:v0.1.13"),
			KubeDNS:                   m("gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.7"),
			DNSmasq:                   m("gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.7"),
			KubeDNSSidecar:            m("gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.7"),
			KubeDNSAutoscaler:         m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
			Flannel:                   m("quay.io/coreos/flannel:v0.9.1"),
			FlannelCNI:                m("quay.io/coreos/flannel-cni:v0.2.0"),
			CalicoNode:                m("quay.io/calico/node:v3.1.1"),
			CalicoCNI:                 m("quay.io/calico/cni:v3.1.1"),
			CalicoCtl:                 m("quay.io/calico/ctl:v2.0.0"),
			CanalNode:                 m("quay.io/calico/node:v3.1.1"),
			CanalCNI:                  m("quay.io/calico/cni:v3.1.1"),
			CanalFlannel:              m("quay.io/coreos/flannel:v0.9.1"),
			WeaveNode:                 m("weaveworks/weave-kube:2.1.2"),
			WeaveCNI:                  m("weaveworks/weave-npc:2.1.2"),
			PodInfraContainer:         m("gcr.io/google_containers/pause-amd64:3.0"),
			Ingress:                   m("rancher/nginx-ingress-controller:0.16.2-rancher1"),
			IngressBackend:            m("k8s.gcr.io/defaultbackend:1.4"),
			MetricsServer:             m("gcr.io/google_containers/metrics-server-amd64:v0.2.1"),
		},
		"v1.10.0-rancher1-1": {
			Etcd:                      m("quay.io/coreos/etcd:v3.1.12"),
			Kubernetes:                m("rancher/hyperkube:v1.10.0-rancher1"),
			Alpine:                    m("rancher/rke-tools:v0.1.4"),
			NginxProxy:                m("rancher/rke-tools:v0.1.4"),
			CertDownloader:            m("rancher/rke-tools:v0.1.4"),
			KubernetesServicesSidecar: m("rancher/rke-tools:v0.1.4"),
			KubeDNS:                   m("gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.8"),
			DNSmasq:                   m("gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.8"),
			KubeDNSSidecar:            m("gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.8"),
			KubeDNSAutoscaler:         m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
			Flannel:                   m("quay.io/coreos/flannel:v0.9.1"),
			FlannelCNI:                m("quay.io/coreos/flannel-cni:v0.2.0"),
			CalicoNode:                m("quay.io/calico/node:v3.1.1"),
			CalicoCNI:                 m("quay.io/calico/cni:v3.1.1"),
			CalicoCtl:                 m("quay.io/calico/ctl:v2.0.0"),
			CanalNode:                 m("quay.io/calico/node:v3.1.1"),
			CanalCNI:                  m("quay.io/calico/cni:v3.1.1"),
			CanalFlannel:              m("quay.io/coreos/flannel:v0.9.1"),
			WeaveNode:                 m("weaveworks/weave-kube:2.1.2"),
			WeaveCNI:                  m("weaveworks/weave-npc:2.1.2"),
			PodInfraContainer:         m("gcr.io/google_containers/pause-amd64:3.1"),
			Ingress:                   m("rancher/nginx-ingress-controller:0.10.2-rancher3"),
			IngressBackend:            m("k8s.gcr.io/defaultbackend:1.4"),
		},
		"v1.10.1-rancher1": {
			Etcd:                      m("quay.io/coreos/etcd:v3.1.12"),
			Kubernetes:                m("rancher/hyperkube:v1.10.1-rancher2"),
			Alpine:                    m("rancher/rke-tools:v0.1.4"),
			NginxProxy:                m("rancher/rke-tools:v0.1.4"),
			CertDownloader:            m("rancher/rke-tools:v0.1.4"),
			KubernetesServicesSidecar: m("rancher/rke-tools:v0.1.4"),
			KubeDNS:                   m("gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.8"),
			DNSmasq:                   m("gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.8"),
			KubeDNSSidecar:            m("gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.8"),
			KubeDNSAutoscaler:         m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
			Flannel:                   m("quay.io/coreos/flannel:v0.9.1"),
			FlannelCNI:                m("quay.io/coreos/flannel-cni:v0.2.0"),
			CalicoNode:                m("quay.io/calico/node:v3.1.1"),
			CalicoCNI:                 m("quay.io/calico/cni:v3.1.1"),
			CalicoCtl:                 m("quay.io/calico/ctl:v2.0.0"),
			CanalNode:                 m("quay.io/calico/node:v3.1.1"),
			CanalCNI:                  m("quay.io/calico/cni:v3.1.1"),
			CanalFlannel:              m("quay.io/coreos/flannel:v0.9.1"),
			WeaveNode:                 m("weaveworks/weave-kube:2.1.2"),
			WeaveCNI:                  m("weaveworks/weave-npc:2.1.2"),
			PodInfraContainer:         m("gcr.io/google_containers/pause-amd64:3.1"),
			Ingress:                   m("rancher/nginx-ingress-controller:0.10.2-rancher3"),
			IngressBackend:            m("k8s.gcr.io/defaultbackend:1.4"),
		},
		"v1.10.1-rancher2-1": {
			Etcd:                      m("quay.io/coreos/etcd:v3.1.12"),
			Kubernetes:                m("rancher/hyperkube:v1.10.1-rancher2"),
			Alpine:                    m("rancher/rke-tools:v0.1.8"),
			NginxProxy:                m("rancher/rke-tools:v0.1.8"),
			CertDownloader:            m("rancher/rke-tools:v0.1.8"),
			KubernetesServicesSidecar: m("rancher/rke-tools:v0.1.8"),
			KubeDNS:                   m("gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.8"),
			DNSmasq:                   m("gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.8"),
			KubeDNSSidecar:            m("gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.8"),
			KubeDNSAutoscaler:         m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
			Flannel:                   m("quay.io/coreos/flannel:v0.9.1"),
			FlannelCNI:                m("quay.io/coreos/flannel-cni:v0.2.0"),
			CalicoNode:                m("quay.io/calico/node:v3.1.1"),
			CalicoCNI:                 m("quay.io/calico/cni:v3.1.1"),
			CalicoCtl:                 m("quay.io/calico/ctl:v2.0.0"),
			CanalNode:                 m("quay.io/calico/node:v3.1.1"),
			CanalCNI:                  m("quay.io/calico/cni:v3.1.1"),
			CanalFlannel:              m("quay.io/coreos/flannel:v0.9.1"),
			WeaveNode:                 m("weaveworks/weave-kube:2.1.2"),
			WeaveCNI:                  m("weaveworks/weave-npc:2.1.2"),
			PodInfraContainer:         m("gcr.io/google_containers/pause-amd64:3.1"),
			Ingress:                   m("rancher/nginx-ingress-controller:0.10.2-rancher3"),
			IngressBackend:            m("k8s.gcr.io/defaultbackend:1.4"),
		},
		"v1.10.3-rancher2-1": {
			Etcd:                      m("quay.io/coreos/etcd:v3.1.12"),
			Kubernetes:                m("rancher/hyperkube:v1.10.3-rancher2"),
			Alpine:                    m("rancher/rke-tools:v0.1.10"),
			NginxProxy:                m("rancher/rke-tools:v0.1.10"),
			CertDownloader:            m("rancher/rke-tools:v0.1.10"),
			KubernetesServicesSidecar: m("rancher/rke-tools:v0.1.10"),
			KubeDNS:                   m("gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.8"),
			DNSmasq:                   m("gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.8"),
			KubeDNSSidecar:            m("gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.8"),
			KubeDNSAutoscaler:         m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
			Flannel:                   m("quay.io/coreos/flannel:v0.9.1"),
			FlannelCNI:                m("quay.io/coreos/flannel-cni:v0.2.0"),
			CalicoNode:                m("quay.io/calico/node:v3.1.1"),
			CalicoCNI:                 m("quay.io/calico/cni:v3.1.1"),
			CalicoCtl:                 m("quay.io/calico/ctl:v2.0.0"),
			CanalNode:                 m("quay.io/calico/node:v3.1.1"),
			CanalCNI:                  m("quay.io/calico/cni:v3.1.1"),
			CanalFlannel:              m("quay.io/coreos/flannel:v0.9.1"),
			WeaveNode:                 m("weaveworks/weave-kube:2.1.2"),
			WeaveCNI:                  m("weaveworks/weave-npc:2.1.2"),
			PodInfraContainer:         m("gcr.io/google_containers/pause-amd64:3.1"),
			Ingress:                   m("rancher/nginx-ingress-controller:0.10.2-rancher3"),
			IngressBackend:            m("k8s.gcr.io/defaultbackend:1.4"),
		},
		"v1.10.5-rancher1-1": {
			Etcd:                      m("quay.io/coreos/etcd:v3.1.12"),
			Kubernetes:                m("rancher/hyperkube:v1.10.5-rancher1"),
			Alpine:                    m("rancher/rke-tools:v0.1.10"),
			NginxProxy:                m("rancher/rke-tools:v0.1.10"),
			CertDownloader:            m("rancher/rke-tools:v0.1.10"),
			KubernetesServicesSidecar: m("rancher/rke-tools:v0.1.10"),
			KubeDNS:                   m("gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.8"),
			DNSmasq:                   m("gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.8"),
			KubeDNSSidecar:            m("gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.8"),
			KubeDNSAutoscaler:         m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
			Flannel:                   m("quay.io/coreos/flannel:v0.9.1"),
			FlannelCNI:                m("quay.io/coreos/flannel-cni:v0.2.0"),
			CalicoNode:                m("quay.io/calico/node:v3.1.1"),
			CalicoCNI:                 m("quay.io/calico/cni:v3.1.1"),
			CalicoCtl:                 m("quay.io/calico/ctl:v2.0.0"),
			CanalNode:                 m("quay.io/calico/node:v3.1.1"),
			CanalCNI:                  m("quay.io/calico/cni:v3.1.1"),
			CanalFlannel:              m("quay.io/coreos/flannel:v0.9.1"),
			WeaveNode:                 m("weaveworks/weave-kube:2.1.2"),
			WeaveCNI:                  m("weaveworks/weave-npc:2.1.2"),
			PodInfraContainer:         m("gcr.io/google_containers/pause-amd64:3.1"),
			Ingress:                   m("rancher/nginx-ingress-controller:0.10.2-rancher3"),
			IngressBackend:            m("k8s.gcr.io/defaultbackend:1.4"),
		},
		"v1.10.5-rancher1-2": {
			Etcd:                      m("quay.io/coreos/etcd:v3.1.12"),
			Kubernetes:                m("rancher/hyperkube:v1.10.5-rancher1"),
			Alpine:                    m("rancher/rke-tools:v0.1.13"),
			NginxProxy:                m("rancher/rke-tools:v0.1.13"),
			CertDownloader:            m("rancher/rke-tools:v0.1.13"),
			KubernetesServicesSidecar: m("rancher/rke-tools:v0.1.13"),
			KubeDNS:                   m("gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.8"),
			DNSmasq:                   m("gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.8"),
			KubeDNSSidecar:            m("gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.8"),
			KubeDNSAutoscaler:         m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
			Flannel:                   m("quay.io/coreos/flannel:v0.9.1"),
			FlannelCNI:                m("quay.io/coreos/flannel-cni:v0.2.0"),
			CalicoNode:                m("quay.io/calico/node:v3.1.1"),
			CalicoCNI:                 m("quay.io/calico/cni:v3.1.1"),
			CalicoCtl:                 m("quay.io/calico/ctl:v2.0.0"),
			CanalNode:                 m("quay.io/calico/node:v3.1.1"),
			CanalCNI:                  m("quay.io/calico/cni:v3.1.1"),
			CanalFlannel:              m("quay.io/coreos/flannel:v0.9.1"),
			WeaveNode:                 m("weaveworks/weave-kube:2.1.2"),
			WeaveCNI:                  m("weaveworks/weave-npc:2.1.2"),
			PodInfraContainer:         m("gcr.io/google_containers/pause-amd64:3.1"),
			Ingress:                   m("rancher/nginx-ingress-controller:0.16.2-rancher1"),
			IngressBackend:            m("k8s.gcr.io/defaultbackend:1.4"),
			MetricsServer:             m("gcr.io/google_containers/metrics-server-amd64:v0.2.1"),
		},
		"v1.11.1-rancher1-1": {
			Etcd:                      m("quay.io/coreos/etcd:v3.2.18"),
			Kubernetes:                m("rancher/hyperkube:v1.11.1-rancher1"),
			Alpine:                    m("rancher/rke-tools:v0.1.13"),
			NginxProxy:                m("rancher/rke-tools:v0.1.13"),
			CertDownloader:            m("rancher/rke-tools:v0.1.13"),
			KubernetesServicesSidecar: m("rancher/rke-tools:v0.1.13"),
			KubeDNS:                   m("gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.10"),
			DNSmasq:                   m("gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.10"),
			KubeDNSSidecar:            m("gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.10"),
			KubeDNSAutoscaler:         m("gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"),
			Flannel:                   m("quay.io/coreos/flannel:v0.9.1"),
			FlannelCNI:                m("quay.io/coreos/flannel-cni:v0.2.0"),
			CalicoNode:                m("quay.io/calico/node:v3.1.1"),
			CalicoCNI:                 m("quay.io/calico/cni:v3.1.1"),
			CalicoCtl:                 m("quay.io/calico/ctl:v2.0.0"),
			CanalNode:                 m("quay.io/calico/node:v3.1.1"),
			CanalCNI:                  m("quay.io/calico/cni:v3.1.1"),
			CanalFlannel:              m("quay.io/coreos/flannel:v0.9.1"),
			WeaveNode:                 m("weaveworks/weave-kube:2.1.2"),
			WeaveCNI:                  m("weaveworks/weave-npc:2.1.2"),
			PodInfraContainer:         m("gcr.io/google_containers/pause-amd64:3.1"),
			Ingress:                   m("rancher/nginx-ingress-controller:0.16.2-rancher1"),
			IngressBackend:            m("k8s.gcr.io/defaultbackend:1.4"),
			MetricsServer:             m("gcr.io/google_containers/metrics-server-amd64:v0.2.1"),
		},
	}
)

func init() {
	badVersions := map[string]bool{
		"v1.9.7-rancher1":    true,
		"v1.10.1-rancher1":   true,
		"v1.8.11-rancher1":   true,
		"v1.8.10-rancher1-1": true,
	}

	if K8sVersionToRKESystemImages != nil {
		panic("Do not initialize or add values to K8sVersionToRKESystemImages")
	}

	K8sVersionToRKESystemImages = map[string]RKESystemImages{}

	for version, images := range AllK8sVersions {
		if badVersions[version] {
			continue
		}

		longName := "rancher/hyperkube:" + version
		if !strings.HasPrefix(longName, images.Kubernetes) {
			panic(fmt.Sprintf("For K8s version %s, the Kubernetes image tag should be a substring of %s, currently it is %s", version, version, images.Kubernetes))
		}
	}

	for _, latest := range k8sVersionsCurrent {
		images, ok := AllK8sVersions[latest]
		if !ok {
			panic("K8s version " + " is not found in AllK8sVersions map")
		}
		K8sVersionToRKESystemImages[latest] = images
	}

	if _, ok := K8sVersionToRKESystemImages[DefaultK8s]; !ok {
		panic("Default K8s version " + DefaultK8s + " is not found in k8sVersionsCurrent list")
	}
}
