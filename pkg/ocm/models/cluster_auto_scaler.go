package models

type AutoscalerConfig struct {
	BalanceSimilarNodeGroups    bool
	SkipNodesWithLocalStorage   bool
	LogVerbosity                int
	MaxPodGracePeriod           int
	PodPriorityThreshold        int
	IgnoreDaemonsetsUtilization bool
	MaxNodeProvisionTime        string
	BalancingIgnoredLabels      []string
	ResourceLimits              ResourceLimits
	ScaleDown                   ScaleDownConfig
}

type ResourceLimits struct {
	MaxNodesTotal int
	Cores         ResourceRange
	Memory        ResourceRange
	GPULimits     []GPULimit
}

type ResourceRange struct {
	Min int
	Max int
}

type GPULimit struct {
	Type  string
	Range ResourceRange
}

type ScaleDownConfig struct {
	Enabled              bool
	UnneededTime         string
	UtilizationThreshold float64
	DelayAfterAdd        string
	DelayAfterDelete     string
	DelayAfterFailure    string
}
