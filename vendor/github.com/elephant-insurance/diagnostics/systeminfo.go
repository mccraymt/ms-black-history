package diagnostics

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

// SystemInfo is the output type of the diagnostics page. With this struct, the results
// of a diagnostic page load can be loaded into memory and analyzed by another process.
type SystemInfo struct {
	HostName           string              `json:"HostName"`                    // The hostname (or container ID) for the machine we're running on.
	AppName            string              `json:"AppName"`                     // The filename of the running executible.
	Version            string              `json:"Version"`                     // The version of the app (set at initialization).
	DiagnosticsVersion string              `json:"DiagnosticsVersion"`          // The version of the diagnostics library itself.
	Modified           time.Time           `json:"Modified"`                    // The modified time of the running executible.
	Uptime             string              `json:"Uptime"`                      // How long it's been since the app started.
	SuccessfulTests    *int                `json:"SuccessfulTests,omitempty"`   // The number of DiagnosticTests that succeeded, if any.
	FailedTests        *int                `json:"FailedTests,omitempty"`       // The number of DiagnosticTests that failed, if any.
	MemStats           *memStat            `json:"MemStats"`                    // Memory usage statistics.
	DiagnosticResults  []*DiagnosticResult `json:"DiagnosticResults,omitempty"` // The full results of any DiagnosticTests that ran, if any.
	Notes              *string             `json:"Notes,omitempty"`             // Human-readable notes to make the diagnostics page more useful.
}

func getCurrentSystemInfo(runTests bool) SystemInfo {
	var stats runtime.MemStats
	numTests := len(testArray)
	runtime.ReadMemStats(&stats)
	shortStats := abbreviate(&stats)
	rtn := SystemInfo{
		HostName:           hostname,
		AppName:            procName,
		Version:            version,
		DiagnosticsVersion: diagnosticsPackageVersion,
		Modified:           appModTime,
		Uptime:             (time.Since(startTime)).String(),
		MemStats:           shortStats,
	}

	if runTests && numTests > 0 {
		results := make([]*DiagnosticResult, numTests)
		var wg sync.WaitGroup
		wg.Add(numTests)
		for i := 0; i < numTests; i++ {
			thisTestName := testArray[i]
			thisTest := diagnosticTests[thisTestName]
			thisResult := DiagnosticResult{}
			results[i] = &thisResult
			go runDiagnostic(thisTestName, thisTest, results[i], &wg)
		}

		wg.Wait()

		var ns, nf int
		for i := 0; i < numTests; i++ {
			if results[i] != nil {
				s := results[i].Success
				if s != nil {
					if *s {
						ns++
					} else {
						nf++
					}
				}
			}
		}

		if ns > 0 {
			rtn.SuccessfulTests = &ns
		}
		if nf > 0 {
			rtn.FailedTests = &nf
		}

		rtn.DiagnosticResults = results
	} else if numTests > 0 {
		note := fmt.Sprintf("This app has %v diagnostic tests that did not run: %v. \nTo run these tests, call this page again with the query string '?%v=true'.", numTests, strings.Join(testArray, ", "), runTestsParameterName)
		rtn.Notes = &note
	}

	if logger != nil {
		resultJSON, err := json.MarshalIndent(rtn, "", "  ")
		if err != nil {
			logger.Printf("Diagnostics: Error marshaling result: %v", err.Error())
		} else {
			logger.Printf("Diagnostics: Results of poll at %v:\n%v", time.Now(), string(resultJSON))
		}
	}

	return rtn
}
