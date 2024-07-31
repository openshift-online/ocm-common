package validations

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
)

var _ = Describe("ValidateMachinePoolTaintEffect", func() {
	It("should return nil for recognized taint effects", func() {
		Expect(ValidateMachinePoolTaintEffect(v1.TaintEffectNoExecute)).To(Succeed())
		Expect(ValidateMachinePoolTaintEffect(v1.TaintEffectNoSchedule)).To(Succeed())
		Expect(ValidateMachinePoolTaintEffect(v1.TaintEffectPreferNoSchedule)).To(Succeed())
	})

	It("should return an error for unrecognized taint effects", func() {
		unrecognizedEffect := v1.TaintEffect("unrecognized")
		Expect(ValidateMachinePoolTaintEffect(unrecognizedEffect)).To(MatchError(
			MatchRegexp("Unrecognized taint effect 'unrecognized', only the following" +
				" effects are supported: 'NoExecute', 'NoSchedule', 'PreferNoSchedule'")))
	})
})
