package http

import (
	k8s "github.com/YakLabs/k8s-client"
	"github.com/pkg/errors"
)

func daemonsetGeneratePath(namespace, name string) string {
	if name == "" {
		return "/apis/extensions/v1beta1/namespaces/" + namespace + "/daemonsets"
	}
	return "/apis/extensions/v1beta1/namespaces/" + namespace + "/daemonsets/" + name
}

// GetDaemonSet fetches a single DaemonSet
func (c *Client) GetDaemonSet(namespace, name string) (*k8s.DaemonSet, error) {
	var out k8s.DaemonSet
	_, err := c.do("GET", daemonsetGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get DaemonSet")
	}
	return &out, nil
}

// CreateDaemonSet creates a new DaemonSet. This will fail if it already exists.
func (c *Client) CreateDaemonSet(namespace string, item *k8s.DaemonSet) (*k8s.DaemonSet, error) {
	item.TypeMeta.Kind = "DaemonSet"
	item.TypeMeta.APIVersion = "extensions/v1beta1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.DaemonSet
	_, err := c.do("POST", daemonsetGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create DaemonSet")
	}
	return &out, nil
}

// ListDaemonSets lists all DaemonSets in a namespace
func (c *Client) ListDaemonSets(namespace string, opts *k8s.ListOptions) (*k8s.DaemonSetList, error) {
	var out k8s.DaemonSetList
	_, err := c.do("GET", daemonsetGeneratePath(namespace, "")+"?"+listOptionsQuery(opts), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list DaemonSets")
	}
	return &out, nil
}

// DeleteDaemonSet deletes a single DaemonSet. It will error if the DaemonSet does not exist.
func (c *Client) DeleteDaemonSet(namespace, name string) error {
	_, err := c.do("DELETE", daemonsetGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete DaemonSet")
}

// UpdateDaemonSet will update in place a single DaemonSet. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdateDaemonSet(namespace string, item *k8s.DaemonSet) (*k8s.DaemonSet, error) {
	item.TypeMeta.Kind = "DaemonSet"
	item.TypeMeta.APIVersion = "extensions/v1beta1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.DaemonSet
	_, err := c.do("PUT", daemonsetGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update DaemonSet")
	}
	return &out, nil
}
