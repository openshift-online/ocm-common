package utils

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Password Utils", func() {
	Context("when generating/hashing password", func() {
		It("should not return error", func() {
			password, err := GenerateRandomPassword()
			Expect(err).NotTo(HaveOccurred())
			_, err = GenerateHTPasswdCompatibleHash(password)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
