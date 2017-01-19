package http

import (
	k8s "github.com/YakLabs/k8s-client"
	"github.com/pkg/errors"
)

func endpointsGeneratePath(namespace, name string) string {
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
	_, err := c.do("GET", endpointsGeneratePath(namespace, "")+"?"+listOptionsQuery(opts), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list Endpointss")
	}
	return &out, nil
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
