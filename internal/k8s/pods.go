package k8s

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PodSummary is the compact, display-oriented view of a pod used by the
// `list` command. The json tags keep the --output json format stable.
type PodSummary struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Phase     string `json:"phase"`    // e.g. "Running", "Pending"
	Ready     string `json:"ready"`    // ready/total containers, e.g. "1/2"
	Restarts  int32  `json:"restarts"` // sum of restarts across containers
	Age       string `json:"age"`      // e.g. "5m", "3h", "2d"
}

// ListPods returns a summary of every pod in namespace ns.
//
// The cs parameter is the kubernetes.Interface interface rather than the
// concrete *kubernetes.Clientset. *Clientset satisfies the interface, so
// production code passes the real client, while tests can pass the client-go
// fake clientset without a cluster.
//
// It performs a single read-only list call.
func ListPods(ctx context.Context, cs kubernetes.Interface, ns string) ([]PodSummary, error) {
	list, err := cs.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing pods in namespace %q: %w", ns, err)
	}

	now := time.Now()
	summaries := make([]PodSummary, 0, len(list.Items))
	for _, pod := range list.Items {
		var ready, restarts int32
		for _, status := range pod.Status.ContainerStatuses {
			if status.Ready {
				ready++
			}
			restarts += status.RestartCount
		}

		summaries = append(summaries, PodSummary{
			Namespace: pod.Namespace,
			Name:      pod.Name,
			Phase:     string(pod.Status.Phase),
			Ready:     fmt.Sprintf("%d/%d", ready, len(pod.Status.ContainerStatuses)),
			Restarts:  restarts,
			Age:       formatAge(now.Sub(pod.CreationTimestamp.Time)),
		})
	}
	return summaries, nil
}

// GetPod fetches a single pod by its exact name in namespace ns. This differs
// from ListPods, which lists every pod in a namespace and returns summaries;
// GetPod returns the full, raw pod object for one known pod.
//
// It performs a single read-only get call.
func GetPod(ctx context.Context, cs kubernetes.Interface, ns, name string) (*corev1.Pod, error) {
	pod, err := cs.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("getting pod %s/%s: %w", ns, name, err)
	}
	return pod, nil
}

// formatAge renders a duration compactly using only its largest unit,
// truncated (not rounded): "30s", "5m", "3h", "2d". Negative durations
// (clock skew) are treated as zero.
func formatAge(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}
