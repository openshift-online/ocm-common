package utils_test

import (
	"regexp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openshift-online/ocm-common/pkg/utils"
)

var _ = Describe("Utils", func() {
	var _ = Describe("Validates Random Label function", func() {
		var _ = Context("when generating random labels", func() {

			It("generates empty label given size 0", func() {
				label := RandomLabel(0)
				Expect("").To(Equal(label))
			})

			It("generates random labels of given size", func() {
				for i := 1; i < 11; i++ {
					label := RandomLabel(i)
					Expect(i).To(Equal(len(label)))
					regex, _ := regexp.Compile("^[a-zA-Z0-9]+$")
					Expect(true).To(Equal(regex.MatchString(label)))
				}
			})
		})
	})
})
