package diagnostics

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type diagnosticTest func() (DiagnosticResult, error)

// AddDiagnosticTest appends a new diagnostic test to be run with each invocation.
// The function will run if the runTests query string parameter is specified.
// It should return a useful result, which will be displayed on the diagnostics page.
func AddDiagnosticTest(name string, t diagnosticTest) {
	testArray = append(testArray, name)
	diagnosticTests[name] = t
}

func runDiagnostic(n string, t diagnosticTest, r *DiagnosticResult, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		if recov := recover(); recov != nil {
			f := false
			d := fmt.Sprintf("Recovered from PANIC: %v", recov)
			logger.Printf("Recovered from PANIC: %v", recov)
			pr := DiagnosticResult{
				Name:        n,
				Success:     &f,
				Description: &d,
			}
			*r = pr
		}
	}()
	startTime := time.Now()
	// we could just have the test func set result.Error,
	// but returning it separately is more idiomatic
	result, err := t()
	result.Name = n
	if err != nil {
		result.Error = err
	}
	elapsed := time.Since(startTime)
	result.Elapsed = elapsed.String()
	*r = result
}

// AddConnectionTest adds a simple connectivity test
// The url parameter should start with a protocol specifier (e.g., "http://").
// It will attempt to GET the specified URL. If the GET returns a 200, the result will have Success = true.
// Any code other than 200 will result in Success = false.
func AddConnectionTest(url string) {
	testName := fmt.Sprintf("Verify GET %v", url)
	test := func() (DiagnosticResult, error) {

		rtn := NewResult()
		res, err := http.Get(url)
		if err != nil {
			rtn.Fail().SetDescriptionf("Error attempting to GET %v: \n%v", url, err.Error())
			return *rtn, err
		}

		if res.StatusCode != http.StatusOK {
			rtn.Fail().SetDescriptionf("GET %v returned non-200 code: %v", url, res.StatusCode)
			return *rtn, nil
		}

		_, err = ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			rtn.Fail().SetDescriptionf("Error attempting to read response from GET %v: \n%v", url, err.Error())
			return *rtn, err
		}

		rtn.Succeed().SetDescriptionf("GET %v returned 200 OK", url)
		return *rtn, nil
	}

	AddDiagnosticTest(testName, test)
}
