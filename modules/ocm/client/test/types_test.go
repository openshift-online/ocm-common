package test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

var _ = Describe("OCM Client Types", func() {

	It("Allows customisation of KubeletConfig", func() {
		kubeletConfig, err := NewKubeletConfig(func(k *v1.KubeletConfigBuilder) {
			k.PodPidsLimit(5000)
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(kubeletConfig.PodPidsLimit()).To(Equal(5000))
	})

	It("Allows customisation of MachinePool", func() {
		machinePool, err := NewMachinePool(func(k *v1.MachinePoolBuilder) {
			k.InstanceType("MyInstanceType")
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(machinePool.InstanceType()).To(Equal("MyInstanceType"))
	})

	It("Allows customisation of NodePool", func() {
		nodePool, err := NewNodePool(func(k *v1.NodePoolBuilder) {
			k.AutoRepair(true)
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(nodePool.AutoRepair()).To(BeTrue())
	})

	It("Allows customisation of ClusterAutoscaler", func() {
		autoscaler, err := NewClusterAutoscaler(func(k *v1.ClusterAutoscalerBuilder) {
			k.LogVerbosity(10)
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(autoscaler.LogVerbosity()).To(Equal(10))
	})

})
