package client

type (
	NodeInterface interface {
		CreateNode(item *Node) (*Node, error)
		GetNode(name string) (result *Node, err error)
		ListNodes(opts *ListOptions) (*NodeList, error)
		DeleteNode(name string) error
		UpdateNode(item *Node) (*Node, error)
	}

	NodeSpec struct {
		PodCIDR       string `json:"podCIDR,omitempty"`
		ExternalID    string `json:"externalID,omitempty"`
		ProviderID    string `json:"providerID,omitempty"`
		Unschedulable bool   `json:"unschedulable,omitempty"`
	}

	NodePhase         string
	NodeConditionType string
	NodeAddressType   string

	NodeAddress struct {
		Type    NodeAddressType `json:"type"`
		Address string          `json:"address"`
	}

	NodeCondition struct {
		Type    NodeConditionType `json:"type"`
		Status  ConditionStatus   `json:"status"`
		Reason  string            `json:"reason,omitempty"`
		Message string            `json:"message,omitempty"`
	}

	NodeStatus struct {
		Phase      NodePhase       `json:"phase,omitempty"`
		Conditions []NodeCondition `json:"conditions,omitempty"`
		Addresses  []NodeAddress   `json:"addresses,omitempty"`
	}

	Node struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`
		Spec       NodeSpec   `json:"spec,omitempty"`
		Status     NodeStatus `json:"status,omitempty"`
	}

	NodeList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`
		Items    []Node `json:"items"`
	}
)
