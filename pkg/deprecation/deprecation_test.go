package deprecation

import (
	"errors"
	"net/http"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-online/ocm-common/pkg/ocm/consts"
	sdk "github.com/openshift-online/ocm-sdk-go"
)

// MockResponse implements the ResponseWithHeaders interface for testing
type MockResponse struct {
	headers http.Header
}

func (m *MockResponse) Header() http.Header {
	return m.headers
}

// MockRequest implements RequestInterface for testing
type MockRequest struct {
	shouldError   bool
	errorToReturn error
	response      *sdk.Response
}

func (m *MockRequest) Send() (*sdk.Response, error) {
	if m.shouldError {
		return nil, m.errorToReturn
	}
	return m.response, nil
}

// MockTypedRequest implements TypedRequestInterface for testing
type MockTypedRequest struct {
	shouldError   bool
	errorToReturn error
	response      *MockResponse
}

func (m *MockTypedRequest) Send() (*MockResponse, error) {
	if m.shouldError {
		return nil, m.errorToReturn
	}
	return m.response, nil
}

var _ = Describe("Deprecation", func() {
	Describe("SendAndHandleDeprecation", func() {
		It("should handle success case without errors", func() {
			// Note: Testing SendAndHandleDeprecation is challenging because it requires
			// a real sdk.Response which is not easily mockable. This test focuses on
			// the error handling path, and the actual deprecation handling is tested
			// through the other HandleDeprecationWarning tests.

			// Create a mock request that succeeds with nil response (simulating success)
			mockRequest := &MockRequest{
				shouldError: false,
				response:    nil, // sdk.Response cannot be easily mocked
			}

			// Call the function
			_, err := SendAndHandleDeprecation(mockRequest)

			// Verify results (main focus is that it doesn't crash and propagates correctly)
			Expect(err).NotTo(HaveOccurred())
			// Response will be nil since we can't easily mock sdk.Response
			// The deprecation handling is tested in other test functions
		})

		It("should propagate errors from the request", func() {
			// Create a mock request that returns an error
			testError := errors.New("test error")
			mockRequest := &MockRequest{
				shouldError:   true,
				errorToReturn: testError,
			}

			// Call the function
			response, err := SendAndHandleDeprecation(mockRequest)

			// Verify error is propagated
			Expect(err).To(Equal(testError))
			Expect(response).To(BeNil())
		})
	})

	Describe("SendTypedAndHandleDeprecation", func() {
		It("should handle typed response with deprecation warning", func() {
			// Capture stderr for testing
			old := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Create a mock typed response with deprecation headers
			headers := http.Header{}
			headers.Set(consts.DeprecationHeader, "true")
			headers.Set(consts.OCMDeprecationMessage, "Typed test deprecation message")

			mockTypedResponse := &MockResponse{
				headers: headers,
			}

			// Create a mock typed request that returns the response
			mockRequest := &MockTypedRequest{
				shouldError: false,
				response:    mockTypedResponse,
			}

			// Call the function
			response, err := SendTypedAndHandleDeprecation(mockRequest)

			// Restore stderr
			w.Close()
			os.Stderr = old

			// Read the output
			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			// Verify results
			Expect(err).NotTo(HaveOccurred())
			Expect(response).NotTo(BeNil())
			Expect(output).To(Equal("Warning: Typed test deprecation message\n"))
		})

		It("should propagate errors from typed requests", func() {
			// Create a mock typed request that returns an error
			testError := errors.New("typed test error")
			mockRequest := &MockTypedRequest{
				shouldError:   true,
				errorToReturn: testError,
			}

			// Call the function
			response, err := SendTypedAndHandleDeprecation(mockRequest)

			// Verify error is propagated
			Expect(err).To(Equal(testError))
			Expect(response).To(BeNil())
		})

		It("should handle typed response without deprecation headers", func() {
			// Capture stderr for testing
			old := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Create a mock typed response without deprecation headers
			headers := http.Header{}

			mockTypedResponse := &MockResponse{
				headers: headers,
			}

			// Create a mock typed request that returns the response
			mockRequest := &MockTypedRequest{
				shouldError: false,
				response:    mockTypedResponse,
			}

			// Call the function
			response, err := SendTypedAndHandleDeprecation(mockRequest)

			// Restore stderr
			w.Close()
			os.Stderr = old

			// Read the output
			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			// Verify results
			Expect(err).NotTo(HaveOccurred())
			Expect(response).NotTo(BeNil())
			Expect(output).To(BeEmpty())
		})
	})

	Describe("HandleDeprecationWarningFromTypedResponse", func() {
		It("should handle OCM deprecation message header", func() {
			// Capture stderr for testing
			old := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Test with OCM deprecation message header (master message)
			headers := http.Header{}
			headers.Set(consts.DeprecationHeader, "1234567890")
			headers.Set(consts.OCMDeprecationMessage, "This endpoint is deprecated. Use /v2/ instead.")

			response := &MockResponse{
				headers: headers,
			}

			HandleDeprecationWarningFromTypedResponse(response)

			// Restore stderr
			w.Close()
			os.Stderr = old

			// Read the output
			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			expected := "Warning: This endpoint is deprecated. Use /v2/ instead.\n"
			Expect(output).To(Equal(expected))
		})

		It("should handle generic deprecation header", func() {
			// Capture stderr for testing
			old := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Test with only deprecation header (no custom message)
			headers := http.Header{}
			headers.Set(consts.DeprecationHeader, "true")

			response := &MockResponse{
				headers: headers,
			}

			HandleDeprecationWarningFromTypedResponse(response)

			// Restore stderr
			w.Close()
			os.Stderr = old

			// Read the output
			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			expected := "Warning: Deprecated endpoint was used\n"
			Expect(output).To(Equal(expected))
		})

		It("should handle future timestamp deprecation", func() {
			// Capture stderr for testing
			old := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Test with future timestamp (use UTC)
			futureTime := time.Now().UTC().Add(24 * time.Hour)
			headers := http.Header{}
			headers.Set(consts.DeprecationHeader, futureTime.Format(time.RFC3339))

			response := &MockResponse{
				headers: headers,
			}

			HandleDeprecationWarningFromTypedResponse(response)

			// Restore stderr
			w.Close()
			os.Stderr = old

			// Read the output
			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			Expect(output).To(ContainSubstring("Warning: This endpoint will be deprecated on"))
		})

		It("should handle no deprecation headers", func() {
			// Capture stderr for testing
			old := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Test with no deprecation header
			response := &MockResponse{
				headers: http.Header{},
			}

			HandleDeprecationWarningFromTypedResponse(response)

			// Restore stderr
			w.Close()
			os.Stderr = old

			// Read the output
			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			// Should be empty since no deprecation header was present
			Expect(output).To(BeEmpty())
		})
	})

	Describe("HandleDeprecationWarning", func() {
		It("should handle nil SDK response without panicking", func() {
			// Test that the function works with nil response
			// This should not panic and should not output anything
			HandleDeprecationWarning(nil)
		})
	})
})
