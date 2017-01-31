package http

import (
	k8s "github.com/bakins/k8s-client"
	"github.com/pkg/errors"
)

type (
	watchEventEndpoints struct {
		raw    k8s.WatchEvent
		object *k8s.Endpoints
	}
)

func (w *watchEventEndpoints) Type() k8s.WatchEventType {
	return w.raw.Type
}

func (w *watchEventEndpoints) Object() (*k8s.Endpoints, error) {
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
	var object k8s.Endpoints
	if err := w.raw.UnmarshalObject(&object); err != nil {
		return nil, errors.Wrap(err, "failed to decode Endpoints")
	}
	w.object = &object
	return &object, nil
}

func endpointsGeneratePath(namespace, name string) string {
	if namespace == "" && name == "" {
		return "/api/v1/endpoints"
	}
	if name == "" {
		return "/api/v1/namespaces/" + namespace + "/endpoints"
	}
	return "/api/v1/namespaces/" + namespace + "/endpoints/" + name
}

// GetEndpoints fetches a single Endpoints
func (c *Client) GetEndpoints(namespace, name string) (*k8s.Endpoints, error) {
	var out k8s.Endpoints
	_, err := c.do("GET", endpointsGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Endpoints")
	}
	return &out, nil
}

// CreateEndpoints creates a new Endpoints. This will fail if it already exists.
func (c *Client) CreateEndpoints(namespace string, item *k8s.Endpoints) (*k8s.Endpoints, error) {
	item.TypeMeta.Kind = "Endpoints"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Endpoints
	_, err := c.do("POST", endpointsGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Endpoints")
	}
	return &out, nil
}

// ListEndpoints lists all Endpointss in a namespace
func (c *Client) ListEndpoints(namespace string, opts *k8s.ListOptions) (*k8s.EndpointsList, error) {
	var out k8s.EndpointsList
	_, err := c.do("GET", endpointsGeneratePath(namespace, "")+"?"+listOptionsQuery(opts, nil), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list Endpointss")
	}
	return &out, nil
}

// WatchEndpoints watches all Endpoints changes in a namespace
func (c *Client) WatchEndpoints(namespace string, opts *k8s.WatchOptions, events chan k8s.EndpointsWatchEvent) error {
	if events == nil {
		return errors.New("events must not be nil")
	}
	rawEvents := make(chan k8s.WatchEvent)
	go func() {
		for rawEvent := range rawEvents {
			events <- &watchEventEndpoints{raw: rawEvent}
		}
		close(events)
	}()
	_, err := c.doWatch("GET", endpointsGeneratePath(namespace, "")+"?"+watchOptionsQuery(opts), nil, rawEvents)
	if err != nil {
		return errors.Wrap(err, "failed to watch Endpointss")
	}
	return nil
}

// DeleteEndpoints deletes a single Endpoints. It will error if the Endpoints does not exist.
func (c *Client) DeleteEndpoints(namespace, name string) error {
	_, err := c.do("DELETE", endpointsGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete Endpoints")
}

// UpdateEndpoints will update in place a single Endpoints. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdateEndpoints(namespace string, item *k8s.Endpoints) (*k8s.Endpoints, error) {
	item.TypeMeta.Kind = "Endpoints"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Endpoints
	_, err := c.do("PUT", endpointsGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Endpoints")
	}
	return &out, nil
}
