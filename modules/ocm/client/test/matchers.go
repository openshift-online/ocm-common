package test

// Custom Gomega/GoMock Matchers that make it easier to assert interactions with the OCM API

import v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

const failureMessage = "Expected KubeletConfig Does Not Match"
const negatedFailureMessage = "Expected KubeletConfig Not to Match"

type KubeletConfigMatcher struct {
	expected *v1.KubeletConfig
}

// MatchKubeletConfig returns a Matcher that asserts that the received KubeletConfig
// matches the expected KubeletConfig
func MatchKubeletConfig(expected *v1.KubeletConfig) KubeletConfigMatcher {
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
	return failureMessage
}

func (k KubeletConfigMatcher) Match(actual interface{}) (success bool, err error) {
	return k.Matches(actual), nil
}

func (k KubeletConfigMatcher) FailureMessage(_ interface{}) (message string) {
	return failureMessage
}

func (k KubeletConfigMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return negatedFailureMessage
}
