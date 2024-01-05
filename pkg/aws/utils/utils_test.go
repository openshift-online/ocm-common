package utils_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openshift-online/ocm-common/pkg/aws/utils"
)

var _ = Describe("AWS Utils", func() {
	var _ = Describe("Validates GetPathFromArn function", func() {
		It("Gets the path from arn", func() {
			path, err := GetPathFromArn(
				"arn:partition:service:region:account-id:resource-type/test-path/resource-id")
			Expect(err).ToNot(HaveOccurred())
			Expect(path).To(Equal("/test-path/"))
		})
		It("Retrieves empty when there's no path", func() {
			path, err := GetPathFromArn(
				"arn:partition:service:region:account-id:resource-type/resource-id")
			Expect(err).ToNot(HaveOccurred())
			Expect(path).To(Equal(""))
		})
		It("Errors if arn isn't valid", func() {
			_, err := GetPathFromArn("aaaa")
			Expect(err).To(HaveOccurred())
		})
	})
})
