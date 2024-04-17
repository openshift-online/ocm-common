package utils

import (
	"regexp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	var _ = Describe("Validates Random Label function", func() {
		var _ = Context("when generating random labels", func() {

			It("generates empty label given size 0", func() {
				label := utils.RandomLabel(0)
				Expect("").To(Equal(label))
			})

			It("generates random labels of given size", func() {
				for i := 1; i < 11; i++ {
					label := utils.RandomLabel(i)
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
			truncated := utils.Truncate(smallerThanByteInLength, 5)
			Expect(smallerThanByteInLength).To(Equal(truncated))
		})
		It("Truncate when bigger than supplied value", func() {
			smallerThanByteInLength := "aaaaaa"
			truncated := utils.Truncate(smallerThanByteInLength, 5)
			Expect(len(truncated)).To(Equal(5))
			Expect(smallerThanByteInLength).To(Equal(truncated + "a"))
		})
	})
})
