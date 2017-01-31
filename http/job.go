package http

import (
	k8s "github.com/bakins/k8s-client"
	"github.com/pkg/errors"
)

type (
	watchEventJob struct {
		raw    k8s.WatchEvent
		object *k8s.Job
	}
)

func (w *watchEventJob) Type() k8s.WatchEventType {
	return w.raw.Type
}

func (w *watchEventJob) Object() (*k8s.Job, error) {
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
	var object k8s.Job
	if err := w.raw.UnmarshalObject(&object); err != nil {
		return nil, errors.Wrap(err, "failed to decode Job")
	}
	w.object = &object
	return &object, nil
}

func jobGeneratePath(namespace, name string) string {
	if namespace == "" && name == "" {
		return "/apis/batch/v1/jobs"
	}
	if name == "" {
		return "/apis/batch/v1/namespaces/" + namespace + "/jobs"
	}
	return "/apis/batch/v1/namespaces/" + namespace + "/jobs/" + name
}

// GetJob fetches a single Job
func (c *Client) GetJob(namespace, name string) (*k8s.Job, error) {
	var out k8s.Job
	_, err := c.do("GET", jobGeneratePath(namespace, name), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Job")
	}
	return &out, nil
}

// CreateJob creates a new Job. This will fail if it already exists.
func (c *Client) CreateJob(namespace string, item *k8s.Job) (*k8s.Job, error) {
	item.TypeMeta.Kind = "Job"
	item.TypeMeta.APIVersion = "batch/v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Job
	_, err := c.do("POST", jobGeneratePath(namespace, ""), item, &out, 201)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Job")
	}
	return &out, nil
}

// ListJobs lists all Jobs in a namespace
func (c *Client) ListJobs(namespace string, opts *k8s.ListOptions) (*k8s.JobList, error) {
	var out k8s.JobList
	_, err := c.do("GET", jobGeneratePath(namespace, "")+"?"+listOptionsQuery(opts, nil), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list Jobs")
	}
	return &out, nil
}

// WatchJobs watches all Job changes in a namespace
func (c *Client) WatchJobs(namespace string, opts *k8s.WatchOptions, events chan k8s.JobWatchEvent) error {
	if events == nil {
		return errors.New("events must not be nil")
	}
	rawEvents := make(chan k8s.WatchEvent)
	go func() {
		for rawEvent := range rawEvents {
			events <- &watchEventJob{raw: rawEvent}
		}
		close(events)
	}()
	_, err := c.doWatch("GET", jobGeneratePath(namespace, "")+"?"+watchOptionsQuery(opts), nil, rawEvents)
	if err != nil {
		return errors.Wrap(err, "failed to watch Jobs")
	}
	return nil
}

// DeleteJob deletes a single Job. It will error if the Job does not exist.
func (c *Client) DeleteJob(namespace, name string) error {
	_, err := c.do("DELETE", jobGeneratePath(namespace, name), nil, nil)
	return errors.Wrap(err, "failed to delete Job")
}

// UpdateJob will update in place a single Job. Generally, you should call
// Get and then use that object for updates to ensure resource versions
// avoid update conflicts
func (c *Client) UpdateJob(namespace string, item *k8s.Job) (*k8s.Job, error) {
	item.TypeMeta.Kind = "Job"
	item.TypeMeta.APIVersion = "batch/v1"
	item.ObjectMeta.Namespace = namespace

	var out k8s.Job
	_, err := c.do("PUT", jobGeneratePath(namespace, item.Name), item, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Job")
	}
	return &out, nil
}
