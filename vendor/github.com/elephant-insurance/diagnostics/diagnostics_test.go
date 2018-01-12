package diagnostics

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
	"time"
)

var rightNow = time.Now()
var halfHourAgo = rightNow.Add(time.Duration(-30) * time.Minute)
var hourAgo = rightNow.Add(time.Duration(-60) * time.Minute)
var logHasPrinted = false
var loggerInstance DiagnosticLogger

const testProcName = "testProcName"
const testVersion = "testVersion"
const testHostname = "testHostname"
const testFuncName = "testFuncName"
const testFuncDescription = "testFuncDescription"
const testFuncErrorText = "testFuncErrorText"
const testPanicMessage = "A Bad Thing has happened. Peace out."
const testGoodDiagnosticName = "testGoodDiagnosticName"
const testBadDiagnosticName = "testBadDiagnosticName"
const testPanicDiagnosticName = "testPanicDiagnosticName"

func TestInitializationFailure(t *testing.T) {
	fmt.Println("TestInitializationFailure")

	// attempt to initialize with a missing version string
	_, err := Initialize("", nil)
	if err == nil {
		fmt.Println("Failed to detect missing process name")
		t.Fail()
	}

	// attempt to initialize while running in test mode (should fail)
	_, err = Initialize("foo", nil)
	if err == nil {
		fmt.Println("Failed to detect test mode")
		t.Fail()
	}
}

func TestLoggerAndSystemInfo(t *testing.T) {
	fmt.Println("TestLoggerAndSystemInfo")
	_, err := initializeFromConfig(createConfig())
	if err != nil {
		fmt.Println("Failed to initialize without errors")
		t.Fail()
	}
	AddDiagnosticTest(testGoodDiagnosticName, GoodDiagnostic)
	AddDiagnosticTest(testBadDiagnosticName, BadDiagnostic)
	AddDiagnosticTest(testPanicDiagnosticName, PanicDiagnostic)

	logHasPrinted = false
	info := getCurrentSystemInfo(true)
	if !logHasPrinted {
		fmt.Println("Logger failed to print")
		t.Fail()
	}

	if info.AppName != testProcName {
		fmt.Println("Correct procName not returned")
		t.Fail()
	}

	if info.Version != testVersion {
		fmt.Println("Correct version not returned")
		t.Fail()
	}

	if info.HostName != testHostname {
		fmt.Println("Correct hostname not returned")
		t.Fail()
	}

	if info.Modified != hourAgo {
		fmt.Printf("Correct modified time not returned: %v\n", info.Modified)
		t.Fail()
	}

	if !strings.HasPrefix(info.Uptime, "30m") {
		fmt.Printf("Correct uptime not returned: %v\n", info.Uptime)
		t.Fail()
	}

	if len(info.DiagnosticResults) != 3 {
		fmt.Printf("Incorrect number of diagnostic results returned: %v\n", len(info.DiagnosticResults))
		t.Fail()
		t.Skip("Skipping tests of diagnostic results due to previous error")
	}

	if info.SuccessfulTests == nil || *info.SuccessfulTests != 1 {
		fmt.Printf("Incorrect number of success results returned (should be 1): %v\n", info.SuccessfulTests)
		t.Fail()
	}

	if info.FailedTests == nil || *info.FailedTests != 2 {
		fmt.Printf("Incorrect number of fail results returned (should be 2): %v\n", info.FailedTests)
		t.Fail()
	}

	goodResult := info.DiagnosticResults[0]
	if goodResult.Name != testGoodDiagnosticName {
		fmt.Printf("Correct test name not returned: %v\n", goodResult.Name)
		t.Fail()
	}

	if goodResult.Description == nil || *goodResult.Description != testFuncDescription {
		fmt.Printf("Correct test description not returned: %v\n", goodResult.Description)
		t.Fail()
	}

	if goodResult.Success == nil || *goodResult.Success != true {
		fmt.Printf("Correct test success not returned: %v\n", goodResult.Success)
		t.Fail()
	}

	badResult := info.DiagnosticResults[1]
	if badResult.Success == nil || *badResult.Success {
		fmt.Printf("Correct test success not returned: %v\n", badResult.Success)
		t.Fail()
	}

	if badResult.Error == nil || badResult.Error.Error() != testFuncErrorText {
		fmt.Printf("Correct test error not returned: %v\n", badResult.Error)
		t.Fail()
	}

	panicResult := info.DiagnosticResults[2]
	if panicResult.Success == nil || *panicResult.Success {
		fmt.Printf("Correct panic test success not returned: %v\n", panicResult.Success)
		t.Fail()
	}

	if panicResult.Description == nil || !strings.HasPrefix(*panicResult.Description, "Recovered from PANIC") {
		fmt.Printf("Correct panic test description not returned: %v\n", panicResult.Description)
		t.Fail()
	}
}

func TestHTTPHandler(t *testing.T) {
	// We can initialize as many times as we like to start fresh
	fmt.Println("TestLoggerAndSystemInfo")
	h, err := initializeFromConfig(createConfig())
	if err != nil {
		fmt.Println("Failed to initialize without errors")
		t.Fail()
	}

	// add a test to verify that we DON'T run it
	AddDiagnosticTest(testGoodDiagnosticName, GoodDiagnostic)

	// fake Request
	rq, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		fmt.Printf("Failed to create HTTP request: %v \n", err.Error())
		t.Fail()
	}

	// fake writer (recorder)
	wr := httptest.NewRecorder()

	h.ServeHTTP(wr, rq)

	if wr == nil {
		fmt.Printf("Handler returned nil response.\n")
		t.Fail()
	}

	if wr.Code != http.StatusOK {
		fmt.Printf("Handler returned non-200 response: %v \n", wr.Code)
		t.Fail()
	}

	si := SystemInfo{}
	err = json.Unmarshal(wr.Body.Bytes(), &si)
	if err != nil {
		fmt.Printf("Failed to unmarshal response: %v \n", err.Error())
		t.Fail()
	}

	// we've gotten back a properly-structured response, so let's make a quick sanity check
	if si.AppName != testProcName {
		fmt.Printf("Correct procName not returned, response body: %v \n", wr.Body.String())
		t.Fail()
	}

	if len(si.DiagnosticResults) > 0 {
		fmt.Printf("Test should not have run but %v results were reurned.\n", len(si.DiagnosticResults))
		t.Fail()
	}

	// check that a response with failures gives us a 500 response
	AddDiagnosticTest(testBadDiagnosticName, BadDiagnostic)
	wr = httptest.NewRecorder()
	rq, _ = http.NewRequest(http.MethodGet, "?runTests=true", nil)
	h.ServeHTTP(wr, rq)
	if wr.Code != http.StatusInternalServerError {
		fmt.Printf("Test should have returned a 500 response code: %v \n", wr.Code)
		t.Fail()
	}

	// check the exception handling of the handler
	logHasPrinted = false
	h.ServeHTTP(httptest.NewRecorder(), nil)
	if logHasPrinted == false || !strings.HasPrefix(lastLogEntry, "Diagnostics: Nil request sent") {
		fmt.Printf("Test should have returned a nil request error: %v \n", lastLogEntry)
		t.Fail()
	}

	logHasPrinted = false
	rq, err = http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		fmt.Printf("Failed to create HTTP request: %v \n", err.Error())
		t.Fail()
	}
	rq.URL = nil
	h.ServeHTTP(httptest.NewRecorder(), rq)
	if logHasPrinted == false || !strings.HasPrefix(lastLogEntry, "Diagnostics: Nil URL sent") {
		fmt.Printf("Test should have returned a nil request URL error: %v \n", lastLogEntry)
		t.Fail()
	}
}

func TestEnableDisableLogging(t *testing.T) {
	// initialize without a logger

	fmt.Println("TestEnableDisableLogging")
	c := createConfig()
	c.Logger = nil
	_, err := initializeFromConfig(c)
	if err != nil {
		fmt.Println("Failed to initialize without errors")
		t.Fail()
	}

	logHasPrinted = false
	_ = getCurrentSystemInfo(false)
	if logHasPrinted {
		fmt.Println("Logger printed when it shouldn't have - logger not specified.")
		t.Fail()
	}

	loggerInstance := testLogger{}
	EnableLogging(loggerInstance)
	logHasPrinted = false
	_ = getCurrentSystemInfo(false)
	if !logHasPrinted {
		fmt.Println("Logger didn't print when it should have.")
		t.Fail()
	}

	DisableLogging()
	logHasPrinted = false
	_ = getCurrentSystemInfo(false)
	if logHasPrinted {
		fmt.Println("Logger printed when it shouldn't have - DisableLogging() has been called.")
		t.Fail()
	}
}

func TestAbbreviate(t *testing.T) {
	// a trivial test that just increases coverage 8^)
	var ms *runtime.MemStats
	if abbreviate(ms) != nil {
		t.Fail()
	}
}

func TestGetRootFilename(t *testing.T) {
	fn := getBaseFileName("/foo/bar/quux")
	if fn != "quux" {
		t.Fail()
	}
	fn = getBaseFileName("/foo/bar/quux/")
	if fn != "quux" {
		t.Fail()
	}
	fn = getBaseFileName("quux")
	if fn != "quux" {
		t.Fail()
	}
	fn = getBaseFileName("./quux")
	if fn != "quux" {
		t.Fail()
	}
	fn = getBaseFileName(`c:\foo\bar\quux.exe`)
	if fn != "quux.exe" {
		t.Fail()
	}
	fn = getBaseFileName(`c:\\foo\\bar\\quux.exe`)
	if fn != "quux.exe" {
		t.Fail()
	}
	fn = getBaseFileName(`\\foo\bar\quux.exe`)
	if fn != "quux.exe" {
		t.Fail()
	}
	fn = getBaseFileName(`\\foo\bar\quux\`)
	if fn != "quux" {
		t.Fail()
	}
	fn = getBaseFileName("quux.exe")
	if fn != "quux.exe" {
		t.Fail()
	}
	fn = getBaseFileName(`.\quux.exe`)
	if fn != "quux.exe" {
		t.Fail()
	}
}

func TestNullResultHandling(t *testing.T) {
	var ptr *DiagnosticResult
	if ptr.SetDescription("") != nil || ptr.SetDescriptionf("%v", "foo") != nil || ptr.Succeed() != nil || ptr.Fail() != nil {
		t.Fail()
	}
}

func TestConnectionTester(t *testing.T) {
	fmt.Println("TestConnectionTester")
	_, err := initializeFromConfig(createConfig())
	if err != nil {
		fmt.Println("Failed to initialize without errors")
		t.Fail()
	}

	// should always return 200
	AddConnectionTest("http://www.yahoo.com")

	// should always return 418
	AddConnectionTest("http://httpstat.us/418")

	info := getCurrentSystemInfo(true)

	if len(info.DiagnosticResults) != 2 {
		fmt.Printf("Incorrect number of diagnostic results returned: %v\n", len(info.DiagnosticResults))
		t.Fail()
		t.Skip("Skipping tests of diagnostic results due to previous error")
	}

	if info.SuccessfulTests == nil || *info.SuccessfulTests != 1 {
		fmt.Printf("Incorrect number of success results returned (should be 1): %v\n", info.SuccessfulTests)
		t.Fail()
	}

	if info.FailedTests == nil || *info.FailedTests != 1 {
		fmt.Printf("Incorrect number of fail results returned (should be 1): %v\n", info.FailedTests)
		t.Fail()
	}

	goodResult := info.DiagnosticResults[0]
	badResult := info.DiagnosticResults[1]
	if goodResult.Success == nil || !*goodResult.Success {
		fmt.Printf("Failed to verify successful connection: %v\n", goodResult.Name)
		t.Fail()
	}
	if badResult.Success == nil || *badResult.Success {
		fmt.Printf("Failed to verify unsuccessful connection: %v\n", badResult.Name)
		t.Fail()
	}
}

func createConfig() diagnosticConfig {
	loggerInstance = testLogger{}
	return diagnosticConfig{
		ProcName:    testProcName,
		Version:     testVersion,
		Hostname:    testHostname,
		StartTime:   &halfHourAgo,
		AppModTime:  hourAgo,
		Diagnostics: map[string]diagnosticTest{},
		Logger:      loggerInstance,
	}
}

type testLogger struct {
}

func (t testLogger) Printf(format string, args ...interface{}) {
	// verify that the logger is firing
	logHasPrinted = true
	lastLogEntry = fmt.Sprintf(format, args)
}

func GoodDiagnostic() (DiagnosticResult, error) {
	r := NewResult().SetDescriptionf("%v", testFuncDescription).Succeed()

	return *r, nil
}

func BadDiagnostic() (DiagnosticResult, error) {

	r := NewResult().SetDescription(testFuncDescription).Fail()

	return *r, errors.New(testFuncErrorText)
}

func PanicDiagnostic() (DiagnosticResult, error) {
	panic(testPanicMessage)
}
