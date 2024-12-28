// Package testutil provides testing utilities and helpers
package testutil

// Window size constants
const (
	DefaultTestWindowWidth  = 80
	DefaultTestWindowHeight = 24
	LargeTestWindowWidth    = 120
	LargeTestWindowHeight   = 40
	SmallTestWindowWidth    = 40
	SmallTestWindowHeight   = 10
)

// Test timeouts and intervals
const (
	DefaultTestTimeout  = 5  // seconds
	DefaultTestInterval = 50 // milliseconds
)

// Log configuration constants
const (
	DefaultTestLogLevel      = "debug"
	DefaultTestLogMaxSize    = 1
	DefaultTestLogMaxBackups = 0
	DefaultTestLogMaxAge     = 1
)
