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
	var _ = Describe("Validates Truncate function", func() {
		It("Doesn't do anything", func() {
			smallerThanByteInLength := "aaaaa"
			truncated := Truncate(smallerThanByteInLength, 5)
			Expect(smallerThanByteInLength).To(Equal(truncated))
		})
		It("Truncate when bigger than supplied value", func() {
			smallerThanByteInLength := "aaaaaa"
			truncated := Truncate(smallerThanByteInLength, 5)
			Expect(len(truncated)).To(Equal(5))
			Expect(smallerThanByteInLength).To(Equal(truncated + "a"))
		})
	})
	var _ = Describe("Validates GetBoolNotNil function", func() {
		boolPtr := func(b bool) *bool { return &b }

		DescribeTable("returns correct boolean value",
			func(input *bool, expected bool) {
				result := GetBoolNotNil(input)
				Expect(result).To(Equal(expected))
			},
			Entry("nil pointer returns false", nil, false),
			Entry("pointer to true returns true", boolPtr(true), true),
			Entry("pointer to false returns false", boolPtr(false), false),
		)
	})
})
