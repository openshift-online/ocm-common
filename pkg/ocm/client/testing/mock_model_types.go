package testing

import v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

// MockKubeletConfig creates an empty KubeletConfig that can be used in tests. Tests can
// mutate the KubeletConfig to their requirements via the variadic list of functions
func MockKubeletConfig(modifyFn ...func(k *v1.KubeletConfig)) (*v1.KubeletConfig, error) {
	builder := v1.KubeletConfigBuilder{}
	kubeletConfig, err := builder.Build()
	if err != nil {
		return nil, err
	}

	for _, f := range modifyFn {
		f(kubeletConfig)
	}

	return kubeletConfig, nil
}
