package builders

import (
	"fmt"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

	"github.com/openshift-online/ocm-common/pkg/ocm/models"
)

func BuildClusterAutoscaler(config *models.AutoscalerConfig) *cmv1.ClusterAutoscalerBuilder {
	if config == nil {
		return nil
	}

	var gpuLimits []*cmv1.AutoscalerResourceLimitsGPULimitBuilder
	for _, gpuLimit := range config.ResourceLimits.GPULimits {
		gpuLimits = append(
			gpuLimits,
			cmv1.NewAutoscalerResourceLimitsGPULimit().
				Type(gpuLimit.Type).
				Range(cmv1.NewResourceRange().
					Min(gpuLimit.Range.Min).
					Max(gpuLimit.Range.Max)),
		)
	}

	return cmv1.NewClusterAutoscaler().
		BalanceSimilarNodeGroups(config.BalanceSimilarNodeGroups).
		SkipNodesWithLocalStorage(config.SkipNodesWithLocalStorage).
		LogVerbosity(config.LogVerbosity).
		MaxPodGracePeriod(config.MaxPodGracePeriod).
		PodPriorityThreshold(config.PodPriorityThreshold).
		IgnoreDaemonsetsUtilization(config.IgnoreDaemonsetsUtilization).
		MaxNodeProvisionTime(config.MaxNodeProvisionTime).
		BalancingIgnoredLabels(config.BalancingIgnoredLabels...).
		ResourceLimits(cmv1.NewAutoscalerResourceLimits().
			MaxNodesTotal(config.ResourceLimits.MaxNodesTotal).
			Cores(cmv1.NewResourceRange().
				Min(config.ResourceLimits.Cores.Min).
				Max(config.ResourceLimits.Cores.Max)).
			Memory(cmv1.NewResourceRange().
				Min(config.ResourceLimits.Memory.Min).
				Max(config.ResourceLimits.Memory.Max)).
			GPUS(gpuLimits...)).
		ScaleDown(cmv1.NewAutoscalerScaleDownConfig().
			Enabled(config.ScaleDown.Enabled).
			UnneededTime(config.ScaleDown.UnneededTime).
			UtilizationThreshold(fmt.Sprintf("%f", config.ScaleDown.UtilizationThreshold)).
			DelayAfterAdd(config.ScaleDown.DelayAfterAdd).
			DelayAfterDelete(config.ScaleDown.DelayAfterDelete).
			DelayAfterFailure(config.ScaleDown.DelayAfterFailure))
}
