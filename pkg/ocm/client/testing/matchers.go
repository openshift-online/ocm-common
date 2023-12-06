package testing

// Custom Gomega Matchers that make it easier to assert interactions with the OCM API

import v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

type KubeletConfigMatcher struct {
	expected *v1.KubeletConfig
}

// KubeletConfigMatches returns a Matcher that asserts that the received KubeletConfig
// matches the expected KubeletConfig
func KubeletConfigMatches(expected *v1.KubeletConfig) KubeletConfigMatcher {
	return KubeletConfigMatcher{
		expected: expected,
	}
}

func (k KubeletConfigMatcher) Matches(x interface{}) bool {
	if kubeletConfig, ok := x.(*v1.KubeletConfig); ok {
		return k.expected.PodPidsLimit() == kubeletConfig.PodPidsLimit()
	}
	return false
}

func (k KubeletConfigMatcher) String() string {
	return "Expected KubeletConfig Does Not Match"
}
