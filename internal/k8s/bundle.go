package k8s

// EvidenceBundle is everything KubeSay gathers about one failing pod. It is
// the input to the reasoner. Components that could not be found are left
// empty rather than treated as errors.
type EvidenceBundle struct {
	Pod            PodStatus       `json:"pod"`             // the pod's status snapshot
	Events         []EventRecord   `json:"events"`          // empty if none
	CurrentLogs    []ContainerLog  `json:"current_logs"`    // one per container
	PreviousLogs   []ContainerLog  `json:"previous_logs"`   // empty if no prior run
	OwnerChain     []OwnerRef      `json:"owner_chain"`     // pod -> ... -> rootmost owner
	ResourceLimits []ResourceLimit `json:"resource_limits"` // per container
}

// PodStatus is a trimmed view of a pod's identity and current state.
type PodStatus struct {
	Namespace  string            `json:"namespace"`
	Name       string            `json:"name"`
	Phase      string            `json:"phase"`      // e.g. "Running", "Pending", "Failed"
	Conditions []string          `json:"conditions"` // human-readable pod conditions
	Containers []ContainerStatus `json:"containers"` // one entry per container
}

// ContainerStatus describes the state of one container in the pod.
type ContainerStatus struct {
	Name         string `json:"name"`
	Ready        bool   `json:"ready"`         // whether the container passes its readiness check
	RestartCount int32  `json:"restart_count"` // how many times it has been restarted
	State        string `json:"state"`         // Running|Waiting|Terminated
	Reason       string `json:"reason"`        // e.g. CrashLoopBackOff, OOMKilled, ImagePullBackOff
	// ExitCode is a pointer so that nil can mean "has not exited". A plain
	// int32 would report 0, which looks like a clean exit.
	ExitCode *int32 `json:"exit_code,omitempty"`
}

// EventRecord is one Kubernetes event related to the pod.
type EventRecord struct {
	Type    string `json:"type"`    // Normal|Warning
	Reason  string `json:"reason"`  // short machine-style cause, e.g. "BackOff"
	Message string `json:"message"` // human-readable detail
	Count   int32  `json:"count"`   // how many times the event has occurred
}

// ContainerLog holds the log text of one container.
type ContainerLog struct {
	Container string `json:"container"` // container name
	Previous  bool   `json:"previous"`  // true if from the previous (crashed) run
	Lines     string `json:"lines"`     // tail-limited log text
}

// OwnerRef identifies one resource that owns another, such as the
// ReplicaSet that owns a pod.
type OwnerRef struct {
	Kind string `json:"kind"` // e.g. "ReplicaSet", "Deployment"
	Name string `json:"name"`
	UID  string `json:"uid"` // Kubernetes unique ID of the owner
}

// ResourceLimit records a container's CPU and memory requests and limits.
// Values are Kubernetes quantity strings such as "500m" or "256Mi"; a field
// is omitted when it is not set.
type ResourceLimit struct {
	Container     string `json:"container"`
	CPURequest    string `json:"cpu_request,omitempty"`
	CPULimit      string `json:"cpu_limit,omitempty"`
	MemoryRequest string `json:"memory_request,omitempty"`
	MemoryLimit   string `json:"memory_limit,omitempty"`
}
