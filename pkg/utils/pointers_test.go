package utils_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-online/ocm-common/pkg/utils"
)

var _ = Describe("Pointer utility functions", func() {
	Describe("New function", func() {
		Context("with different primitive types", func() {
			It("should create pointer to int", func() {
				value := 42
				ptr := utils.New(value)
				Expect(ptr).ToNot(BeNil())
				Expect(*ptr).To(Equal(42))
				// Verify it's a different memory location (different pointer address)
				Expect(ptr != &value).To(BeTrue())
			})

			It("should create pointer to int32", func() {
				value := int32(123)
				ptr := utils.New(value)
				Expect(ptr).ToNot(BeNil())
				Expect(*ptr).To(Equal(int32(123)))
			})

			It("should create pointer to int64", func() {
				value := int64(456)
				ptr := utils.New(value)
				Expect(ptr).ToNot(BeNil())
				Expect(*ptr).To(Equal(int64(456)))
			})

			It("should create pointer to float32", func() {
				value := float32(3.14)
				ptr := utils.New(value)
				Expect(ptr).ToNot(BeNil())
				Expect(*ptr).To(Equal(float32(3.14)))
			})

			It("should create pointer to float64", func() {
				value := 2.718
				ptr := utils.New(value)
				Expect(ptr).ToNot(BeNil())
				Expect(*ptr).To(Equal(2.718))
			})

			It("should create pointer to string", func() {
				value := "hello world"
				ptr := utils.New(value)
				Expect(ptr).ToNot(BeNil())
				Expect(*ptr).To(Equal("hello world"))
			})

			It("should create pointer to bool", func() {
				value := true
				ptr := utils.New(value)
				Expect(ptr).ToNot(BeNil())
				Expect(*ptr).To(BeTrue())

				value2 := false
				ptr2 := utils.New(value2)
				Expect(ptr2).ToNot(BeNil())
				Expect(*ptr2).To(BeFalse())
			})

			It("should create pointer to time.Time", func() {
				now := time.Now()
				ptr := utils.New(now)
				Expect(ptr).ToNot(BeNil())
				Expect(*ptr).To(Equal(now))
			})
		})

		Context("with zero values", func() {
			It("should handle zero int", func() {
				ptr := utils.New(0)
				Expect(ptr).ToNot(BeNil())
				Expect(*ptr).To(Equal(0))
			})

			It("should handle empty string", func() {
				ptr := utils.New("")
				Expect(ptr).ToNot(BeNil())
				Expect(*ptr).To(Equal(""))
			})

			It("should handle zero time", func() {
				var zeroTime time.Time
				ptr := utils.New(zeroTime)
				Expect(ptr).ToNot(BeNil())
				Expect(*ptr).To(Equal(zeroTime))
			})
		})
	})

	Describe("NewFormat function", func() {
		It("should create formatted string pointer", func() {
			ptr := utils.NewFormat("Hello %s, you are %d years old", "Alice", 30)
			Expect(ptr).ToNot(BeNil())
			Expect(*ptr).To(Equal("Hello Alice, you are 30 years old"))
		})

		It("should handle format without arguments", func() {
			ptr := utils.NewFormat("Simple string")
			Expect(ptr).ToNot(BeNil())
			Expect(*ptr).To(Equal("Simple string"))
		})

		It("should handle empty format", func() {
			ptr := utils.NewFormat("")
			Expect(ptr).ToNot(BeNil())
			Expect(*ptr).To(Equal(""))
		})

		It("should handle complex formatting", func() {
			ptr := utils.NewFormat("Number: %d, Float: %.2f, Bool: %t", 42, 3.14159, true)
			Expect(ptr).ToNot(BeNil())
			Expect(*ptr).To(Equal("Number: 42, Float: 3.14, Bool: true"))
		})
	})

	Describe("NewStringArray function", func() {
		It("should convert string array to pointer array", func() {
			input := []string{"apple", "banana", "cherry"}
			result := utils.NewStringArray(input)

			Expect(result).ToNot(BeNil())
			Expect(len(result)).To(Equal(3))

			Expect(result[0]).ToNot(BeNil())
			Expect(*result[0]).To(Equal("apple"))

			Expect(result[1]).ToNot(BeNil())
			Expect(*result[1]).To(Equal("banana"))

			Expect(result[2]).ToNot(BeNil())
			Expect(*result[2]).To(Equal("cherry"))
		})

		It("should handle empty array", func() {
			input := []string{}
			result := utils.NewStringArray(input)
			Expect(result).To(BeNil())
		})

		It("should handle nil array", func() {
			result := utils.NewStringArray(nil)
			Expect(result).To(BeNil())
		})

		It("should handle array with empty strings", func() {
			input := []string{"", "test", ""}
			result := utils.NewStringArray(input)

			Expect(result).ToNot(BeNil())
			Expect(len(result)).To(Equal(3))

			Expect(result[0]).ToNot(BeNil())
			Expect(*result[0]).To(Equal(""))

			Expect(result[1]).ToNot(BeNil())
			Expect(*result[1]).To(Equal("test"))

			Expect(result[2]).ToNot(BeNil())
			Expect(*result[2]).To(Equal(""))
		})

		It("should handle single element array", func() {
			input := []string{"single"}
			result := utils.NewStringArray(input)

			Expect(result).ToNot(BeNil())
			Expect(len(result)).To(Equal(1))
			Expect(*result[0]).To(Equal("single"))
		})
	})

	Describe("NewByteArray function", func() {
		It("should copy byte array", func() {
			input := []byte{1, 2, 3, 4, 5}
			result := utils.NewByteArray(input)

			Expect(result).ToNot(BeNil())
			Expect(len(result)).To(Equal(5))
			Expect(result).To(Equal([]byte{1, 2, 3, 4, 5}))

			// Verify it's a copy, not the same slice
			input[0] = 99
			Expect(result[0]).To(Equal(byte(1))) // Should still be 1
		})

		It("should handle empty byte array", func() {
			input := []byte{}
			result := utils.NewByteArray(input)
			Expect(result).To(BeNil())
		})

		It("should handle nil byte array", func() {
			result := utils.NewByteArray(nil)
			Expect(result).To(BeNil())
		})

		It("should handle single byte", func() {
			input := []byte{42}
			result := utils.NewByteArray(input)

			Expect(result).ToNot(BeNil())
			Expect(len(result)).To(Equal(1))
			Expect(result[0]).To(Equal(byte(42)))
		})

		It("should handle byte array with zeros", func() {
			input := []byte{0, 1, 0, 2, 0}
			result := utils.NewByteArray(input)

			Expect(result).ToNot(BeNil())
			Expect(len(result)).To(Equal(5))
			Expect(result).To(Equal([]byte{0, 1, 0, 2, 0}))
		})
	})

	Describe("Memory independence", func() {
		It("should create independent copies for New function", func() {
			original := 100
			ptr1 := utils.New(original)
			ptr2 := utils.New(original)

			// Different pointers (different memory addresses)
			Expect(ptr1 != ptr2).To(BeTrue())

			// Same values
			Expect(*ptr1).To(Equal(*ptr2))
			Expect(*ptr1).To(Equal(original))

			// Modifying original doesn't affect pointers
			//nolint:ineffassign
			original = 200
			Expect(*ptr1).To(Equal(100))
			Expect(*ptr2).To(Equal(100))
		})

		It("should create independent string copies in NewStringArray", func() {
			input := []string{"test"}
			result := utils.NewStringArray(input)

			// Modify original slice
			input[0] = "modified"

			// Result should be unchanged
			Expect(*result[0]).To(Equal("test"))
		})
	})
})
