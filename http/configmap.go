package http

import (
	k8s "github.com/YakLabs/k8s-client"
	"github.com/pkg/errors"
)

type (
	watchEventConfigMap struct {
		raw    k8s.WatchEvent
		object *k8s.ConfigMap
	}
)

func (w *watchEventConfigMap) Type() k8s.WatchEventType {
	return w.raw.Type
}

func (w *watchEventConfigMap) Object() (*k8s.ConfigMap, error) {
	if w.object != nil {
		return w.object, nil
	}
	if w.raw.Type == k8s.WatchEventTypeError {
		var status k8s.Status
		if err := w.raw.UnmarshalObject(&status); err != nil {
			return nil, errors.Wrap(err, "failed to decode Status")
		}
		return nil, &status
	}
	var object k8s.ConfigMap
	if err := w.raw.UnmarshalObject(&object); err != nil {
		return nil, errors.Wrap(err, "failed to decode ConfigMap")
	}
	w.object = &object
	return &object, nil
}

func configmapGeneratePath(namespace, name string) string {
	if namespace == "" && name == "" {
		return "/api/v1/configmaps"
	}
	if name == "" {
		return "/api/v1/namespaces/" + namespace + "/configmaps"
	}
	return "/api/v1/namespaces/" + namespace + "/configmaps/" + name
}

// GetConfigMap fetches a single ConfigMap
func (c *Client) GetConfigMap(namespace, name string) (*k8s.ConfigMap, error) {
	var out k8s.ConfigMap
	_, err := c.do("GET", configmapGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ConfigMap")
	}
	return &out, nil
}

// CreateConfigMap creates a new ConfigMap. This will fail if it already exists.
func (c *Client) CreateConfigMap(namespace string, item *k8s.ConfigMap) (*k8s.ConfigMap, error) {
	item.TypeMeta.Kind = "ConfigMap"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.ConfigMap
	_, err := c.do("POST", configmapGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ConfigMap")
	}
	return &out, nil
}

// ListConfigMaps lists all ConfigMaps in a namespace
func (c *Client) ListConfigMaps(namespace string, opts *k8s.ListOptions) (*k8s.ConfigMapList, error) {
	var out k8s.ConfigMapList
	_, err := c.do("GET", configmapGeneratePath(namespace, "")+"?"+listOptionsQuery(opts, nil), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list ConfigMaps")
	}
	return &out, nil
}

// WatchConfigMaps watches all ConfigMap changes in a namespace
func (c *Client) WatchConfigMaps(namespace string, opts *k8s.WatchOptions, events chan k8s.ConfigMapWatchEvent) error {
	if events == nil {
		return errors.New("events must not be nil")
	}
	rawEvents := make(chan k8s.WatchEvent)
	go func() {
		for rawEvent := range rawEvents {
			events <- &watchEventConfigMap{raw: rawEvent}
		}
		close(events)
	}()
	_, err := c.doWatch("GET", configmapGeneratePath(namespace, "")+"?"+watchOptionsQuery(opts), nil, rawEvents)
	if err != nil {
		return errors.Wrap(err, "failed to watch ConfigMaps")
	}
	return nil
}

// DeleteConfigMap deletes a single ConfigMap. It will error if the ConfigMap does not exist.
func (c *Client) DeleteConfigMap(namespace, name string) error {
	_, err := c.do("DELETE", configmapGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete ConfigMap")
}

// UpdateConfigMap will update in place a single ConfigMap. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdateConfigMap(namespace string, item *k8s.ConfigMap) (*k8s.ConfigMap, error) {
	item.TypeMeta.Kind = "ConfigMap"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.ConfigMap
	_, err := c.do("PUT", configmapGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update ConfigMap")
	}
	return &out, nil
}
