package http

import (
	k8s "github.com/bakins/k8s-client"
	"github.com/pkg/errors"
)

type (
	watchEventEvent struct {
		raw    k8s.WatchEvent
		object *k8s.Event
	}
)

func (w *watchEventEvent) Type() k8s.WatchEventType {
	return w.raw.Type
}

func (w *watchEventEvent) Object() (*k8s.Event, error) {
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
	var object k8s.Event
	if err := w.raw.UnmarshalObject(&object); err != nil {
		return nil, errors.Wrap(err, "failed to decode Event")
	}
	w.object = &object
	return &object, nil
}

func eventGeneratePath(namespace, name string) string {
	if namespace == "" && name == "" {
		return "/api/v1/events"
	}
	if name == "" {
		return "/api/v1/namespaces/" + namespace + "/events"
	}
	return "/api/v1/namespaces/" + namespace + "/events/" + name
}

// GetEvent fetches a single Event
func (c *Client) GetEvent(namespace, name string) (*k8s.Event, error) {
	var out k8s.Event
	_, err := c.do("GET", eventGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Event")
	}
	return &out, nil
}

// CreateEvent creates a new Event. This will fail if it already exists.
func (c *Client) CreateEvent(namespace string, item *k8s.Event) (*k8s.Event, error) {
	item.TypeMeta.Kind = "Event"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Event
	_, err := c.do("POST", eventGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Event")
	}
	return &out, nil
}

// ListEvents lists all Events in a namespace
func (c *Client) ListEvents(namespace string, opts *k8s.ListOptions) (*k8s.EventList, error) {
	var out k8s.EventList
	_, err := c.do("GET", eventGeneratePath(namespace, "")+"?"+listOptionsQuery(opts, nil), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list Events")
	}
	return &out, nil
}

// WatchEvents watches all Event changes in a namespace
func (c *Client) WatchEvents(namespace string, opts *k8s.WatchOptions, events chan k8s.EventWatchEvent) error {
	if events == nil {
		return errors.New("events must not be nil")
	}
	rawEvents := make(chan k8s.WatchEvent)
	go func() {
		for rawEvent := range rawEvents {
			events <- &watchEventEvent{raw: rawEvent}
		}
		close(events)
	}()
	_, err := c.doWatch("GET", eventGeneratePath(namespace, "")+"?"+watchOptionsQuery(opts), nil, rawEvents)
	if err != nil {
		return errors.Wrap(err, "failed to watch Events")
	}
	return nil
}

// DeleteEvent deletes a single Event. It will error if the Event does not exist.
func (c *Client) DeleteEvent(namespace, name string) error {
	_, err := c.do("DELETE", eventGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete Event")
}

// UpdateEvent will update in place a single Event. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdateEvent(namespace string, item *k8s.Event) (*k8s.Event, error) {
	item.TypeMeta.Kind = "Event"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Event
	_, err := c.do("PUT", eventGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Event")
	}
	return &out, nil
}
