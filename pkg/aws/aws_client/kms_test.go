package aws_client_test

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/openshift-online/ocm-common/pkg/aws/aws_client"
)

// stubKMSTransport is a local stand-in for the KMS client's HTTP transport. It
// never reaches AWS: every request is recorded and answered with a canned
// response, which keeps these specs hermetic and deterministic.
type stubKMSTransport struct {
	statusCode int
	body       string
	headers    http.Header

	requestBodies []string
}

func (t *stubKMSTransport) Do(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		t.requestBodies = append(t.requestBodies, string(body))
	}

	header := http.Header{"Content-Type": []string{"application/x-amz-json-1.1"}}
	for key, values := range t.headers {
		for _, value := range values {
			header.Add(key, value)
		}
	}

	return &http.Response{
		StatusCode: t.statusCode,
		Status:     http.StatusText(t.statusCode),
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(t.body)),
		Request:    req,
	}, nil
}

func newStubbedKMSAWSClient(transport *stubKMSTransport) *AWSClient {
	return &AWSClient{
		KmsClient: kms.New(kms.Options{
			Region:           "us-east-1",
			Credentials:      credentials.NewStaticCredentialsProvider("AKID", "SECRET", ""),
			HTTPClient:       transport,
			RetryMaxAttempts: 1,
		}),
	}
}

var _ = Describe("PutKMSPolicy", func() {
	It("returns a zero-value output and the original error when PutKeyPolicy fails", func() {
		transport := &stubKMSTransport{
			statusCode: http.StatusBadRequest,
			body:       `{"__type":"NotFoundException","message":"Key 'key-id' does not exist"}`,
			headers:    http.Header{"X-Amzn-Errortype": []string{"NotFoundException"}},
		}
		client := newStubbedKMSAWSClient(transport)

		// Before the fix this panicked: the failed call returns a nil
		// *kms.PutKeyPolicyOutput, which was dereferenced unconditionally.
		output, err := client.PutKMSPolicy("key-id", "", `{"Version":"2012-10-17"}`)

		Expect(err).To(HaveOccurred())
		var notFound *types.NotFoundException
		Expect(errors.As(err, &notFound)).To(BeTrue(), "the original AWS error must be propagated")
		Expect(output).To(Equal(kms.PutKeyPolicyOutput{}))
		Expect(transport.requestBodies).To(HaveLen(1))
		Expect(transport.requestBodies[0]).To(ContainSubstring(`"PolicyName":"default"`))
	})

	It("returns the output and no error when PutKeyPolicy succeeds", func() {
		transport := &stubKMSTransport{
			statusCode: http.StatusOK,
			body:       `{}`,
		}
		client := newStubbedKMSAWSClient(transport)

		output, err := client.PutKMSPolicy("key-id", "custom-policy", `{"Version":"2012-10-17"}`)

		Expect(err).NotTo(HaveOccurred())
		// The SDK populates ResultMetadata on a successful call, so a
		// successful output is never the zero value the error path returns.
		Expect(output).NotTo(Equal(kms.PutKeyPolicyOutput{}))
		Expect(transport.requestBodies).To(HaveLen(1))
		Expect(transport.requestBodies[0]).To(ContainSubstring(`"PolicyName":"custom-policy"`))
		Expect(transport.requestBodies[0]).To(ContainSubstring(`"KeyId":"key-id"`))
	})
})
