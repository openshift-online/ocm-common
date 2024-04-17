package validations

import (
	"fmt"
	"regexp"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Password Validator", func() {
	Context("Valid password", func() {
		It("should not return an error", func() {
			err := PasswordValidator("Abcdefg123456@")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Password contains special characters", func() {
		It("should return an error", func() {
			password := "AbcdefAbcdef@日本語"
			err := PasswordValidator(password)
			re := regexp.MustCompile(`[^\x20-\x7E]`)
			invalidChars := re.FindAllString(password, -1)
			expectedErrMsg := fmt.Sprintf("Password must not contain special characters [%s]", strings.Join(invalidChars, ", "))
			Expect(err.Error()).To(Equal(expectedErrMsg))
		})
	})

	Context("Password contains whitespace", func() {
		It("should return an error", func() {
			expectedErrMsg := "Password must not contain whitespace"
			err := PasswordValidator("Abc defg123456@")
			Expect(err.Error()).To(Equal(expectedErrMsg))
		})
	})

	Context("Password is too short", func() {
		It("should return an error", func() {
			password := "Abcd12@"
			expectedErrMsg := fmt.Sprintf("Password must be at least 14 characters (got %d)", len(password))
			err := PasswordValidator(password)
			Expect(err.Error()).To(Equal(expectedErrMsg))
		})
	})

	Context("Password does not contain uppercase letters", func() {
		It("should return an error", func() {
			expectedErrMsg := "Password must include uppercase letters, lowercase letters, and numbers or symbols (ASCII-standard characters only)"
			err := PasswordValidator("abcdefg123456@")
			Expect(err.Error()).To(Equal(expectedErrMsg))
		})
	})

	Context("Password does not contain lowercase letters", func() {
		It("should return an error", func() {
			expectedErrMsg := "Password must include uppercase letters, lowercase letters, and numbers or symbols (ASCII-standard characters only)"
			err := PasswordValidator("ABCDEFG123456@")
			Expect(err.Error()).To(Equal(expectedErrMsg))
		})
	})

	Context("Password does not contain numbers or symbols", func() {
		It("should return an error", func() {
			expectedErrMsg := "Password must include uppercase letters, lowercase letters, and numbers or symbols (ASCII-standard characters only)"
			err := PasswordValidator("Abcdefgabcdefg")
			Expect(err.Error()).To(Equal(expectedErrMsg))
		})
	})

	Context("White space and too short password errors", func() {
		It("should return an error", func() {
			password := "Abc def"
			err := PasswordValidator(password)

			expectedErrMsg := "Password must not contain whitespace, and must be at least 14 characters (got 7)"
			Expect(err.Error()).To(Equal(expectedErrMsg))
		})
	})

	Context("Multiple password errors", func() {
		It("should return an error", func() {
			password := "Abc 語def"
			err := PasswordValidator(password)
			re := regexp.MustCompile(`[^\x20-\x7E]`)
			invalidChars := re.FindAllString(password, -1)

			expectedErrMsg := fmt.Sprintf("Password must not contain special characters [%s], must not contain whitespace, and must be at least 14 characters (got %d)", strings.Join(invalidChars, ", "), len(password))
			Expect(err.Error()).To(Equal(expectedErrMsg))
		})
	})
})
