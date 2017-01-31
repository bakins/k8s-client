package http

import (
	k8s "github.com/bakins/k8s-client"
	"github.com/pkg/errors"
)

type (
	watchEventHorizontalPodAutoscaler struct {
		raw    k8s.WatchEvent
		object *k8s.HorizontalPodAutoscaler
	}
)

func (w *watchEventHorizontalPodAutoscaler) Type() k8s.WatchEventType {
	return w.raw.Type
}

func (w *watchEventHorizontalPodAutoscaler) Object() (*k8s.HorizontalPodAutoscaler, error) {
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
	var object k8s.HorizontalPodAutoscaler
	if err := w.raw.UnmarshalObject(&object); err != nil {
		return nil, errors.Wrap(err, "failed to decode HorizontalPodAutoscaler")
	}
	w.object = &object
	return &object, nil
}

func horizontalpodautoscalerGeneratePath(namespace, name string) string {
	if namespace == "" && name == "" {
		return "/apis/autoscaling/v1/horizontalpodautoscalers"
	}
	if name == "" {
		return "/apis/autoscaling/v1/namespaces/" + namespace + "/horizontalpodautoscalers"
	}
	return "/apis/autoscaling/v1/namespaces/" + namespace + "/horizontalpodautoscalers/" + name
}

// GetHorizontalPodAutoscaler fetches a single HorizontalPodAutoscaler
func (c *Client) GetHorizontalPodAutoscaler(namespace, name string) (*k8s.HorizontalPodAutoscaler, error) {
	var out k8s.HorizontalPodAutoscaler
	_, err := c.do("GET", horizontalpodautoscalerGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get HorizontalPodAutoscaler")
	}
	return &out, nil
}

// CreateHorizontalPodAutoscaler creates a new HorizontalPodAutoscaler. This will fail if it already exists.
func (c *Client) CreateHorizontalPodAutoscaler(namespace string, item *k8s.HorizontalPodAutoscaler) (*k8s.HorizontalPodAutoscaler, error) {
	item.TypeMeta.Kind = "HorizontalPodAutoscaler"
	item.TypeMeta.APIVersion = "autoscaling/v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.HorizontalPodAutoscaler
	_, err := c.do("POST", horizontalpodautoscalerGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HorizontalPodAutoscaler")
	}
	return &out, nil
}

// ListHorizontalPodAutoscalers lists all HorizontalPodAutoscalers in a namespace
func (c *Client) ListHorizontalPodAutoscalers(namespace string, opts *k8s.ListOptions) (*k8s.HorizontalPodAutoscalerList, error) {
	var out k8s.HorizontalPodAutoscalerList
	_, err := c.do("GET", horizontalpodautoscalerGeneratePath(namespace, "")+"?"+listOptionsQuery(opts, nil), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list HorizontalPodAutoscalers")
	}
	return &out, nil
}

// WatchHorizontalPodAutoscalers watches all HorizontalPodAutoscaler changes in a namespace
func (c *Client) WatchHorizontalPodAutoscalers(namespace string, opts *k8s.WatchOptions, events chan k8s.HorizontalPodAutoscalerWatchEvent) error {
	if events == nil {
		return errors.New("events must not be nil")
	}
	rawEvents := make(chan k8s.WatchEvent)
	go func() {
		for rawEvent := range rawEvents {
			events <- &watchEventHorizontalPodAutoscaler{raw: rawEvent}
		}
		close(events)
	}()
	_, err := c.doWatch("GET", horizontalpodautoscalerGeneratePath(namespace, "")+"?"+watchOptionsQuery(opts), nil, rawEvents)
	if err != nil {
		return errors.Wrap(err, "failed to watch HorizontalPodAutoscalers")
	}
	return nil
}

// DeleteHorizontalPodAutoscaler deletes a single HorizontalPodAutoscaler. It will error if the HorizontalPodAutoscaler does not exist.
func (c *Client) DeleteHorizontalPodAutoscaler(namespace, name string) error {
	_, err := c.do("DELETE", horizontalpodautoscalerGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete HorizontalPodAutoscaler")
}

// UpdateHorizontalPodAutoscaler will update in place a single HorizontalPodAutoscaler. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdateHorizontalPodAutoscaler(namespace string, item *k8s.HorizontalPodAutoscaler) (*k8s.HorizontalPodAutoscaler, error) {
	item.TypeMeta.Kind = "HorizontalPodAutoscaler"
	item.TypeMeta.APIVersion = "autoscaling/v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.HorizontalPodAutoscaler
	_, err := c.do("PUT", horizontalpodautoscalerGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update HorizontalPodAutoscaler")
	}
	return &out, nil
}
