package http

import (
	k8s "github.com/bakins/k8s-client"
	"github.com/pkg/errors"
)

type (
	watchEventDeployment struct {
		raw    k8s.WatchEvent
		object *k8s.Deployment
	}
)

func (w *watchEventDeployment) Type() k8s.WatchEventType {
	return w.raw.Type
}

func (w *watchEventDeployment) Object() (*k8s.Deployment, error) {
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
	var object k8s.Deployment
	if err := w.raw.UnmarshalObject(&object); err != nil {
		return nil, errors.Wrap(err, "failed to decode Deployment")
	}
	w.object = &object
	return &object, nil
}

func deploymentGeneratePath(namespace, name string) string {
	if namespace == "" && name == "" {
		return "/apis/extensions/v1beta1/deployments"
	}
	if name == "" {
		return "/apis/extensions/v1beta1/namespaces/" + namespace + "/deployments"
	}
	return "/apis/extensions/v1beta1/namespaces/" + namespace + "/deployments/" + name
}

// GetDeployment fetches a single Deployment
func (c *Client) GetDeployment(namespace, name string) (*k8s.Deployment, error) {
	var out k8s.Deployment
	_, err := c.do("GET", deploymentGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Deployment")
	}
	return &out, nil
}

// CreateDeployment creates a new Deployment. This will fail if it already exists.
func (c *Client) CreateDeployment(namespace string, item *k8s.Deployment) (*k8s.Deployment, error) {
	item.TypeMeta.Kind = "Deployment"
	item.TypeMeta.APIVersion = "extensions/v1beta1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Deployment
	_, err := c.do("POST", deploymentGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Deployment")
	}
	return &out, nil
}

// ListDeployments lists all Deployments in a namespace
func (c *Client) ListDeployments(namespace string, opts *k8s.ListOptions) (*k8s.DeploymentList, error) {
	var out k8s.DeploymentList
	_, err := c.do("GET", deploymentGeneratePath(namespace, "")+"?"+listOptionsQuery(opts, nil), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list Deployments")
	}
	return &out, nil
}

// WatchDeployments watches all Deployment changes in a namespace
func (c *Client) WatchDeployments(namespace string, opts *k8s.WatchOptions, events chan k8s.DeploymentWatchEvent) error {
	if events == nil {
		return errors.New("events must not be nil")
	}
	rawEvents := make(chan k8s.WatchEvent)
	go func() {
		for rawEvent := range rawEvents {
			events <- &watchEventDeployment{raw: rawEvent}
		}
		close(events)
	}()
	_, err := c.doWatch("GET", deploymentGeneratePath(namespace, "")+"?"+watchOptionsQuery(opts), nil, rawEvents)
	if err != nil {
		return errors.Wrap(err, "failed to watch Deployments")
	}
	return nil
}

// DeleteDeployment deletes a single Deployment. It will error if the Deployment does not exist.
func (c *Client) DeleteDeployment(namespace, name string) error {
	_, err := c.do("DELETE", deploymentGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete Deployment")
}

// UpdateDeployment will update in place a single Deployment. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdateDeployment(namespace string, item *k8s.Deployment) (*k8s.Deployment, error) {
	item.TypeMeta.Kind = "Deployment"
	item.TypeMeta.APIVersion = "extensions/v1beta1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Deployment
	_, err := c.do("PUT", deploymentGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Deployment")
	}
	return &out, nil
}
