package utils_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-online/ocm-common/pkg/ocm/consts"
	. "github.com/openshift-online/ocm-common/pkg/ocm/utils"
)

var _ = Describe("OCM Utils", func() {
	var _ = Describe("Validates CreateVersionId function", func() {
		It("Generates the version ID on stable channel", func() {
			versionId := CreateVersionId("4.10.32", consts.DefaultChannelGroup)
			Expect(versionId).To(Equal("openshift-v4.10.32"))
		})

		It("Generates the version ID on other channel", func() {
			versionId := CreateVersionId("4.10.32", "candidate")
			Expect(versionId).To(Equal("openshift-v4.10.32-candidate"))
		})
	})
})
