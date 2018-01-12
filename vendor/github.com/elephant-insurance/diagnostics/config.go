package diagnostics

import "time"

// DiagnosticConfig provides all needed variables for Diagnostics,
// primarily for testing Diagnostics itself.
type diagnosticConfig struct {
	ProcName    string
	Version     string
	Hostname    string
	StartTime   *time.Time
	AppModTime  time.Time
	Diagnostics map[string]diagnosticTest
	TestArray   []string
	Logger      DiagnosticLogger
}
