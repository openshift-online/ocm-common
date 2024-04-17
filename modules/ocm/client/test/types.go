package test

import v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

// NewKubeletConfig creates an empty KubeletConfig that can be used in tests. Tests can
// mutate the KubeletConfig to their requirements via the variadic list of functions
func NewKubeletConfig(modifyFn ...func(k *v1.KubeletConfigBuilder)) (*v1.KubeletConfig, error) {
	builder := &v1.KubeletConfigBuilder{}
	for _, f := range modifyFn {
		f(builder)
	}
	return builder.Build()
}

// NewClusterAutoscaler creates an empty ClusterAutoscaler that can be used in tests. Tests can
// mutate the ClusterAutoscaler to their requirements via the variadic list of functions
func NewClusterAutoscaler(modifyFn ...func(k *v1.ClusterAutoscalerBuilder)) (*v1.ClusterAutoscaler, error) {
	builder := &v1.ClusterAutoscalerBuilder{}
	for _, f := range modifyFn {
		f(builder)
	}
	return builder.Build()
}

// NewMachinePool creates an empty MachinePool that can be used in tests. Tests can
// mutate the MachinePool to their requirements via the variadic list of functions
func NewMachinePool(modifyFn ...func(k *v1.MachinePoolBuilder)) (*v1.MachinePool, error) {
	builder := &v1.MachinePoolBuilder{}
	for _, f := range modifyFn {
		f(builder)
	}
	return builder.Build()
}

// NewNodePool creates an empty NodePool that can be used in tests. Tests can
// mutate the NodePool to their requirements via the variadic list of functions
func NewNodePool(modifyFn ...func(k *v1.NodePoolBuilder)) (*v1.NodePool, error) {
	builder := &v1.NodePoolBuilder{}
	for _, f := range modifyFn {
		f(builder)
	}
	return builder.Build()
}
