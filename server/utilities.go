package server

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

// getServerURL extracts the server URL from the Echo context request
// Returns format: {schema}://{host}[:{port}] where port is omitted for standard ports (80, 443)
func getServerURL(c echo.Context) string {
	req := c.Request()
	scheme := c.Scheme()
	host := req.Host

	// Fallback to localhost if host is empty (shouldn't happen in normal scenarios)
	if host == "" {
		// Try to get from X-Forwarded-Host header
		host = req.Header.Get("X-Forwarded-Host")
		if host == "" {
			host = "localhost"
		}
	}

	// Remove port from host if it's the default port for the scheme
	if strings.Contains(host, ":") {
		hostParts := strings.Split(host, ":")
		port := hostParts[1]
		if (scheme == "http" && port == "80") || (scheme == "https" && port == "443") {
			host = hostParts[0]
		}
	}

	return fmt.Sprintf("%s://%s", scheme, host)
}
