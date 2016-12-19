package client

type (
	// DaemonSetInterface has methods to work with DaemonSet resources.
	DaemonSetInterface interface {
		CreateDaemonSet(namespace string, item *DaemonSet) (*DaemonSet, error)
		GetDaemonSet(namespace, name string) (result *DaemonSet, err error)
		ListDaemonSets(namespace string, opts *ListOptions) (*DaemonSetList, error)
		DeleteDaemonSet(namespace, name string) error
		UpdateDaemonSet(namespace string, item *DaemonSet) (*DaemonSet, error)
	}

	// DaemonSet represents the configuration of a daemon set.
	DaemonSet struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`

		// Specification of the desired behavior of the DaemonSet.
		Spec *DaemonSetSpec `json:"spec,omitempty"`

		// Most recently observed status of the DaemonSet.
		Status *DaemonSetStatus `json:"status,omitempty"`
	}

	// DaemonSetSpec is the specification of a daemon set.
	DaemonSetSpec struct {
		// Selector is a label query over pods that are managed by the daemon set.
		// Must match in order to be controlled. If empty, defaulted to labels on Pod template.
		Selector *LabelSelector `json:"selector,omitempty"`

		// Template is the object that describes the pod that will be created.
		// The DaemonSet will create exactly one copy of this pod on every node that matches the template’s
		// node selector (or on every node if no node selector is specified).
		Template PodTemplateSpec `json:"template"`
	}

	// DaemonSetStatus represents the current status of a daemon set.
	DaemonSetStatus struct {
		// CurrentNumberScheduled is the number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod.
		CurrentNumberScheduled int32 `json:"currentNumberScheduled"`
		// NumberMisscheduled is the number of nodes that are running the daemon pod, but are not supposed to run the daemon pod.
		NumberMisscheduled int32 `json:"numberMisscheduled"`
		// DesiredNumberScheduled is the total number of nodes that should be running the daemon pod (including nodes correctly running the daemon pod).
		DesiredNumberScheduled int32 `json:"desiredNumberScheduled"`
		// NumberReady is the number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready.
		NumberReady int32 `json:"numberReady"`
	}

	DaemonSetList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`

		// Items is the list of daemonsets.
		Items []DaemonSet `json:"items"`
	}
)

// NewDaemonSet creates a new DaemonSet struct
func NewDaemonSet(namespace, name string) *DaemonSet {
	return &DaemonSet{
		TypeMeta:   NewTypeMeta("DaemonSet", "extensions/v1beta1"),
		ObjectMeta: NewObjectMeta(namespace, name),
		Spec:       &DaemonSetSpec{},
	}
}
