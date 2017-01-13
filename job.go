package client

type (
	// JobInterface has methods to work with Job resources.
	JobInterface interface {
		CreateJob(namespace string, item *Job) (*Job, error)
		GetJob(namespace, name string) (result *Job, err error)
		ListJobs(namespace string, opts *ListOptions) (*Job, error)
		DeleteJob(namespace, name string) error
		UpdateJob(namespace string, item *Job) (*Job, error)
	}

	// Job represents the configuration of a single job.
	Job struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`

		// Specification of the desired behavior of the Job.
		Spec *JobSpec `json:"spec,omitempty"`

		// Most recently observed status of the DaemonSet.
		Status *JobStatus `json:"status,omitempty"`
	}

	// JobSpec describes how the job execution will look like.
	JobSpec struct {
		// Parallelism specifies the maximum desired number of pods the job should run at any given time.
		// The actual number of pods running in steady state will be less than this number when
		// ((.spec.completions - .status.successful) < .spec.parallelism), i.e. when the work left to do is less than max parallelism.
		Parallelism int32 `json:"parallelism,omitempty"`
		// Completions specifies the desired number of successfully finished pods the job should be run with.
		// Setting to nil means that the success of any pod signals the success of all pods, and allows parallelism to have any positive value.
		// Setting to 1 means that parallelism is limited to 1 and the success of that pod signals the success of the job.
		Completions int32 `json:"completions,omitempty"`
		// Optional duration in seconds relative to the startTime that the job may be active before the system tries to terminate it; value must be positive integer
		ActiveDeadlineSeconds int64 `json:"activeDeadlineSeconds,omitempty"`
		// Selector is a label query over pods that should match the pod count. Normally, the system sets this field for you.
		Selector *LabelSelector `json:"selector,omitempty"`
		// ManualSelector controls generation of pod labels and pod selectors.
		// Leave manualSelector unset unless you are certain what you are doing.
		// When false or unset, the system pick labels unique to this job and appends those labels to the pod template.
		// When true, the user is responsible for picking unique labels and specifying the selector.
		// Failure to pick a unique label may cause this and other jobs to not function correctly.
		// However, You may see manualSelector=true in jobs that were created with the old extensions/v1beta1 API.
		ManualSelector bool `json:"manualSelector,omitempty"`
		// Template is the object that describes the pod that will be created when executing a job.
		Template PodTemplateSpec `json:"template"`
	}

	// JobStatus represents the current status of a job.
	JobStatus struct {
		// Conditions represent the latest available observations of an object’s current state
		Conditions []JobCondition `json:"conditions,omitempty"`
		// StartTime represents time when the job was acknowledged by the Job Manager.
		// It is not guaranteed to be set in happens-before order across separate operations.
		// It is represented in RFC3339 form and is in UTC.
		StartTime Time `json:"startTime,omitempty"`
		// CompletionTime represents time when the job was completed.
		// It is not guaranteed to be set in happens-before order across separate operations.
		// It is represented in RFC3339 form and is in UTC.
		CompletionTime Time `json:"completionTime,omitempty"`
		// Active is the number of actively running pods.
		Active int32 `json:"active,omitempty"`
		// Succeeded is the number of pods which reached Phase Succeeded.
		Succeeded int32 `json:"succeeded,omitempty"`
		// Failed is the number of pods which reached Phase Failed.
		Failed int32 `json:"failed,omitempty"`
	}

	// JobCondition describes current state of a job.
	JobCondition struct {
		// Type of job condition, Complete or Failed.
		Type string `json:"type"`
		// Status of the condition, one of True, False, Unknown.
		Status string `json:"status"`
		// Last time the condition was checked.
		LastProbeTime Time `json:"lastProbeTime,omitempty"`
		// Last time the condition transit from one status to another.
		LastTransitionTime Time `json:"lastTransitionTime,omitempty"`
		// (brief) reason for the condition’s last transition.
		Reason string `json:"reason,omitempty"`
		// Human readable message indicating details about last transition.
		Message string `json:"message,omitempty"`
	}

	JobList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`

		// Items is the list of jobs.
		Items []Job `json:"items"`
	}
)

// NewJob creates a new Job struct
func NewJob(namespace, name string) *Job {
	return &Job{
		TypeMeta:   NewTypeMeta("Job", "batch/v1"),
		ObjectMeta: NewObjectMeta(namespace, name),
		Spec:       &JobSpec{},
	}
}
