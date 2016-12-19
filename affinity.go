package client

const (
	// AffinityAnnotationKey represents the key of affinity data (json serialized)
	// in the Annotations of a Pod.
	AffinityAnnotationKey string = "scheduler.alpha.kubernetes.io/affinity"
)

type (
	// Affinity is a group of affinity scheduling rules.
	Affinity struct {
		// Describes node affinity scheduling rules for the pod.
		NodeAffinity *NodeAffinity `json:"nodeAffinity,omitempty"`
		// Describes pod affinity scheduling rules (e.g. co-locate this pod in the same node, zone, etc. as some other pod(s)).
		PodAffinity *PodAffinity `json:"podAffinity,omitempty"`
		// Describes pod anti-affinity scheduling rules (e.g. avoid putting this pod in the same node, zone, etc. as some other pod(s)).
		PodAntiAffinity *PodAntiAffinity `json:"podAntiAffinity,omitempty"`
	}

	// NodeAffinity is a group of node affinity scheduling rules.
	NodeAffinity struct {
		// If the affinity requirements specified by this field are not met at
		// scheduling time, the pod will not be scheduled onto the node.
		// If the affinity requirements specified by this field cease to be met
		// at some point during pod execution (e.g. due to an update), the system
		// may or may not try to eventually evict the pod from its node.
		RequiredDuringSchedulingIgnoredDuringExecution *NodeSelector `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`

		// The scheduler will prefer to schedule pods to nodes that satisfy
		// the affinity expressions specified by this field, but it may choose
		// a node that violates one or more of the expressions. The node that is
		// most preferred is the one with the greatest sum of weights, i.e.
		// for each node that meets all of the scheduling requirements (resource
		// request, requiredDuringScheduling affinity expressions, etc.),
		// compute a sum by iterating through the elements of this field and adding
		// "weight" to the sum if the node matches the corresponding matchExpressions; the
		// node(s) with the highest sum are the most preferred.
		PreferredDuringSchedulingIgnoredDuringExecution []PreferredSchedulingTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`
	}

	// PodAffinity is a group of inter pod affinity scheduling rules.
	PodAffinity struct {
		// RequiredDuringSchedulingRequiredDuringExecution []PodAffinityTerm  `json:"requiredDuringSchedulingRequiredDuringExecution,omitempty"`
		// If the affinity requirements specified by this field are not met at
		// scheduling time, the pod will not be scheduled onto the node.
		// If the affinity requirements specified by this field cease to be met
		// at some point during pod execution (e.g. due to a pod label update), the
		// system may or may not try to eventually evict the pod from its node.
		// When there are multiple elements, the lists of nodes corresponding to each
		// podAffinityTerm are intersected, i.e. all terms must be satisfied.
		RequiredDuringSchedulingIgnoredDuringExecution []PodAffinityTerm `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`

		// The scheduler will prefer to schedule pods to nodes that satisfy
		// the affinity expressions specified by this field, but it may choose
		// a node that violates one or more of the expressions. The node that is
		// most preferred is the one with the greatest sum of weights, i.e.
		// for each node that meets all of the scheduling requirements (resource
		// request, requiredDuringScheduling affinity expressions, etc.),
		// compute a sum by iterating through the elements of this field and adding
		// "weight" to the sum if the node has pods which matches the corresponding podAffinityTerm; the
		// node(s) with the highest sum are the most preferred.
		PreferredDuringSchedulingIgnoredDuringExecution []WeightedPodAffinityTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`
	}

	// PodAntiAffinity is a group of inter pod anti affinity scheduling rules.
	PodAntiAffinity struct {
		// RequiredDuringSchedulingRequiredDuringExecution []PodAffinityTerm  `json:"requiredDuringSchedulingRequiredDuringExecution,omitempty"`
		// If the anti-affinity requirements specified by this field are not met at
		// scheduling time, the pod will not be scheduled onto the node.
		// If the anti-affinity requirements specified by this field cease to be met
		// at some point during pod execution (e.g. due to a pod label update), the
		// system may or may not try to eventually evict the pod from its node.
		// When there are multiple elements, the lists of nodes corresponding to each
		// podAffinityTerm are intersected, i.e. all terms must be satisfied.
		RequiredDuringSchedulingIgnoredDuringExecution []PodAffinityTerm `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`

		// The scheduler will prefer to schedule pods to nodes that satisfy
		// the anti-affinity expressions specified by this field, but it may choose
		// a node that violates one or more of the expressions. The node that is
		// most preferred is the one with the greatest sum of weights, i.e.
		// for each node that meets all of the scheduling requirements (resource
		// request, requiredDuringScheduling anti-affinity expressions, etc.),
		// compute a sum by iterating through the elements of this field and adding
		// "weight" to the sum if the node has pods which matches the corresponding podAffinityTerm; the
		// node(s) with the highest sum are the most preferred.
		PreferredDuringSchedulingIgnoredDuringExecution []WeightedPodAffinityTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`
	}

	// NodeSelector represents the union of the results of one or more label queries
	// over a set of nodes; that is, it represents the OR of the selectors represented
	// by the node selector terms.
	NodeSelector struct {
		// Required. A list of node selector terms. The terms are ORed.
		NodeSelectorTerms []NodeSelectorTerm `json:"nodeSelectorTerms"`
	}

	// An empty preferred scheduling term matches all objects with implicit weight 0
	// (i.e. it's a no-op). A null preferred scheduling term matches no objects (i.e. is also a no-op).
	PreferredSchedulingTerm struct {
		// Weight associated with matching the corresponding nodeSelectorTerm, in the range 1-100.
		Weight int32 `json:"weight,omitempty"`

		// A node selector term, associated with the corresponding weight.
		Preference *NodeSelectorTerm `json:"preference,omitempty"`
	}

	// Defines a set of pods (namely those matching the labelSelector
	// relative to the given namespace(s)) that this pod should be
	// co-located (affinity) or not co-located (anti-affinity) with,
	// where co-located is defined as running on a node whose value of
	// the label with key <topologyKey> tches that of any node on which
	// a pod of the set of pods is running
	PodAffinityTerm struct {
		// A label query over a set of resources, in this case pods.
		LabelSelector *LabelSelector `json:"labelSelector,omitempty"`

		// namespaces specifies which namespaces the labelSelector applies to (matches against);
		// nil list means "this pod's namespace," empty list means "all namespaces"
		// The json tag here is not "omitempty" since we need to distinguish nil and empty.
		// See https://golang.org/pkg/encoding/json/#Marshal for more details.
		Namespaces []string `json:"namespaces,omitempty"`

		// This pod should be co-located (affinity) or not co-located (anti-affinity) with the pods matching
		// the labelSelector in the specified namespaces, where co-located is defined as running on a node
		// whose value of the label with key topologyKey matches that of any node on which any of the
		// selected pods is running.
		// For PreferredDuringScheduling pod anti-affinity, empty topologyKey is interpreted as "all topologies"
		// ("all topologies" here means all the topologyKeys indicated by scheduler command-line argument --failure-domains);
		// for affinity and for RequiredDuringScheduling pod anti-affinity, empty topologyKey is not allowed.
		TopologyKey string `json:"topologyKey,omitempty"`
	}

	// The weights of all of the matched WeightedPodAffinityTerm fields are added per-node to find the most preferred node(s)
	WeightedPodAffinityTerm struct {
		// weight associated with matching the corresponding podAffinityTerm, in the range 1-100.
		Weight int32 `json:"weight,omitempty"`

		// Required. A pod affinity term, associated with the corresponding weight.
		PodAffinityTerm PodAffinityTerm `json:"podAffinityTerm"`
	}

	// A null or empty node selector term matches no objects.
	NodeSelectorTerm struct {
		// Required. A list of node selector requirements. The requirements are ANDed.
		MatchExpressions []NodeSelectorRequirement `json:"matchExpressions"`
	}

	// A node selector requirement is a selector that contains values, a key, and an operator
	// that relates the key and values.
	NodeSelectorRequirement struct {
		// The label key that the selector applies to.
		Key string `json:"key,omitempty"`

		// Represents a key's relationship to a set of values.
		// Valid operators are In, NotIn, Exists, DoesNotExist. Gt, and Lt.
		Operator string `json:"operator,omitempty"`

		// An array of string values. If the operator is In or NotIn,
		// the values array must be non-empty. If the operator is Exists or DoesNotExist,
		// the values array must be empty. If the operator is Gt or Lt, the values
		// array must have a single element, which will be interpreted as an integer.
		// This array is replaced during a strategic merge patch.
		Values []string `json:"values,omitempty"`
	}
)
