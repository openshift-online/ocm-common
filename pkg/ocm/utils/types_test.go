package utils

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Contains", func() {
	When("the element is present in the slice", func() {
		It("should return true", func() {
			slice := []int{1, 2, 3, 4, 5}
			element := 3
			Expect(Contains(slice, element)).To(BeTrue())
		})

		It("should return true for a custom type", func() {
			type CustomStruct struct {
				Name string
				Age  int
			}

			slice := []CustomStruct{
				{Name: "Alice", Age: 25},
				{Name: "Bob", Age: 30},
			}
			element := CustomStruct{Name: "Alice", Age: 30}
			Expect(Contains(slice, element)).To(BeTrue())
		})
	})

	When("the element is not present in the slice", func() {
		It("should return false", func() {
			slice := []string{"apple", "orange", "banana"}
			element := "grape"
			Expect(Contains(slice, element)).To(BeFalse())
		})

		It("should return false for a custom type", func() {
			type CustomStruct struct {
				Name string
				Age  int
			}

			slice := []CustomStruct{
				{Name: "Alice", Age: 28},
				{Name: "Bob", Age: 35},
			}
			element := CustomStruct{Name: "Charlie", Age: 30}
			Expect(Contains(slice, element)).To(BeFalse())
		})
	})
})
