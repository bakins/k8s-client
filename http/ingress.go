package http

import (
	k8s "github.com/YakLabs/k8s-client"
	"github.com/pkg/errors"
)

func ingressGeneratePath(namespace, name string) string {
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

// ListIngresss lists all Ingresss in a namespace
func (c *Client) ListIngresss(namespace string, opts *k8s.ListOptions) (*k8s.IngressList, error) {
	var out k8s.IngressList
	_, err := c.do("GET", ingressGeneratePath(namespace, "")+"?"+listOptionsQuery(opts), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list Ingresss")
	}
	return &out, nil
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
