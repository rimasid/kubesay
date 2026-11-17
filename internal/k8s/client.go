package k8s

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClient builds a Kubernetes clientset from a kubeconfig and returns it
// together with the default namespace of the kubeconfig's current context.
//
// The kubeconfig is resolved in this order:
//  1. kubeconfigPath, if non-empty
//  2. the KUBECONFIG environment variable
//  3. ~/.kube/config
//
// Any failure to locate, parse, or build a client from the kubeconfig is
// wrapped with ErrKubeconfig, so callers can test it with errors.Is.
// NewClient only constructs the client; it sends no requests to the cluster.
func NewClient(kubeconfigPath string) (*kubernetes.Clientset, string, error) {
	// The default loading rules already honor KUBECONFIG and ~/.kube/config.
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeconfigPath != "" {
		// An explicit path takes precedence over the defaults.
		rules.ExplicitPath = kubeconfigPath
	}

	// Non-interactive: never prompt for credentials on stdin.
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		rules, &clientcmd.ConfigOverrides{})

	// Namespace() returns the current context's namespace, or "default".
	namespace, _, err := clientConfig.Namespace()
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrKubeconfig, err)
	}

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrKubeconfig, err)
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrKubeconfig, err)
	}

	return clientset, namespace, nil
}
