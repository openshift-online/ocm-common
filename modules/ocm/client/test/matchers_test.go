package test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

var _ = Describe("Matchers Testing", func() {

	Context("KubeletConfigMatcher", func() {

		kubeletConfig := func(k *v1.KubeletConfigBuilder) {
			k.PodPidsLimit(5000)
		}
		kubeletConfig2 := func(k *v1.KubeletConfigBuilder) {
			k.PodPidsLimit(10000)
		}

		It("Matches when KubeletConfig are the same", func() {
			config, err := NewKubeletConfig(kubeletConfig)
			Expect(err).NotTo(HaveOccurred())

			config2, err := NewKubeletConfig(kubeletConfig)
			Expect(err).NotTo(HaveOccurred())

			Expect(config).To(MatchKubeletConfig(config2))
		})

		It("Does not match when KubeletConfigs are different", func() {
			config, err := NewKubeletConfig(kubeletConfig)
			Expect(err).NotTo(HaveOccurred())

			config2, err := NewKubeletConfig(kubeletConfig2)
			Expect(err).NotTo(HaveOccurred())

			Expect(config).NotTo(MatchKubeletConfig(config2))
		})
	})
})
