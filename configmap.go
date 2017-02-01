package client

type (
	ConfigMapInterface interface {
		CreateConfigMap(namespace string, item *ConfigMap) (*ConfigMap, error)
		GetConfigMap(namespace, name string) (result *ConfigMap, err error)
		ListConfigMaps(namespace string, opts *ListOptions) (*ConfigMapList, error)
		WatchConfigMaps(namespace string, opts *WatchOptions, events chan ConfigMapWatchEvent) error
		DeleteConfigMap(namespace, name string) error
		UpdateConfigMap(namespace string, item *ConfigMap) (*ConfigMap, error)
	}

	ConfigMapWatchEvent interface {
		Type() WatchEventType
		Object() (*ConfigMap, error)
	}

	ConfigMapType string

	ConfigMap struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`
		Data       map[string]string `json:"data,omitempty"`
	}

	ConfigMapList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`

		Items []ConfigMap `json:"items"`
	}
)

// NewConfigMap creates a new ConfigMap struct
func NewConfigMap(namespace, name string) *ConfigMap {
	return &ConfigMap{
		TypeMeta:   NewTypeMeta("ConfigMap", "v1"),
		ObjectMeta: NewObjectMeta(namespace, name),
		Data:       make(map[string]string),
	}
}
