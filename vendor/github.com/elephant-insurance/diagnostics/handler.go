package diagnostics

import (
	"encoding/json"
	"net/http"
)

type handler struct{}

// ServeHTTP is the implementation of the http.Handler interface. It serves. diagnostics page requests.
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r == nil {
		// not much we can do here
		logger.Printf("Diagnostics: Nil request sent to ServeHTTP")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	url := r.URL
	if url == nil {
		logger.Printf("Diagnostics: Nil URL sent to ServeHTTP")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	q := url.Query().Get(runTestsParameterName)

	code := http.StatusOK
	runTests := q != "" && q != "false"

	result := getCurrentSystemInfo(runTests)
	if runTests && result.FailedTests != nil && *result.FailedTests > 0 {
		code = http.StatusInternalServerError
	}

	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		if logger != nil {
			logger.Printf("Diagnostics: Error serializing result set:\n%v", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(resultJSON)
}
