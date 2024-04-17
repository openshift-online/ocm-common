package validations

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cluster Node Validations", func() {
	Context("MinReplicasValidator", func() {
		It("minReplicas negative (failure)", func() {
			err := MinReplicasValidator(-1, false, false, 0)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("The value for the number of cluster nodes must be non-negative"))
		})
		It("hostedCP privateSubnetsCount zero (failure)", func() {
			err := MinReplicasValidator(1, false, true, 0)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Hosted clusters require at least a private subnet"))
		})
		It("hostedCP compute nodes not multiple of private subnets (failure)", func() {
			err := MinReplicasValidator(3, false, true, 2)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Hosted clusters require that the number of compute nodes be a multiple of " +
				"the number of private subnets 2, instead received: 3"))
		})
		It("MultipleAZ minReplicas smaller than 3 (failure)", func() {
			err := MinReplicasValidator(2, true, false, 0)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Multi AZ cluster requires at least 3 compute nodes"))
		})
		It("MultipleAZ minReplicas not multiple of 3 (failure)", func() {
			err := MinReplicasValidator(4, true, false, 0)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Multi AZ clusters require that the number of compute nodes be a multiple of 3"))
		})
		It("minReplicas zero (failure)", func() {
			err := MinReplicasValidator(0, false, false, 0)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Cluster requires at least 2 compute nodes"))
		})
		It("not MultipleAZ minReplicas smaller than 2 (failure)", func() {
			err := MinReplicasValidator(1, false, false, 0)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Cluster requires at least 2 compute nodes"))
		})
		It("MultipleAZ minReplicas multiple of 3 (success)", func() {
			err := MinReplicasValidator(3, true, false, 0)
			Expect(err).NotTo(HaveOccurred())
		})
		It("not MultipleAZ minReplicas 2 (success)", func() {
			err := MinReplicasValidator(2, false, false, 0)
			Expect(err).NotTo(HaveOccurred())
		})
		It("hostedCP compute nodes multiple of private subnets (success)", func() {
			err := MinReplicasValidator(4, false, true, 2)
			Expect(err).NotTo(HaveOccurred())
		})
		It("Not hostedCP compute nodes count is 3 (success)", func() {
			err := MinReplicasValidator(3, true, false, 0)
			Expect(err).NotTo(HaveOccurred())
		})

	})
	Context("MaxReplicasValidator", func() {
		It("maxReplicas smaller than minReplicas (failure)", func() {
			err := MaxReplicasValidator(4, 3, false, false, 0)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("max-replicas must be greater or equal to min-replicas"))
		})
		It("hostedCP compute nodes not multiple of private subnets (failure)", func() {
			err := MaxReplicasValidator(2, 3, false, true, 2)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Hosted clusters require that the number of compute nodes be a multiple of " +
				"the number of private subnets 2, instead received: 3"))
		})
		It("hostedCP compute nodes multiple of private subnets (success)", func() {
			err := MaxReplicasValidator(2, 4, false, true, 2)
			Expect(err).NotTo(HaveOccurred())
		})
		It("MultipleAZ maxReplicas not multiple of 3 (failure)", func() {
			err := MaxReplicasValidator(3, 4, true, false, 0)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Multi AZ clusters require that the number of compute nodes be a multiple of 3"))
		})
		It("MultipleAZ maxReplicas multiple of 3 (success)", func() {
			err := MaxReplicasValidator(3, 6, true, false, 0)
			Expect(err).NotTo(HaveOccurred())
		})
		It("Not MultipleAZ and not hostedCP maxReplicas bigger than minReplicas (success)", func() {
			err := MaxReplicasValidator(2, 5, false, false, 0)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("ValidateAvailabilityZonesCount", func() {
		It("Not MultipleAZ (success)", func() {
			err := ValidateAvailabilityZonesCount(false, 1)
			Expect(err).NotTo(HaveOccurred())
		})
		It("Not MultipleAZ (failure)", func() {
			err := ValidateAvailabilityZonesCount(false, 2)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("The number of availability zones for " +
				"a single AZ cluster should be 1, instead received: 2"))
		})
		It("MultipleAZ (success)", func() {
			err := ValidateAvailabilityZonesCount(true, 3)
			Expect(err).NotTo(HaveOccurred())
		})
		It("MultipleAZ (failure)", func() {
			err := ValidateAvailabilityZonesCount(true, 2)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("The number of availability zones for " +
				"a multi AZ cluster should be 3, instead received: 2"))
		})
	})
})
