package client

import "github.com/YakLabs/k8s-client/intstr"

const (
	RestartPolicyAlways    RestartPolicy = "Always"
	RestartPolicyOnFailure RestartPolicy = "OnFailure"
	RestartPolicyNever     RestartPolicy = "Never"
	DNSClusterFirst        DNSPolicy     = "ClusterFirst"
	DNSDefault             DNSPolicy     = "Default"
	// URISchemeHTTP means that the scheme used will be http://
	URISchemeHTTP URIScheme = "HTTP"
	// URISchemeHTTPS means that the scheme used will be https://
	URISchemeHTTPS URIScheme = "HTTPS"
	// CPU, in cores. (500m = .5 cores)
	ResourceCPU ResourceName = "cpu"
	// Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceMemory ResourceName = "memory"
	// Volume size, in bytes (e,g. 5Gi = 5GiB = 5 * 1024 * 1024 * 1024)
	ResourceStorage ResourceName = "storage"
	// PullAlways means that kubelet always attempts to pull the latest image.  Container will fail If the pull fails.
	PullAlways PullPolicy = "Always"
	// PullNever means that kubelet never pulls an image, but only uses a local image.  Container will fail if the image isn't present
	PullNever PullPolicy = "Never"
	// PullIfNotPresent means that kubelet pulls if the image isn't present on disk. Container will fail if the image isn't present and the pull fails.
	PullIfNotPresent PullPolicy = "IfNotPresent"
	// PodPending means the pod has been accepted by the system, but one or more of the containers
	// has not been started. This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	PodPending PodPhase = "Pending"
	// PodRunning means the pod has been bound to a node and all of the containers have been started.
	// At least one container is still running or is in the process of being restarted.
	PodRunning PodPhase = "Running"
	// PodSucceeded means that all containers in the pod have voluntarily terminated
	// with a container exit code of 0, and the system is not going to restart any of these containers.
	PodSucceeded PodPhase = "Succeeded"
	// PodFailed means that all containers in the pod have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	PodFailed PodPhase = "Failed"
	// PodUnknown means that for some reason the state of the pod could not be obtained, typically due
	// to an error in communicating with the host of the pod.
	PodUnknown PodPhase = "Unknown"
)

type (
	PodInterface interface {
		CreatePod(namespace string, item *Pod) (*Pod, error)
		GetPod(namespace, name string) (result *Pod, err error)
		ListPods(namespace string, opts *ListOptions) (*PodList, error)
		DeletePod(namespace, name string) error
		UpdatePod(namespace string, item *Pod) (*Pod, error)
	}

	Pod struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`
		Spec       *PodSpec   `json:"spec,omitempty"`
		Status     *PodStatus `json:"status,omitempty"`
	}

	PodList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`
		Items    []Pod `json:"items"`
	}

	// PodSpec is a description of a pod.
	PodSpec struct {
		// List of volumes that can be mounted by containers belonging to the pod.
		Volumes []Volume `json:"volumes,omitempty"`
		// List of containers belonging to the pod. Containers cannot currently be added or removed.
		// There must be at least one container in a Pod. Cannot be updated.
		Containers []Container `json:"containers"`
		// Restart policy for all containers within the pod. One of Always, OnFailure, Never. Default to Always.
		RestartPolicy RestartPolicy `json:"restartPolicy,omitempty"`
		// Optional duration in seconds the pod needs to terminate gracefully. May be decreased in delete request.
		// Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil,
		// the default grace period will be used instead.
		//  The grace period is the duration in seconds after the processes running in the pod are sent a termination
		// signal and the time when the processes are forcibly halted with a kill signal.
		// Set this value longer than the expected cleanup time for your process. Defaults to 30 seconds.
		TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`
		// Optional duration in seconds the pod may be active on the node relative to StartTime before the system will
		// actively try to mark it failed and kill associated containers. Value must be a positive integer.
		ActiveDeadlineSeconds *int64 `json:"activeDeadlineSeconds,omitempty"`
		// Set DNS policy for containers within the pod. One of ClusterFirst or Default. Defaults to "ClusterFirst".
		DNSPolicy DNSPolicy `json:"dnsPolicy,omitempty"`
		// NodeSelector is a selector which must be true for the pod to fit on a node.
		// Selector which must match a node’s labels for the pod to be scheduled on that node.
		NodeSelector map[string]string `json:"nodeSelector,omitempty"`
		// ServiceAccountName is the name of the ServiceAccount to use to run this pod.
		ServiceAccountName string `json:"serviceAccountName,omitempty"`
		// NodeName is a request to schedule this pod onto a specific node.
		// If it is non-empty, the scheduler simply schedules this pod onto that node, assuming that it fits resource requirements.
		NodeName string `json:"nodeName,omitempty"`
		// Host networking requested for this pod. Use the host’s network namespace.
		// If this option is set, the ports that will be used must be specified. Default to false.
		HostNetwork bool `json:"hostNetwork,omitempty"`
		// Use the host’s pid namespace. Optional: Default to false.
		HostPID bool `json:"hostPID,omitempty"`
		// Use the host’s ipc namespace. Optional: Default to false.
		HostIPC bool `json:"hostIPC,omitempty"`
		// SecurityContext holds pod-level security attributes and common container settings.
		// Optional: Defaults to empty. See type description for default values of each field.
		SecurityContext *PodSecurityContext `json:"securityContext,omitempty"`
		// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
		// If specified, these secrets will be passed to individual puller implementations for them to use.
		// For example, in the case of docker, only DockerConfig type secrets are honored.
		ImagePullSecrets []LocalObjectReference `json:"imagePullSecrets,omitempty"`
		// Specifies the hostname of the Pod If not specified, the pod’s hostname will be set to a system-defined value.
		Hostname string `json:"hostname,omitempty"`
		// If specified, the fully qualified Pod hostname will be "<hostname>.<subdomain>.<pod namespace>.svc.<cluster domain>".
		// If not specified, the pod will not have a domainname at all.
		Subdomain string `json:"subdomain,omitempty"`
	}

	PodStatus struct {
		Phase             PodPhase          `json:"phase,omitempty"`
		Conditions        []PodCondition    `json:"conditions,omitempty"`
		Message           string            `json:"message,omitempty"`
		Reason            string            `json:"reason,omitempty"`
		HostIP            string            `json:"hostIP,omitempty"`
		PodIP             string            `json:"podIP,omitempty"`
		StartTime         *Time             `json:"startTime,omitempty"`
		ContainerStatuses []ContainerStatus `json:"containerStatuses,omitempty"`
	}

	PodCondition struct {
		Type               PodConditionType `json:"type"`
		Status             ConditionStatus  `json:"status"`
		LastProbeTime      Time             `json:"lastProbeTime,omitempty"`
		LastTransitionTime Time             `json:"lastTransitionTime,omitempty"`
		Reason             string           `json:"reason,omitempty"`
		Message            string           `json:"message,omitempty"`
	}

	// PodSecurityContext holds pod-level security attributes and common container settings.
	// Some fields are also present in container.securityContext.
	// Field values of container.securityContext take precedence over field values of PodSecurityContext.
	PodSecurityContext struct {
		// The SELinux context to be applied to all containers. If unspecified, the container runtime will
		// allocate a random SELinux context for each container. May also be set in SecurityContext.
		// If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence for that container.
		SELinuxOptions *SELinuxOptions `json:"seLinuxOptions,omitempty"`
		// The UID to run the entrypoint of the container process.
		// Defaults to user specified in image metadata if unspecified. May also be set in SecurityContext.
		// If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence for that container.
		RunAsUser int64 `json:"runAsUser,omitempty"`
		// Indicates that the container must run as a non-root user. If true, the Kubelet will validate the image at runtime to ensure
		// that it does not run as UID 0 (root) and fail to start the container if it does.
		// If unset or false, no such validation will be performed. May also be set in SecurityContext.
		// If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
		RunAsNonRoot bool `json:"runAsNonRoot,omitempty"`
		// A list of groups applied to the first process run in each container, in addition to the container’s primary GID.
		// If unspecified, no groups will be added to any container.
		SupplementalGroups []int32 `json:"supplementalGroups,omitempty"`
		// A special supplemental group that applies to all containers in a pod.
		// Some volume types allow the Kubelet to change the ownership of that volume to be owned by the pod:
		// 1. The owning GID will be the FSGroup
		// 2. The setgid bit is set (new files created in the volume will be owned by FSGroup)
		// 3. The permission bits are OR’d with rw-rw
		FSGroup int64 `json:"fsGroup,omitempty"`
	}

	// SELinuxOptions are the labels to be applied to the container
	SELinuxOptions struct {
		// User is a SELinux user label that applies to the container.
		User string `json:"user,omitempty"`
		// Role is a SELinux role label that applies to the container.
		Role string `json:"role,omitempty"`
		// Type is a SELinux type label that applies to the container.
		Type string `json:"type,omitempty"`
		// Level is SELinux level label that applies to the container.
		Level string `json:"level,omitempty"`
	}

	PodConditionType string
	StorageMedium    string
	RestartPolicy    string
	DNSPolicy        string
	Protocol         string
	URIScheme        string

	Volume struct {
		Name         string `json:"name"`
		VolumeSource `json:",inline,omitempty"`
	}

	VolumeSource struct {
		EmptyDir *EmptyDirVolumeSource `json:"emptyDir,omitempty"`
		HostPath *HostPathVolumeSource `json:"hostPath,omitempty"`
		Secret   *SecretVolumeSource   `json:"secret,omitempty"`
	}

	// Represents an empty directory for a pod. Empty directory volumes support ownership management and SELinux relabeling.
	EmptyDirVolumeSource struct {
		Medium StorageMedium `json:"medium,omitempty"`
	}

	// Represents a host path mapped into a pod. Host path volumes do not support ownership management or SELinux relabeling.
	HostPathVolumeSource struct {
		// Path of the directory on the host.
		Path string `json:"path"`
	}

	// Adapts a Secret into a volume. The contents of the target Secret’s Data field will be presented in a volume as files
	// using the keys in the Data field as the file names.
	// Secret volumes support ownership management and SELinux relabeling.
	SecretVolumeSource struct {
		// Name of the secret in the pod’s namespace to use.
		SecretName string `json:"secretName,omitempty"`
		// If unspecified, each key-value pair in the Data field of the referenced Secret will be projected into the volume
		// as a file whose name is the key and content is the value.
		// If specified, the listed keys will be projected into the specified paths, and unlisted keys will not be present.
		// If a key is specified which is not present in the Secret, the volume setup will error.
		// Paths must be relative and may not contain the .. path or start with ...
		Items []KeyToPath `json:"items,omitempty"`
		// Optional: mode bits to use on created files by default. Must be a value between 0 and 0777.
		// Defaults to 0644. Directories within the path are not affected by this setting.
		// This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.
		DefaultMode int32 `json:"defaultMode,omitempty"`
	}

	// Maps a string key to a path within a volume.
	KeyToPath struct {
		// The key to project.
		Key string `json:"key"`
		// The relative path of the file to map the key to. May not be an absolute path. May not contain the path element ... May not start with the string ...
		Path string `json:"path"`
		// Optional: mode bits to use on this file, must be a value between 0 and 0777.
		// If not specified, the volume defaultMode will be used. This might be in conflict with other options
		// that affect the file mode, like fsGroup, and the result can be other mode bits set.
		Mode int32 `json:"mode,omitempty"`
	}

	Container struct {
		Name                   string                `json:"name"`
		Image                  string                `json:"image"`
		Command                []string              `json:"command,omitempty"`
		Args                   []string              `json:"args,omitempty"`
		WorkingDir             string                `json:"workingDir,omitempty"`
		Ports                  []ContainerPort       `json:"ports,omitempty"`
		Env                    []EnvVar              `json:"env,omitempty"`
		Resources              *ResourceRequirements `json:"resources,omitempty"`
		VolumeMounts           []VolumeMount         `json:"volumeMounts,omitempty"`
		LivenessProbe          *Probe                `json:"livenessProbe,omitempty"`
		ReadinessProbe         *Probe                `json:"readinessProbe,omitempty"`
		Lifecycle              *Lifecycle            `json:"lifecycle,omitempty"`
		TerminationMessagePath string                `json:"terminationMessagePath,omitempty"`
		ImagePullPolicy        PullPolicy            `json:"imagePullPolicy"`
		Stdin                  bool                  `json:"stdin,omitempty"`
		StdinOnce              bool                  `json:"stdinOnce,omitempty"`
		TTY                    bool                  `json:"tty,omitempty"`
	}

	ContainerPort struct {
		Name          string   `json:"name,omitempty"`
		HostPort      int      `json:"hostPort,omitempty"`
		ContainerPort int      `json:"containerPort"`
		Protocol      Protocol `json:"protocol,omitempty"`
		HostIP        string   `json:"hostIP,omitempty"`
	}

	EnvVar struct {
		Name      string        `json:"name"`
		Value     string        `json:"value,omitempty"`
		ValueFrom *EnvVarSource `json:"valueFrom,omitempty"`
	}

	EnvVarSource struct {
		FieldRef        *ObjectFieldSelector  `json:"fieldRef,omitempty"`
		ConfigMapKeyRef *ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`
		SecretKeyRef    *SecretKeySelector    `json:"secretKeyRef,omitempty"`
	}

	// VolumeMount describes a mounting of a Volume within a container.
	VolumeMount struct {
		// Required: This must match the Name of a Volume [above].
		Name string `json:"name"`
		// Optional: Defaults to false (read-write).
		ReadOnly bool `json:"readOnly,omitempty"`
		// Required. Must not contain ':'.
		MountPath string `json:"mountPath"`
	}

	// Probe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.
	Probe struct {
		// The action taken to determine the health of a container
		Handler `json:",inline"`
		// Length of time before health checking is activated.  In seconds.
		InitialDelaySeconds int `json:"initialDelaySeconds,omitempty"`
		// Length of time before health checking times out.  In seconds.
		TimeoutSeconds int `json:"timeoutSeconds,omitempty"`
		// How often (in seconds) to perform the probe.
		PeriodSeconds int `json:"periodSeconds,omitempty"`
		// Minimum consecutive successes for the probe to be considered successful after having failed.
		// Must be 1 for liveness.
		SuccessThreshold int `json:"successThreshold,omitempty"`
		// Minimum consecutive failures for the probe to be considered failed after having succeeded.
		FailureThreshold int `json:"failureThreshold,omitempty"`
	}

	// Handler defines a specific action that should be taken TODO: pass structured data to these actions, and document that data here.
	Handler struct {
		// One and only one of the following should be specified.
		// Exec specifies the action to take.
		Exec *ExecAction `json:"exec,omitempty"`
		// HTTPGet specifies the http request to perform.
		HTTPGet *HTTPGetAction `json:"httpGet,omitempty"`
		// TCPSocket specifies an action involving a TCP port.
		// TODO: implement a realistic TCP lifecycle hook
		TCPSocket *TCPSocketAction `json:"tcpSocket,omitempty"`
	}

	// ExecAction describes a "run in container" action.
	ExecAction struct {
		// Command is the command line to execute inside the container, the working directory for the
		// command  is root ('/') in the container's filesystem.  The command is simply exec'd, it is
		// not run inside a shell, so traditional shell instructions ('|', etc) won't work.  To use
		// a shell, you need to explicitly call out to that shell.
		Command []string `json:"command,omitempty"`
	}

	// HTTPGetAction describes an action based on HTTP Get requests.
	HTTPGetAction struct {
		// Optional: Path to access on the HTTP server.
		Path string `json:"path,omitempty"`
		// Required: Name or number of the port to access on the container.
		Port intstr.IntOrString `json:"port,omitempty"`
		// Optional: Host name to connect to, defaults to the pod IP. You
		// probably want to set "Host" in httpHeaders instead.
		Host string `json:"host,omitempty"`
		// Optional: Scheme to use for connecting to the host, defaults to HTTP.
		Scheme URIScheme `json:"scheme,omitempty"`
		// Optional: Custom headers to set in the request. HTTP allows repeated headers.
		HTTPHeaders []HTTPHeader `json:"httpHeaders,omitempty"`
	}

	// HTTPHeader describes a custom header to be used in HTTP probes
	HTTPHeader struct {
		// The header field name
		Name string `json:"name"`
		// The header field value
		Value string `json:"value"`
	}

	// ResourceRequirements describes the compute resource requirements.
	ResourceRequirements struct {
		// Limits describes the maximum amount of compute resources allowed.
		Limits ResourceList `json:"limits,omitempty"`
		// Requests describes the minimum amount of compute resources required.
		// If Request is omitted for a container, it defaults to Limits if that is explicitly specified,
		// otherwise to an implementation-defined value
		Requests ResourceList `json:"requests,omitempty"`
	}

	// ResourceList is a set of (resource name, quantity) pairs.
	ResourceList map[ResourceName]string

	// ResourceName is the name identifying various resources in a ResourceList.
	ResourceName string

	// TCPSocketAction describes an action based on opening a socket
	TCPSocketAction struct {
		// Required: Port to connect to.
		Port intstr.IntOrString `json:"port,omitempty"`
	}

	// Lifecycle describes actions that the management system should take in response to container lifecycle events. For the PostStart and PreStop lifecycle handlers, management of the container blocks until the action is complete, unless the container process fails, in which case the handler is aborted.
	Lifecycle struct {
		// PostStart is called immediately after a container is created.  If the handler fails, the container
		// is terminated and restarted.
		PostStart *Handler `json:"postStart,omitempty"`
		// PreStop is called immediately before a container is terminated.  The reason for termination is
		// passed to the handler.  Regardless of the outcome of the handler, the container is eventually terminated.
		PreStop *Handler `json:"preStop,omitempty"`
	}

	// PullPolicy describes a policy for if/when to pull a container image
	PullPolicy string

	// PodPhase is a label for the condition of a pod at the current time.
	PodPhase string

	// ContainerStatus is the status of a single container within a pod.
	ContainerStatus struct {
		// Each container in a pod must have a unique name.
		Name                 string          `json:"name"`
		State                *ContainerState `json:"state,omitempty"`
		LastTerminationState *ContainerState `json:"lastState,omitempty"`
		// Ready specifies whether the container has passed its readiness check.
		Ready bool `json:"ready"`
		// Note that this is calculated from dead containers.  But those containers are subject to
		// garbage collection.  This value will get capped at 5 by GC.
		RestartCount int    `json:"restartCount"`
		Image        string `json:"image"`
		ImageID      string `json:"imageID"`
		ContainerID  string `json:"containerID,omitempty"`
	}

	// ContainerState holds a possible state of container. Only one of its members may be specified. If none of them is specified, the default one is ContainerStateWaiting.
	ContainerState struct {
		Waiting    *ContainerStateWaiting    `json:"waiting,omitempty"`
		Running    *ContainerStateRunning    `json:"running,omitempty"`
		Terminated *ContainerStateTerminated `json:"terminated,omitempty"`
	}

	ContainerStateWaiting struct {
		// A brief CamelCase string indicating details about why the container is in waiting state.
		Reason string `json:"reason,omitempty"`
		// A human-readable message indicating details about why the container is in waiting state.
		Message string `json:"message,omitempty"`
	}

	ContainerStateRunning struct {
		StartedAt Time `json:"startedAt,omitempty"`
	}

	ContainerStateTerminated struct {
		ExitCode    int    `json:"exitCode"`
		Signal      int    `json:"signal,omitempty"`
		Reason      string `json:"reason,omitempty"`
		Message     string `json:"message,omitempty"`
		StartedAt   Time   `json:"startedAt,omitempty"`
		FinishedAt  Time   `json:"finishedAt,omitempty"`
		ContainerID string `json:"containerID,omitempty"`
	}
)
