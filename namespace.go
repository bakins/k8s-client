package client

type (
	NamespaceInterface interface {
		CreateNamespace(item *Namespace) (*Namespace, error)
		GetNamespace(name string) (result *Namespace, err error)
		ListNamespaces(opts *ListOptions) (*NamespaceList, error)
		DeleteNamespace(name string) error
		UpdateNamespace(item *Namespace) (*Namespace, error)
	}

	NamespaceSpec struct {
		Finalizers []FinalizerName
	}

	NamespacePhase string

	NamespaceStatus struct {
		Phase NamespacePhase `json:"phase,omitempty"`
	}

	Namespace struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`
		Spec       *NamespaceSpec   `json:"spec,omitempty"`
		Status     *NamespaceStatus `json:"status,omitempty"`
	}

	NamespaceList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`

		Items []Namespace `json:"items"`
	}
)

// NewNamespace creates a new namespace struct
func NewNamespace(name string) *Namespace {
	return &Namespace{
		TypeMeta:   NewTypeMeta("Namespace", "v1"),
		ObjectMeta: NewObjectMeta(name, name),
		Spec:       &NamespaceSpec{},
	}
}
