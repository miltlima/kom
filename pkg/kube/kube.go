package kube

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

func GetConfig() (*rest.Config, error) {
	kubeconfig := os.Getenv("HOME") + "/.kube/config"

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("Dont possible to get kubeconfig: %v", err)
		}
	}
	return config, nil
}

func NewClientSet(config *rest.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func NewMetricsClientSet(config *rest.Config) (*versioned.Clientset, error) {
	metricsClientset, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return metricsClientset, nil
}
