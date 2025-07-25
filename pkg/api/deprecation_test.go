package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/openshift-online/ocm-common/pkg/ocm/consts"
)

// ErrorTemplate represents a template that can be used to send an error.
type ErrorTemplate struct {
	// Status is the HTTP status of the response.
	Status int

	// ID is the numeric identifier of the error within the service.
	ID int

	// Reason is a function that generates a default reason message for the error from the HTTP
	// request object. If this is nil then there will be no default error message.
	Reason func(r *http.Request) string
}

// Format generates an error reason from the given format and arguments, creates an instance of the
// error template and returns it. If the format string is empty, then the default reason of
// the template will be used.
func (t *ErrorTemplate) Format(r *http.Request, format string, a ...interface{}) Error {

	// Generate the reason message:
	var reason string
	if format != "" {
		reason = fmt.Sprintf(format, a...)
	} else if t.Reason != nil {
		reason = t.Reason(r)
	}

	// Prepare the body:
	return Error{
		ID:        fmt.Sprintf("%d", t.ID),
		Code:      fmt.Sprintf("CLUSTERS-MGMT-%d", t.ID),
		Reason:    reason,
		Timestamp: time.Now().UTC(),
	}
}

// E410 is a template for a generic HTTP `410 Gone` error.
var E410 = &ErrorTemplate{
	Status: http.StatusGone,
	ID:     410,
	Reason: func(r *http.Request) string {
		return fmt.Sprintf(
			"The requested resource '%s' is no longer available and will not be available again",
			r.URL.Path,
		)
	},
}

var _ = Describe("Deprecation Middleware", func() {
	var (
		nextHandler http.Handler
		handler     http.Handler
		rr          *httptest.ResponseRecorder
		nextCalled  bool
	)

	BeforeEach(func() {
		nextCalled = false
		nextHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			w.WriteHeader(http.StatusOK)
		})
		rr = httptest.NewRecorder()
	})

	Context("when endpoint is not deprecated", func() {
		It("should call the next handler without adding headers", func() {
			deprecatedEndpoints := map[string]DeprecatedEndpoint{}
			cfg := MiddlewareConfig{Endpoints: deprecatedEndpoints}
			handler = NewDeprecationMiddleware(cfg)(nextHandler)

			req := httptest.NewRequest("GET", "/api/test", nil)
			handler.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(nextCalled).To(BeTrue())
			Expect(rr.Header().Get(consts.DeprecationHeader)).To(BeEmpty())
			Expect(rr.Header().Get(consts.OcmDeprecationMessage)).To(BeEmpty())
		})
	})

	Context("when endpoint is deprecated but not expired", func() {
		It("should add deprecation headers and call the next handler", func() {
			sunsetDate := time.Now().Add(24 * time.Hour)
			deprecatedEndpoints := map[string]DeprecatedEndpoint{
				"/api/test": {
					Message:    "This is deprecated",
					SunsetDate: sunsetDate,
				},
			}
			cfg := MiddlewareConfig{Endpoints: deprecatedEndpoints}
			handler = NewDeprecationMiddleware(cfg)(nextHandler)

			req := httptest.NewRequest("GET", "/api/test", nil)
			handler.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(nextCalled).To(BeTrue())
			Expect(rr.Header().Get(consts.DeprecationHeader)).To(Equal(sunsetDate.Format(time.RFC3339)))
			Expect(rr.Header().Get(consts.OcmDeprecationMessage)).To(Equal("This is deprecated"))
		})
	})

	Context("when endpoint is expired", func() {
		It("should return 410 Gone and not call the next handler", func() {
			sunsetDate := time.Now().Add(-24 * time.Hour) // Expired
			deprecatedEndpoints := map[string]DeprecatedEndpoint{
				"/api/test": {
					Message:    "This is gone",
					SunsetDate: sunsetDate,
				},
			}

			var sentError *Error
			var createErrorCalled bool
			cfg := MiddlewareConfig{
				Endpoints: deprecatedEndpoints,
				CreateError: func(r *http.Request, format string, a ...interface{}) Error {
					createErrorCalled = true
					return E410.Format(r, "%v", deprecatedEndpoints["/api/test"].Message)
				},
				SendError: func(w http.ResponseWriter, r *http.Request, err *Error) {
					sentError = err
					status, conversionErr := strconv.Atoi(err.ID)
					Expect(conversionErr).ToNot(HaveOccurred())
					w.WriteHeader(status)
					w.Header().Set("Content-Type", "application/json")
				},
			}

			handler = NewDeprecationMiddleware(cfg)(nextHandler)

			req := httptest.NewRequest("GET", "/api/test", nil)
			handler.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusGone))
			Expect(nextCalled).To(BeFalse())
			Expect(createErrorCalled).To(BeTrue())
			Expect(sentError).ToNot(BeNil())
			Expect(sentError.Reason).To(ContainSubstring("This is gone"))
		})
	})

	Context("when endpoint with path parameter is deprecated", func() {
		It("should match the pattern and add deprecation headers", func() {
			sunsetDate := time.Now().Add(24 * time.Hour)
			deprecatedEndpoints := map[string]DeprecatedEndpoint{
				"/api/clusters/{id}": {
					Message:    "Use v2 instead",
					SunsetDate: sunsetDate,
				},
			}
			cfg := MiddlewareConfig{Endpoints: deprecatedEndpoints}
			handler = NewDeprecationMiddleware(cfg)(nextHandler)

			req := httptest.NewRequest("GET", "/api/clusters/12345", nil)
			handler.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(nextCalled).To(BeTrue())
			Expect(rr.Header().Get(consts.DeprecationHeader)).To(Equal(sunsetDate.Format(time.RFC3339)))
			Expect(rr.Header().Get(consts.OcmDeprecationMessage)).To(Equal("Use v2 instead"))
		})
	})
})

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Suite")
}

var _ = Describe("matchesPattern", func() {
	type testCase struct {
		path    string
		pattern string
		matches bool
	}

	DescribeTable("path matching",
		func(tc testCase) {
			Expect(matchesPattern(tc.path, tc.pattern)).To(Equal(tc.matches))
		},
		Entry("should match identical paths", testCase{
			path:    "/api/v1/test",
			pattern: "/api/v1/test",
			matches: true,
		}),
		Entry("should match with path parameter", testCase{
			path:    "/api/v1/clusters/123",
			pattern: "/api/v1/clusters/{id}",
			matches: true,
		}),
		Entry("should not match different paths", testCase{
			path:    "/api/v1/foo",
			pattern: "/api/v1/bar",
			matches: false,
		}),
		Entry("should not match if lengths are different", testCase{
			path:    "/api/v1/clusters/123/nodes",
			pattern: "/api/v1/clusters/{id}",
			matches: false,
		}),
		Entry("should handle multiple path parameters", testCase{
			path:    "/api/v1/clusters/123/nodes/456",
			pattern: "/api/v1/clusters/{cluster_id}/nodes/{node_id}",
			matches: true,
		}),
		Entry("should handle trailing slashes", testCase{
			path:    "/api/v1/test/",
			pattern: "/api/v1/test",
			matches: true,
		}),
	)
})
