package client

import (
	"encoding/json"
	"time"
)

// Common object elements

type (
	UID             string
	FinalizerName   string
	ConditionStatus string

	TypeMeta struct {
		Kind       string `json:"kind,omitempty"`
		APIVersion string `json:"apiVersion,omitempty"`
	}

	ObjectMeta struct {
		Name            string            `json:"name,omitempty"`
		Namespace       string            `json:"namespace,omitempty"`
		SelfLink        string            `json:"selfLink,omitempty"`
		UID             UID               `json:"uid,omitempty"`
		ResourceVersion string            `json:"resourceVersion,omitempty"`
		Generation      int64             `json:"generation,omitempty"`
		Labels          map[string]string `json:"labels,omitempty"`
		Annotations     map[string]string `json:"annotations,omitempty"`
	}

	ListMeta struct {
		SelfLink        string `json:"selfLink,omitempty"`
		ResourceVersion string `json:"resourceVersion,omitempty"`
	}

	ObjectReference struct {
		Kind            string `json:"kind,omitempty"`
		Namespace       string `json:"namespace,omitempty"`
		Name            string `json:"name,omitempty"`
		UID             UID    `json:"uid,omitempty"`
		APIVersion      string `json:"apiVersion,omitempty"`
		ResourceVersion string `json:"resourceVersion,omitempty"`
		FieldPath       string `json:"fieldPath,omitempty"`
	}

	LocalObjectReference struct {
		Name string
	}

	Time struct {
		time.Time
	}

	ObjectFieldSelector struct {
		APIVersion string `json:"apiVersion"`
		FieldPath  string `json:"fieldPath"`
	}

	ConfigMapKeySelector struct {
		LocalObjectReference `json:",inline"`
		Key                  string `json:"key"`
	}

	SecretKeySelector struct {
		LocalObjectReference `json:",inline"`
		Key                  string `json:"key"`
	}

	LabelSelector struct {
		// matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels
		// map is equivalent to an element of matchExpressions, whose key field is "key", the
		// operator is "In", and the values array contains only "value". The requirements are ANDed.
		MatchLabels map[string]string `json:"matchLabels,omitempty" protobuf:"bytes,1,rep,name=matchLabels"`
	}

	FieldSelector map[string]string

	Object interface {
		GetKind() string
	}
	NamespacedObject interface {
		Object
		GetNamespace() string
		GetAnnotations() map[string]string
	}
	ListObject interface {
		Object
		GetItems()
	}
)

func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		// Encode unset/nil objects as JSON's "null".
		return []byte("null"), nil
	}

	return json.Marshal(t.UTC().Format(time.RFC3339))
}

func (t *TypeMeta) GetKind() string {
	return t.Kind
}

func (o *ObjectMeta) GetNamespace() string {
	return o.Namespace
}

func (o *ObjectMeta) GetAnnotations() map[string]string {
	return o.Annotations
}
