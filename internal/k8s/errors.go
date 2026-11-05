// Package k8s provides read-only access to Kubernetes resources for KubeSay.
//
// It only ever issues get/list requests; no code path constructs a write
// request. The evidence-collection types and functions are added in later
// phases; this file defines the sentinel errors that cmd/kubesay maps onto
// process exit codes.
package k8s

import "errors"

var (
	// ErrKubeconfig indicates the kubeconfig could not be located or parsed.
	// It is returned from client construction and maps to a non-zero exit.
	ErrKubeconfig = errors.New("kubeconfig could not be loaded")

	// ErrNoFailingPod indicates a target workload has no resolvable failing
	// pod. It is returned from downward target resolution and maps to a
	// non-zero exit.
	ErrNoFailingPod = errors.New("no failing pod found for target")
)
