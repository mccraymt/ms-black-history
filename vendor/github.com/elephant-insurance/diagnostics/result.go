package diagnostics

import "fmt"

// DiagnosticResult is a simple struct for passing back results from diagnostic tests.
type DiagnosticResult struct {
	Name        string  `json:"Name"`                  // The name of the test, which stays constant.
	Description *string `json:"Description,omitempty"` // The description of the result, including any explanation of the result.
	Success     *bool   `json:"Success,omitempty"`     // Whether the test reported success, failure, or nothing.
	Elapsed     string  `json:"Elapsed,omitempty"`     // How long it took to execute the test.
	Error       error   `json:"Error,omitempty"`       // Any error returned by the test.
}

// NewResult returns a pointer to a new, empty result struct.
func NewResult() *DiagnosticResult {
	rtn := DiagnosticResult{}

	return &rtn
}

// SetDescription adds a Description to the result.
func (r *DiagnosticResult) SetDescription(d string) *DiagnosticResult {
	if r == nil {
		return nil
	}

	// copy for safety
	dc := d
	r.Description = &dc

	return r
}

// SetDescriptionf sets a formatted description string.
func (r *DiagnosticResult) SetDescriptionf(format string, args ...interface{}) *DiagnosticResult {
	if r == nil {
		return nil
	}

	d := fmt.Sprintf(string(format), args...)

	r.SetDescription(d)

	return r
}

// Succeed sets the result.Success to true.
func (r *DiagnosticResult) Succeed() *DiagnosticResult {
	if r == nil {
		return nil
	}

	s := true
	r.Success = &s

	return r
}

// Fail sets the result.Success to false.
func (r *DiagnosticResult) Fail() *DiagnosticResult {
	if r == nil {
		return nil
	}

	s := false
	r.Success = &s

	return r
}
