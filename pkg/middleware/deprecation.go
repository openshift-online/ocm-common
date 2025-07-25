package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/openshift-online/ocm-common/pkg/log"
	"github.com/openshift-online/ocm-common/pkg/ocm/consts"
)

// DeprecatedEndpoint represents a deprecated API endpoint with its message and sunset date.
type DeprecatedEndpoint struct {
	Message    string
	SunsetDate time.Time
}

// NewDeprecationMiddleware creates an HTTP middleware that adds deprecation headers
// and returns errors for expired endpoints. It accepts a map where keys are URL
// patterns and values are the deprecation details.
func NewDeprecationMiddleware(deprecatedEndpoints map[string]DeprecatedEndpoint) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the current request matches any deprecated endpoint
			deprecatedEndpoint, isDeprecated := matchDeprecatedEndpoint(r.URL.Path, deprecatedEndpoints)
			if isDeprecated {
				now := time.Now().UTC()
				// Check if the endpoint is expired (sunset date is in the past)
				if now.After(deprecatedEndpoint.SunsetDate) {
					// Return a standard 410 Gone error for expired endpoints.
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusGone)
					if err := json.NewEncoder(w).Encode(map[string]interface{}{
						"error":               "Gone",
						"deprecation_message": deprecatedEndpoint.Message,
						"sunset_date":         deprecatedEndpoint.SunsetDate.Format(time.RFC3339),
					}); err != nil {
						log.LogError("Failed to encode deprecation error response: %v", err)
					}
					return
				}
				// Add deprecation headers for active but deprecated endpoints
				w.Header().Set(consts.DeprecationHeader, deprecatedEndpoint.SunsetDate.Format(time.RFC3339))
				w.Header().Set(consts.OcmDeprecationMessage, deprecatedEndpoint.Message)
			}
			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// matchDeprecatedEndpoint checks if the given path matches any deprecated endpoint
func matchDeprecatedEndpoint(path string, deprecatedEndpoints map[string]DeprecatedEndpoint) (DeprecatedEndpoint, bool) {
	// Direct match first
	if endpoint, exists := deprecatedEndpoints[path]; exists {
		return endpoint, true
	}
	// Pattern matching for endpoints with path parameters
	for pattern, endpoint := range deprecatedEndpoints {
		if matchesPattern(path, pattern) {
			return endpoint, true
		}
	}
	return DeprecatedEndpoint{}, false
}

// matchesPattern checks if a path matches a pattern with path parameters
func matchesPattern(path, pattern string) bool {
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	if len(pathParts) != len(patternParts) {
		return false
	}
	for i, patternPart := range patternParts {
		// Skip path parameters (enclosed in curly braces)
		if strings.HasPrefix(patternPart, "{") && strings.HasSuffix(patternPart, "}") {
			continue
		}
		// Exact match required for non-parameter parts
		if pathParts[i] != patternPart {
			return false
		}
	}
	return true
}
