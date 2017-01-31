package http

import (
	k8s "github.com/bakins/k8s-client"
	"github.com/pkg/errors"
)

type (
	watchEventIngress struct {
		raw    k8s.WatchEvent
		object *k8s.Ingress
	}
)

func (w *watchEventIngress) Type() k8s.WatchEventType {
	return w.raw.Type
}

func (w *watchEventIngress) Object() (*k8s.Ingress, error) {
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
	var object k8s.Ingress
	if err := w.raw.UnmarshalObject(&object); err != nil {
		return nil, errors.Wrap(err, "failed to decode Ingress")
	}
	w.object = &object
	return &object, nil
}

func ingressGeneratePath(namespace, name string) string {
	if namespace == "" && name == "" {
		return "/apis/extensions/v1beta1/ingresses"
	}
	if name == "" {
		return "/apis/extensions/v1beta1/namespaces/" + namespace + "/ingresses"
	}
	return "/apis/extensions/v1beta1/namespaces/" + namespace + "/ingresses/" + name
}

// GetIngress fetches a single Ingress
func (c *Client) GetIngress(namespace, name string) (*k8s.Ingress, error) {
	var out k8s.Ingress
	_, err := c.do("GET", ingressGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Ingress")
	}
	return &out, nil
}

// CreateIngress creates a new Ingress. This will fail if it already exists.
func (c *Client) CreateIngress(namespace string, item *k8s.Ingress) (*k8s.Ingress, error) {
	item.TypeMeta.Kind = "Ingress"
	item.TypeMeta.APIVersion = "extensions/v1beta1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Ingress
	_, err := c.do("POST", ingressGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Ingress")
	}
	return &out, nil
}

// ListIngresses lists all Ingresss in a namespace
func (c *Client) ListIngresses(namespace string, opts *k8s.ListOptions) (*k8s.IngressList, error) {
	var out k8s.IngressList
	_, err := c.do("GET", ingressGeneratePath(namespace, "")+"?"+listOptionsQuery(opts, nil), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list Ingresss")
	}
	return &out, nil
}

// WatchIngresses watches all Ingress changes in a namespace
func (c *Client) WatchIngresses(namespace string, opts *k8s.WatchOptions, events chan k8s.IngressWatchEvent) error {
	if events == nil {
		return errors.New("events must not be nil")
	}
	rawEvents := make(chan k8s.WatchEvent)
	go func() {
		for rawEvent := range rawEvents {
			events <- &watchEventIngress{raw: rawEvent}
		}
		close(events)
	}()
	_, err := c.doWatch("GET", ingressGeneratePath(namespace, "")+"?"+watchOptionsQuery(opts), nil, rawEvents)
	if err != nil {
		return errors.Wrap(err, "failed to watch Ingresss")
	}
	return nil
}

// DeleteIngress deletes a single Ingress. It will error if the Ingress does not exist.
func (c *Client) DeleteIngress(namespace, name string) error {
	_, err := c.do("DELETE", ingressGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete Ingress")
}

// UpdateIngress will update in place a single Ingress. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdateIngress(namespace string, item *k8s.Ingress) (*k8s.Ingress, error) {
	item.TypeMeta.Kind = "Ingress"
	item.TypeMeta.APIVersion = "extensions/v1beta1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Ingress
	_, err := c.do("PUT", ingressGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Ingress")
	}
	return &out, nil
}
