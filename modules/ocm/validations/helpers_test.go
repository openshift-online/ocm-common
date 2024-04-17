package validations

import (
	"fmt"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ValidateIssuerUrlMatchesAssumePolicyDocument", func() {
	It("should return an error when the policy document does not have a trusted relationship to the issuer URL", func() {
		roleArn := "your-role-arn"
		parsedUrl, _ := url.Parse("http://your-issuer-url.com")
		assumePolicyDocument := "your-assume-policy-document"

		err := ValidateIssuerUrlMatchesAssumePolicyDocument(roleArn, parsedUrl, assumePolicyDocument)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(fmt.Sprintf("Operator role '%s' does not have trusted relationship to '%s' issuer URL",
			roleArn, parsedUrl.Host+parsedUrl.Path)))
	})

	It("should not return an error when the policy document has a trusted relationship to the issuer URL", func() {
		roleArn := "your-role-arn"
		parsedUrl, _ := url.Parse("http://your-issuer-url.com")
		assumePolicyDocument := "your-assume-policy-document with your-issuer-url.com"

		err := ValidateIssuerUrlMatchesAssumePolicyDocument(roleArn, parsedUrl, assumePolicyDocument)

		Expect(err).ToNot(HaveOccurred())
	})
})
