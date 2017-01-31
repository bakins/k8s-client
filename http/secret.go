package http

import (
	k8s "github.com/bakins/k8s-client"
	"github.com/pkg/errors"
)

type (
	watchEventSecret struct {
		raw    k8s.WatchEvent
		object *k8s.Secret
	}
)

func (w *watchEventSecret) Type() k8s.WatchEventType {
	return w.raw.Type
}

func (w *watchEventSecret) Object() (*k8s.Secret, error) {
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
	var object k8s.Secret
	if err := w.raw.UnmarshalObject(&object); err != nil {
		return nil, errors.Wrap(err, "failed to decode Secret")
	}
	w.object = &object
	return &object, nil
}

func secretGeneratePath(namespace, name string) string {
	if namespace == "" && name == "" {
		return "/api/v1/secrets"
	}
	if name == "" {
		return "/api/v1/namespaces/" + namespace + "/secrets"
	}
	return "/api/v1/namespaces/" + namespace + "/secrets/" + name
}

// GetSecret fetches a single Secret
func (c *Client) GetSecret(namespace, name string) (*k8s.Secret, error) {
	var out k8s.Secret
	_, err := c.do("GET", secretGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Secret")
	}
	return &out, nil
}

// CreateSecret creates a new Secret. This will fail if it already exists.
func (c *Client) CreateSecret(namespace string, item *k8s.Secret) (*k8s.Secret, error) {
	item.TypeMeta.Kind = "Secret"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Secret
	_, err := c.do("POST", secretGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Secret")
	}
	return &out, nil
}

// ListSecrets lists all Secrets in a namespace
func (c *Client) ListSecrets(namespace string, opts *k8s.ListOptions) (*k8s.SecretList, error) {
	var out k8s.SecretList
	_, err := c.do("GET", secretGeneratePath(namespace, "")+"?"+listOptionsQuery(opts, nil), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list Secrets")
	}
	return &out, nil
}

// WatchSecrets watches all Secret changes in a namespace
func (c *Client) WatchSecrets(namespace string, opts *k8s.WatchOptions, events chan k8s.SecretWatchEvent) error {
	if events == nil {
		return errors.New("events must not be nil")
	}
	rawEvents := make(chan k8s.WatchEvent)
	go func() {
		for rawEvent := range rawEvents {
			events <- &watchEventSecret{raw: rawEvent}
		}
		close(events)
	}()
	_, err := c.doWatch("GET", secretGeneratePath(namespace, "")+"?"+watchOptionsQuery(opts), nil, rawEvents)
	if err != nil {
		return errors.Wrap(err, "failed to watch Secrets")
	}
	return nil
}

// DeleteSecret deletes a single Secret. It will error if the Secret does not exist.
func (c *Client) DeleteSecret(namespace, name string) error {
	_, err := c.do("DELETE", secretGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete Secret")
}

// UpdateSecret will update in place a single Secret. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdateSecret(namespace string, item *k8s.Secret) (*k8s.Secret, error) {
	item.TypeMeta.Kind = "Secret"
	item.TypeMeta.APIVersion = "v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Secret
	_, err := c.do("PUT", secretGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Secret")
	}
	return &out, nil
}
