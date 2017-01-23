package client

type (
	// EndpointsInterface has methods to work with Endpoints resources.
	EndpointsInterface interface {
		CreateEndpoints(namespace string, item *Endpoints) (*Endpoints, error)
		GetEndpoints(namespace, name string) (result *Endpoints, err error)
		ListEndpoints(namespace string, opts *ListOptions) (*EndpointsList, error)
		DeleteEndpoints(namespace, name string) error
		UpdateEndpoints(namespace string, item *Endpoints) (*Endpoints, error)
	}

	// Endpoints is a collection of endpoints that implement the actual service.
	Endpoints struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`

		// The set of all endpoints is the union of all subsets.
		// Addresses are placed into subsets according to the IPs they share.
		// A single address with multiple ports, some of which are ready and some of which are not
		// (because they come from different containers) will result in the address being displayed in different
		// subsets for the different ports.
		// No address will appear in both Addresses and NotReadyAddresses in the same subset.
		// Sets of addresses and ports that comprise a service
		Subsets []EndpointSubset `json:"subsets,omitempty"`
	}

	// EndpointSubset is a group of addresses with a common set of ports.
	// The expanded set of endpoints is the Cartesian product of Addresses x Ports.
	EndpointSubset struct {
		// IP addresses which offer the related ports that are marked as ready.
		// These endpoints should be considered safe for load balancers and clients to utilize.
		Addresses []EndpointAddress `json:"addresses,omitempty"`

		// IP addresses which offer the related ports but are not currently marked as ready because
		// they have not yet finished starting, have recently failed a readiness check, or have recently failed a liveness check.
		NotReadyAddresses []EndpointAddress `json:"notReadyAddresses,omitempty"`

		// Port numbers available on the related IP addresses
		Ports []EndpointPort `json:"ports,omitempty"`
	}

	EndpointAddress struct {
		// The IP of this endpoint.
		// May not be loopback (127.0.0.0/8), link-local (169.254.0.0/16), or link-local multicast ((224.0.0.0/24).
		// IPv6 is also accepted but not fully supported on all platforms.
		// Also, certain kubernetes components, like kube-proxy, are not IPv6 ready.
		IP string `json:"ip,omitempty"`
		// The Hostname of this endpoint
		Hostname string `json:"hostname,omitempty"`
		// Optional: Node hosting this endpoint. This can be used to determine endpoints local to a node.
		NodeName string `json:"nodeName,omitempty"`
		// Reference to object providing the endpoint.
		TargetRef *ObjectReference `json:"targetRef,omitempty"`
	}

	EndpointPort struct {
		// The name of this port (corresponds to ServicePort.Name). Must be a DNS_LABEL.
		// Optional only if one port is defined.
		Name string `json:"name,omitempty"`
		// The port number of the endpoint.
		Port int32 `json:"port"`
		// The IP protocol for this port. Must be UDP or TCP. Default is TCP.
		Protocol string `json:"protocol,omitempty"`
	}

	// EndpointsList holds a list of endpoints.
	EndpointsList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`

		Items []Endpoints `json:"items"`
	}
)
