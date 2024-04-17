package validations

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validations", func() {
	Describe("validateKMSKeyARN", func() {
		Context("empty kmsKeyARN", func() {
			When("kmsKeyARN is nil", func() {
				It("should not return an error", func() {
					err := ValidateKMSKeyARN(nil)
					Expect(err).ToNot(HaveOccurred())
				})
			})

			When("kmsKeyARN is empty", func() {
				It("should not return an error", func() {
					emptyKmsKeyARN := ""
					err := ValidateKMSKeyARN(&emptyKmsKeyARN)
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})

		Context("kmsKeyARN regex", func() {
			When("kmsKeyARN is not empty and matches the regex", func() {
				It("should not return an error", func() {
					validKmsKeyARN := "arn:aws:kms:us-east-1:111111111111:key/mrk-0123456789abcdef0123456789abcdef"
					err := ValidateKMSKeyARN(&validKmsKeyARN)
					Expect(err).ToNot(HaveOccurred())
				})
			})

			When("kmsKeyARN is not empty but is not prefixed with 'mrk'", func() {
				It("should return an error", func() {
					invalidKmsKeyARN := "arn:aws:notkms:us-west-2:301721915996:key/9fdfaf2f-efb7-4db7-a5c3-0d047c52f094"
					err := ValidateKMSKeyARN(&invalidKmsKeyARN)
					Expect(err).To(HaveOccurred())
				})
			})

			When("when kmsKeyARN is not empty and does not match the regex", func() {
				It("should return an error", func() {
					invalidKmsKeyARN := "invalid-kms-key-arn"
					err := ValidateKMSKeyARN(&invalidKmsKeyARN)
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
