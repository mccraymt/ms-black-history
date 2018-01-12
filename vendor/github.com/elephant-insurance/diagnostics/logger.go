package diagnostics

// DiagnosticLogger defines a simple interface for logging diagnostic events.
// Note that this logger will NOT log general diagnostic operation.
// It simply writes the same thing to the log that it writes to the web client
// every time handler.ServeHTTP is called. If an app appears to be having problems,
// a handle to a filewriter can be passed in to write a separate logfile just for diagnostics.
type DiagnosticLogger interface {
	Printf(format string, args ...interface{})
}

// EnableLogging enables logging to logger.
func EnableLogging(l DiagnosticLogger) {
	logger = l
}

// DisableLogging clears the logger so that logs are no longer written.
func DisableLogging() {
	logger = nil
}
