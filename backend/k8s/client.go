package k8s

import (
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientset      *kubernetes.Clientset
	kubeconfigPath string
}

func NewClient(kubeconfigPath string) *Client {
	if kubeconfigPath == "" {
		// Try environment variable
		kubeconfigPath = os.Getenv("KUBECONFIG")
		if kubeconfigPath == "" {
			// Default to home directory
			kubeconfigPath = os.ExpandEnv("$HOME/.kube/config")
		}
	}
	return &Client{
		kubeconfigPath: kubeconfigPath,
	}
}

func (c *Client) GetClient() (*kubernetes.Clientset, error) {
	if c.clientset != nil {
		return c.clientset, nil
	}

	// Try in-cluster config first
	config, err := rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", c.kubeconfigPath)
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	c.clientset = clientset
	return c.clientset, nil
}

func (c *Client) GetConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", c.kubeconfigPath)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}
