package k8s

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ClientsetInterface defines the interface for Kubernetes clientset
type ClientsetInterface interface {
	// Add methods as needed
}

// Client represents a Kubernetes client
type Client struct {
	clientset kubernetes.Interface
	config    *rest.Config
	namespace string
}

// NewClient creates a new Kubernetes client
func NewClient(kubeConfigPath, namespace string) (*Client, error) {
	var config *rest.Config
	var err error

	// Try to use in-cluster config if kubeConfigPath is empty
	if kubeConfigPath == "" {
		config, err = rest.InClusterConfig()
		if err != nil {
			// Fall back to kubeconfig in home directory
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("failed to get user home directory: %w", err)
			}
			kubeConfigPath = filepath.Join(homeDir, ".kube", "config")
		}
	}

	// If we still need to load from kubeconfig file
	if config == nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
		}
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	// Use default namespace if not specified
	if namespace == "" {
		namespace = "default"
	}

	return &Client{
		clientset: clientset,
		config:    config,
		namespace: namespace,
	}, nil
}

// GetClientset returns the Kubernetes clientset
func (c *Client) GetClientset() kubernetes.Interface {
	return c.clientset
}

// GetConfig returns the Kubernetes REST config
func (c *Client) GetConfig() *rest.Config {
	return c.config
}

// GetNamespace returns the current namespace
func (c *Client) GetNamespace() string {
	return c.namespace
}

// SetNamespace sets the current namespace
func (c *Client) SetNamespace(namespace string) {
	c.namespace = namespace
}