package v1beta

import (
	"github.com/urvil38/kubepaas/schema/util"
)

const (
	Version = "kubepaas/v1beta"
)

type KubepaasConfig struct {
	APIVersion string `json:"apiVersion"`

	// kind represent the type of config.
	// Value is always going to be "config".
	Kind string `json:"kind"`

	Metadata metadata `json:"metadata,omitempty"`

	Build BuildConfig `json:"build,omitempty"`

	Deploy DeployConfig `json:"deploy,omitempty"`
}

func (c *KubepaasConfig) GetVersion() string {
	return c.APIVersion
}

func NewKubePaasConfig() util.VersionedConfig {
	return new(KubepaasConfig)
}

type metadata struct {
	Name        string            `json:"name,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type BuildConfig struct {
	TagPolicy TagPolicy `json:"tagPolicy,omitempty"`
}

type TagPolicy struct {
	GitTagger      *GitTagger      `json:"gitCommit,omitempty"`
	ShaTagger      *ShaTagger      `json:"sha256,omitempty"`
	DateTimeTagger *DateTimeTagger `json:"dateTime,omitempty"`
}

type GitTagger struct{}

type ShaTagger struct{}

type DateTimeTagger struct {
	Format   string `json:"format,omitempty"`
	TimeZone string `json:"timezone,omitempty"`
}

type DeployConfig struct {

	// Runtime on which current deployment is depends on.
	// i.e. runtime can be node.js, python, go, java etc.
	Runtime string `json:"runtime,omitempty"`

	// Port on which deployment is running
	Port string `json:"port,omitempty"`

	// If provided, this docker image will be used for deployment
	Image string `json:"image,omitempty"`

	StaticDir string `json:"static_dir,omitempty"`

	// Relative path to the dockerfile, which will be used to build the artifect(docker image)
	// If Not provided, then default template of particular runtime will be used to generate
	// dockerfile.
	DockerfilePath string `json:"dockerfilePath,omitempty"`

	Resources *ResourceRequirements `json:"resources,omitempty"`

	Envs []EnvVar `json:"env,omitempty"`

	ExternalLogging bool `json:"logging,omitempty"`

	// if true, app must expose the "/metrics" path to be collected by the
	// prometheus metrics collector.
	Metrics MetricsConfig `json:"matrics,omitempty"`

	LivenessProbe *Probe `json:"livenessProbe,omitempty"`

	ReadinessProbe *Probe `json:"readinessProbe,omitempty"`
}

type Probe struct {
	// The action taken to determine the health of a container
	Path string `json:"path,omitempty"`
	// Number of seconds after the container has started before liveness probes are initiated.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	InitialDelaySeconds int32 `json:"initialDelaySeconds,omitempty"`
	// Number of seconds after which the probe times out.
	// Defaults to 1 second. Minimum value is 1.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty"`
	// How often (in seconds) to perform the probe.
	// Default to 10 seconds. Minimum value is 1.
	// +optional
	PeriodSeconds int32 `json:"periodSeconds,omitempty"`
	// Minimum consecutive successes for the probe to be considered successful after having failed.
	// Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
	// +optional
	SuccessThreshold int32 `json:"successThreshold,omitempty"`
	// Minimum consecutive failures for the probe to be considered failed after having succeeded.
	// Defaults to 3. Minimum value is 1.
	// +optional
	FailureThreshold int32 `json:"failureThreshold,omitempty"`
}

// ResourceRequirements describes the resource requirements for the kaniko pod.
type ResourceRequirements struct {
	// Requests [resource requests](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container) for the Kaniko pod.
	Requests *ResourceRequirement `json:"requests,omitempty"`

	// Limits [resource limits](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container) for the Kaniko pod.
	Limits *ResourceRequirement `json:"limits,omitempty"`
}

// ResourceRequirement stores the CPU/Memory requirements for the pod.
type ResourceRequirement struct {
	// CPU the number cores to be used.
	// For example: `2`, `2.0` or `200m`.
	CPU string `json:"cpu,omitempty"`

	// Memory the amount of memory to allocate to the pod.
	// For example: `1Gi` or `1000Mi`.
	Memory string `json:"memory,omitempty"`

	// EphemeralStorage the amount of Ephemeral storage to allocate to the pod.
	// For example: `1Gi` or `1000Mi`.
	EphemeralStorage string `json:"ephemeralStorage,omitempty"`

	// ResourceStorage the amount of resource storage to allocate to the pod.
	// For example: `1Gi` or `1000Mi`.
	ResourceStorage string `json:"resourceStorage,omitempty"`
}

type EnvVar struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type MetricsConfig struct {
	Path            string   `json:"path,omitempty"`
	SrapingInterval string   `json:"scraping_interval,omitempty"`
	HTTPAuth        HTTPAuth `json:"auth,omitempty"`
}

type HTTPAuth struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
