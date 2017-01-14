package http

import (
	k8s "github.com/YakLabs/k8s-client"
	"github.com/pkg/errors"
)

func jobGeneratePath(namespace, name string) string {
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
	_, err := c.do("GET", jobGeneratePath(namespace, "")+"?"+listOptionsQuery(opts), nil, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list Jobs")
	}
	return &out, nil
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
