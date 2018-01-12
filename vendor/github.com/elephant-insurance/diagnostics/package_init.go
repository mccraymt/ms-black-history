package diagnostics

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"
)

const diagnosticsPackageVersion = "1.0.0"
const runTestsParameterName = "runTests"

var procName, version, hostname string
var startTime, appModTime time.Time
var diagnosticTests map[string]diagnosticTest
var testArray []string
var logger DiagnosticLogger
var initialized = false
var lastLogEntry string

func init() {
	startTime = time.Now()
	diagnosticTests = map[string]diagnosticTest{}
}

// Initialize sets up the diagnostic namespace.
// versionName should be set by the config file.
// logger, if not nil, will write diagnostics to the log each time the handler is invoked.
// The returned handler should be assigned to a route in the consuming application.
func Initialize(versionName string, appLogger DiagnosticLogger) (http.Handler, error) {
	processName := os.Args[0]
	diagnosticTests = map[string]diagnosticTest{}
	config := diagnosticConfig{}
	if processName == "" {
		return nil, errors.New("Diagnostics: processName cannot be empty")
	}
	if versionName == "" {
		return nil, errors.New("Diagnostics: versionName cannot be empty")
	}

	config.ProcName = getBaseFileName(processName)

	config.Version = versionName
	config.Logger = appLogger

	host, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("Diagnostics: Could not determine hostname: %v", err.Error())
	}
	config.Hostname = host

	workingDirectory, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Diagnostics: Could not determine working directory for running process %v: %v", procName, err.Error())
	}

	appFilepath := fmt.Sprintf("%v/%v", workingDirectory, config.ProcName)
	appFileHandle, err := os.Stat(appFilepath)
	if err != nil {
		return nil, fmt.Errorf("Diagnostics: Could not load file handle for running file %v: %v", appFilepath, err.Error())
	}

	if appFileHandle == nil || appFileHandle.Sys() == nil || appFileHandle.Mode().IsDir() {
		return nil, fmt.Errorf("Diagnostics: Invalid path for running file: %v", appFilepath)
	}

	config.AppModTime = appFileHandle.ModTime()

	return initializeFromConfig(config)
}

// InitializeFromConfig makes it possible to set up diagnostics without runtime calls, primarily for testing.
func initializeFromConfig(c diagnosticConfig) (http.Handler, error) {
	procName = c.ProcName
	version = c.Version
	hostname = c.Hostname
	if c.StartTime != nil {
		startTime = *c.StartTime
	}
	appModTime = c.AppModTime
	if c.Diagnostics != nil {
		diagnosticTests = c.Diagnostics
	}
	testArray = c.TestArray
	logger = c.Logger
	initialized = true

	return handler{}, nil
}

// Tokenize on path separators. Starting from the end and working toward the start,
//  return the first token with positive length
func getBaseFileName(p string) string {
	reg := regexp.MustCompile(`[/\\]`)
	tokens := reg.Split(p, -1)
	nt := len(tokens)
	for i := 1; i < nt+1; i++ {
		if i == nt || len(tokens[nt-i]) > 0 {
			return tokens[nt-i]
		}
	}
	return ""
}
