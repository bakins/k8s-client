package http

import (
	k8s "github.com/bakins/k8s-client"
	"github.com/pkg/errors"
)

type (
	watchEventNode struct {
		raw    k8s.WatchEvent
		object *k8s.Node
	}
)

func (w *watchEventNode) Type() k8s.WatchEventType {
	return w.raw.Type
}

func (w *watchEventNode) Object() (*k8s.Node, error) {
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
	var object k8s.Node
	if err := w.raw.UnmarshalObject(&object); err != nil {
		return nil, errors.Wrap(err, "failed to decode Node")
	}
	w.object = &object
	return &object, nil
}

// GetNode gets a single node.
func (c *Client) GetNode(name string) (*k8s.Node, error) {
	var out k8s.Node
	_, err := c.do("GET", "/api/v1/nodes/"+name, nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get node")
	}
	return &out, nil
}

// CreateNode creates a single node. It will fail if it already exists.
func (c *Client) CreateNode(item *k8s.Node) (*k8s.Node, error) {
	item.TypeMeta.Kind = "Node"
	item.TypeMeta.APIVersion = "v1"

	var out k8s.Node
	_, err := c.do("POST", "/api/v1/nodes", item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create node")
	}
	return &out, nil
}

// ListNodes list all nodes, optionally filtering.
func (c *Client) ListNodes(opts *k8s.ListOptions) (*k8s.NodeList, error) {
	var out k8s.NodeList
	_, err := c.do("GET", "/api/v1/nodes?"+listOptionsQuery(opts, nil), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list nodes")
	}
	return &out, nil
}

// WatchNodes watches all Nodes changes
func (c *Client) WatchNodes(opts *k8s.WatchOptions, events chan k8s.NodeWatchEvent) error {
	if events == nil {
		return errors.New("events must not be nil")
	}
	rawEvents := make(chan k8s.WatchEvent)
	go func() {
		for rawEvent := range rawEvents {
			events <- &watchEventNode{raw: rawEvent}
		}
		close(events)
	}()
	_, err := c.doWatch("GET", "/api/v1/nodes?"+watchOptionsQuery(opts), nil, rawEvents)
	if err != nil {
		return errors.Wrap(err, "failed to watch Nodes")
	}
	return nil
}

// DeleteNode removes a single node.
func (c *Client) DeleteNode(name string) error {
	_, err := c.do("DELETE", "/api/v1/nodes/"+name, nil, nil)
	return errors.Wrap(err, "failed to delete node")
}

// UpdateNode updates s sinle node.
func (c *Client) UpdateNode(item *k8s.Node) (*k8s.Node, error) {
	item.TypeMeta.Kind = "Node"
	item.TypeMeta.APIVersion = "v1"

	var out k8s.Node
	_, err := c.do("PUT", "/api/v1/nodes/"+item.Name, item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update node")
	}
	return &out, nil
}
