package k8s

import (
	"log"

	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

// NewMockClient creates a new mock Kubernetes client for development and testing
func NewMockClient() *Client {
	log.Println("Creating mock Kubernetes client")
	
	// Create a fake clientset
	clientset := fake.NewSimpleClientset()
	
	// Create a dummy config
	config := &rest.Config{
		Host: "https://mock-kubernetes.local",
	}
	
	// Convert fake clientset to the interface type expected by Client
	return &Client{
		clientset: clientset,
		config:    config,
		namespace: "default",
	}
}