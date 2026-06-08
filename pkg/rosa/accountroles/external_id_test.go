package accountroles_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openshift-online/ocm-common/pkg/rosa/accountroles"
)

var _ = Describe("RequiresSTSExternalIDInTrustPolicy", func() {
	It("returns true for installer and support", func() {
		Expect(RequiresSTSExternalIDInTrustPolicy(InstallerAccountRole)).To(BeTrue())
		Expect(RequiresSTSExternalIDInTrustPolicy(SupportAccountRole)).To(BeTrue())
	})

	It("returns false for worker and control plane", func() {
		Expect(RequiresSTSExternalIDInTrustPolicy(WorkerAccountRole)).To(BeFalse())
		Expect(RequiresSTSExternalIDInTrustPolicy(ControlPlaneAccountRole)).To(BeFalse())
	})

	It("returns false for unknown role keys", func() {
		Expect(RequiresSTSExternalIDInTrustPolicy("unknown")).To(BeFalse())
	})
})
