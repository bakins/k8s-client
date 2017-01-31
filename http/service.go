package http

import (
	k8s "github.com/bakins/k8s-client"
	"github.com/pkg/errors"
)

type (
	watchEventService struct {
		raw    k8s.WatchEvent
		object *k8s.Service
	}
)

func (w *watchEventService) Type() k8s.WatchEventType {
	return w.raw.Type
}

func (w *watchEventService) Object() (*k8s.Service, error) {
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
	var object k8s.Service
	if err := w.raw.UnmarshalObject(&object); err != nil {
		return nil, errors.Wrap(err, "failed to decode Service")
	}
	w.object = &object
	return &object, nil
}

func serviceGeneratePath(namespace, name string) string {
	if namespace == "" && name == "" {
		return "/api/v1/services"
	}
	if name == "" {
		return "/api/v1/namespaces/" + namespace + "/services"
	}
	return "/api/v1/namespaces/" + namespace + "/services/" + name
}

// GetService fetches a single Service
func (c *Client) GetService(namespace, name string) (*k8s.Service, error) {
	var out k8s.Service
	_, err := c.do("GET", serviceGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Service")
	}
	return &out, nil
}

// CreateService creates a new Service. This will fail if it already exists.
func (c *Client) CreateService(namespace string, item *k8s.Service) (*k8s.Service, error) {
	item.TypeMeta.Kind = "Service"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Service
	_, err := c.do("POST", serviceGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Service")
	}
	return &out, nil
}

// ListServices lists all Services in a namespace
func (c *Client) ListServices(namespace string, opts *k8s.ListOptions) (*k8s.ServiceList, error) {
	var out k8s.ServiceList
	_, err := c.do("GET", serviceGeneratePath(namespace, "")+"?"+listOptionsQuery(opts, nil), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list Services")
	}
	return &out, nil
}

// WatchServices watches all Service changes in a namespace
func (c *Client) WatchServices(namespace string, opts *k8s.WatchOptions, events chan k8s.ServiceWatchEvent) error {
	if events == nil {
		return errors.New("events must not be nil")
	}
	rawEvents := make(chan k8s.WatchEvent)
	go func() {
		for rawEvent := range rawEvents {
			events <- &watchEventService{raw: rawEvent}
		}
		close(events)
	}()
	_, err := c.doWatch("GET", serviceGeneratePath(namespace, "")+"?"+watchOptionsQuery(opts), nil, rawEvents)
	if err != nil {
		return errors.Wrap(err, "failed to watch Services")
	}
	return nil
}

// DeleteService deletes a single Service. It will error if the Service does not exist.
func (c *Client) DeleteService(namespace, name string) error {
	_, err := c.do("DELETE", serviceGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete Service")
}

// UpdateService will update in place a single Service. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdateService(namespace string, item *k8s.Service) (*k8s.Service, error) {
	item.TypeMeta.Kind = "Service"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Service
	_, err := c.do("PUT", serviceGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Service")
	}
	return &out, nil
}
