package deprecation_test

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"go.uber.org/mock/gomock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/openshift-online/ocm-common/pkg/deprecation"
	"github.com/openshift-online/ocm-common/pkg/deprecation/test"
	"github.com/openshift-online/ocm-common/pkg/ocm/consts"
)

// captureStderr captures stderr output for test
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
		mockCtrl      *gomock.Controller
		mockTransport *test.MockRoundTripper
		req           *http.Request
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTransport = test.NewMockRoundTripper(mockCtrl)
		wrapper = deprecation.NewTransportWrapper()(mockTransport)
		req, _ = http.NewRequest("GET", "https://api.example.com/test", nil)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when response has deprecation headers", func() {
		It("should print warning with both deprecation and message headers", func() {
			headers := http.Header{}
			headers.Set(consts.DeprecationHeader, "2050-12-31T23:59:59Z")
			headers.Set(consts.OcmDeprecationMessage, "This endpoint is deprecated")

			mockTransport.EXPECT().RoundTrip(req).Return(&http.Response{
				StatusCode: http.StatusOK,
				Header:     headers,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}, nil)

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

			mockTransport.EXPECT().RoundTrip(req).Return(&http.Response{
				StatusCode: http.StatusOK,
				Header:     headers,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}, nil)

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

			mockTransport.EXPECT().RoundTrip(req).Return(&http.Response{
				StatusCode: http.StatusOK,
				Header:     headers,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}, nil)

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

			mockTransport.EXPECT().RoundTrip(req).Return(&http.Response{
				StatusCode: http.StatusOK,
				Header:     headers,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}, nil)

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

			mockTransport.EXPECT().RoundTrip(req).Return(&http.Response{
				StatusCode: http.StatusOK,
				Header:     headers,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}, nil)

			output := captureStderr(func() {
				resp, err := wrapper.RoundTrip(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			Expect(output).To(ContainSubstring("This endpoint will be removed on: 2040-12-31T23:59:59Z"))
		})
		It("should print warning with only field deprecation header", func() {
			headers := http.Header{}

			sunsetDate := time.Now().UTC().Add(time.Hour * 24 * 365)
			fieldDeprecations := deprecation.NewFieldDeprecations()
			err := fieldDeprecations.Add("field", "this field is deprecated", sunsetDate)
			Expect(err).NotTo(HaveOccurred())
			fieldDeprecationsJSON, _ := fieldDeprecations.ToJSON()
			headers.Set(consts.OcmFieldDeprecation, string(fieldDeprecationsJSON))

			mockTransport.EXPECT().RoundTrip(req).Return(&http.Response{
				StatusCode: http.StatusOK,
				Header:     headers,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}, nil)

			output := captureStderr(func() {
				resp, err := wrapper.RoundTrip(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			Expect(output).To(ContainSubstring("WARNING: You are using OCM API fields that have been deprecated"))
			Expect(output).To(ContainSubstring("this field is deprecated"))
			Expect(output).NotTo(ContainSubstring("Deprecation:"))
		})
		It("should error if sunset date is in the past", func() {
			sunsetDate := time.Now().UTC().Add(time.Hour * 24 * -365)
			fieldDeprecations := deprecation.NewFieldDeprecations()
			err := fieldDeprecations.Add("field", "this field is deprecated", sunsetDate)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("this field is deprecated"))
		})
	})

	Context("when response has no deprecation headers", func() {
		It("should not print any warning", func() {
			headers := http.Header{}

			mockTransport.EXPECT().RoundTrip(req).Return(&http.Response{
				StatusCode: http.StatusOK,
				Header:     headers,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}, nil)

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
			mockTransport.EXPECT().RoundTrip(req).Return(nil, http.ErrUseLastResponse)

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
			mockTransport.EXPECT().RoundTrip(req).Return(nil, nil)

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
			mockTransport.EXPECT().RoundTrip(req).Return(&http.Response{
				StatusCode: http.StatusOK,
				Header:     nil,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}, nil)

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

			mockTransport.EXPECT().RoundTrip(req).Return(&http.Response{
				StatusCode: http.StatusOK,
				Header:     headers,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}, nil)

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
				mockTransport.EXPECT().RoundTrip(req).Return(&http.Response{
					StatusCode: statusCode,
					Header:     headers,
					Body:       io.NopCloser(strings.NewReader("{}")),
				}, nil)

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
