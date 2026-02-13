// Package version provides build-time version information.
package version

import "fmt"

// Build-time variables set via -ldflags.
var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

// Info returns a formatted version string.
func Info() string {
	return fmt.Sprintf("go-chat %s (commit: %s, built: %s)", version, commit, buildDate)
}

// Version returns the version string.
func Version() string {
	return version
}
