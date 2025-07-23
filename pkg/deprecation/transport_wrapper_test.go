package deprecation_test

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/openshift-online/ocm-common/pkg/deprecation"
	"github.com/openshift-online/ocm-common/pkg/ocm/consts"
)

// mockRoundTripper implements http.RoundTripper for testing
type mockRoundTripper struct {
	response *http.Response
	err      error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

// captureStderr captures stderr output for testing
func captureStderr(fn func()) string {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	fn()

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	Expect(err).NotTo(HaveOccurred())
	return buf.String()
}

var _ = Describe("TransportWrapper", func() {
	var (
		wrapper       http.RoundTripper
		mockTransport *mockRoundTripper
		req           *http.Request
	)

	BeforeEach(func() {
		mockTransport = &mockRoundTripper{}
		wrapper = deprecation.NewTransportWrapper()(mockTransport)
		req, _ = http.NewRequest("GET", "https://api.example.com/test", nil)
	})

	Context("when response has deprecation headers", func() {
		It("should print warning with both deprecation and message headers", func() {
			headers := http.Header{}
			headers.Set(consts.DeprecationHeader, "2050-12-31T23:59:59Z")
			headers.Set(consts.OcmDeprecationMessage, "This endpoint is deprecated")

			mockTransport.response = &http.Response{
				StatusCode: http.StatusOK,
				Header:     headers,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}

			output := captureStderr(func() {
				resp, err := wrapper.RoundTrip(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			Expect(output).To(ContainSubstring("WARNING: You are using a deprecated OCM API"))
			Expect(output).To(ContainSubstring("This endpoint will be removed on: 2050-12-31T23:59:59Z"))
			Expect(output).To(ContainSubstring("Details: This endpoint is deprecated"))
		})

		It("should print warning with only deprecation header", func() {
			headers := http.Header{}
			headers.Set(consts.DeprecationHeader, "2050-12-31T23:59:59Z")

			mockTransport.response = &http.Response{
				StatusCode: http.StatusOK,
				Header:     headers,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}

			output := captureStderr(func() {
				resp, err := wrapper.RoundTrip(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			Expect(output).To(ContainSubstring("WARNING: You are using a deprecated OCM API"))
			Expect(output).To(ContainSubstring("This endpoint will be removed on: 2050-12-31T23:59:59Z"))
			Expect(output).NotTo(ContainSubstring("Details:"))
		})

		It("should print warning with only message header", func() {
			headers := http.Header{}
			headers.Set(consts.OcmDeprecationMessage, "Use v2 API instead")

			mockTransport.response = &http.Response{
				StatusCode: http.StatusOK,
				Header:     headers,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}

			output := captureStderr(func() {
				resp, err := wrapper.RoundTrip(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			Expect(output).To(ContainSubstring("WARNING: You are using a deprecated OCM API"))
			Expect(output).To(ContainSubstring("Details: Use v2 API instead"))
			Expect(output).NotTo(ContainSubstring("Deprecation:"))
		})

		It("should parse RFC1123 date format correctly", func() {
			headers := http.Header{}
			headers.Set(consts.DeprecationHeader, "Mon, 31 Dec 2050 23:59:59 GMT")

			mockTransport.response = &http.Response{
				StatusCode: http.StatusOK,
				Header:     headers,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}

			output := captureStderr(func() {
				resp, err := wrapper.RoundTrip(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			Expect(output).To(ContainSubstring("Deprecation: Mon, 31 Dec 2050 23:59:59 GMT"))
		})
		It("should parse RFC3339 date format correctly", func() {
			headers := http.Header{}
			headers.Set(consts.DeprecationHeader, "2040-12-31T23:59:59Z")

			mockTransport.response = &http.Response{
				StatusCode: http.StatusOK,
				Header:     headers,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}

			output := captureStderr(func() {
				resp, err := wrapper.RoundTrip(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			Expect(output).To(ContainSubstring("This endpoint will be removed on: 2040-12-31T23:59:59Z"))
		})
	})

	Context("when response has no deprecation headers", func() {
		It("should not print any warning", func() {
			headers := http.Header{}

			mockTransport.response = &http.Response{
				StatusCode: http.StatusOK,
				Header:     headers,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}

			output := captureStderr(func() {
				resp, err := wrapper.RoundTrip(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			Expect(output).To(BeEmpty())
		})
	})

	Context("when underlying transport returns error", func() {
		It("should return the error without processing", func() {
			mockTransport.err = http.ErrUseLastResponse

			output := captureStderr(func() {
				resp, err := wrapper.RoundTrip(req)
				Expect(err).To(Equal(http.ErrUseLastResponse))
				Expect(resp).To(BeNil())
			})

			Expect(output).To(BeEmpty())
		})
	})

	Context("when response is nil", func() {
		It("should not panic and not print warnings", func() {
			mockTransport.response = nil

			output := captureStderr(func() {
				resp, err := wrapper.RoundTrip(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			Expect(output).To(BeEmpty())
		})
	})

	Context("when response headers are nil", func() {
		It("should not panic and not print warnings", func() {
			mockTransport.response = &http.Response{
				StatusCode: http.StatusOK,
				Header:     nil,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}

			output := captureStderr(func() {
				resp, err := wrapper.RoundTrip(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			Expect(output).To(BeEmpty())
		})
	})

	Context("when deprecation header has invalid date format", func() {
		It("should still show warning without parsed date", func() {
			headers := http.Header{}
			headers.Set(consts.DeprecationHeader, "invalid-date")

			mockTransport.response = &http.Response{
				StatusCode: http.StatusOK,
				Header:     headers,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}

			output := captureStderr(func() {
				resp, err := wrapper.RoundTrip(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			Expect(output).To(ContainSubstring("WARNING: You are using a deprecated OCM API"))
			Expect(output).To(ContainSubstring("Deprecation: invalid-date"))
			Expect(output).NotTo(ContainSubstring("This endpoint will be removed on:"))
		})
	})

	Context("when testing different HTTP status codes", func() {
		It("should handle deprecation headers regardless of status code", func() {
			headers := http.Header{}
			headers.Set(consts.OcmDeprecationMessage, "Deprecated")

			// Test with different status codes
			statusCodes := []int{
				http.StatusOK,
				http.StatusCreated,
				http.StatusBadRequest,
				http.StatusInternalServerError,
			}

			for _, statusCode := range statusCodes {
				mockTransport.response = &http.Response{
					StatusCode: statusCode,
					Header:     headers,
					Body:       io.NopCloser(strings.NewReader("{}")),
				}

				output := captureStderr(func() {
					resp, err := wrapper.RoundTrip(req)
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.StatusCode).To(Equal(statusCode))
				})

				Expect(output).To(ContainSubstring("WARNING: You are using a deprecated OCM API"))
			}
		})
	})
})
