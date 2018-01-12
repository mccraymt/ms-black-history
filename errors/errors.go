package errors

import (
	"encoding/json"
	"time"
)

// Error is a struct that implements the Go error interface, adding some additional fields for consistency
type Error struct {
	SystemError *error         `json:"Error,omitempty"`
	Title       *string        `json:"Title,omitempty"`
	Description *string        `json:"Description,omitempty"`
	Code        *int           `json:"Code,omitempty"`
	Timestamp   *time.Time     `json:"Timestamp,omitempty"`
	Elapsed     *time.Duration `json:"Duration,omitempty"`
}

// New initializes a new Error object
func New(title, description *string, code *int, err *error, timestamp *time.Time, elapsed *time.Duration) Error {
	rtn := Error{
		SystemError: err,
		Title:       title,
		Description: description,
		Code:        code,
		Timestamp:   timestamp,
		Elapsed:     elapsed,
	}

	if rtn.Timestamp == nil {
		rightNow := time.Now()
		rtn.Timestamp = &rightNow
	}

	return rtn
}

func (err *Error) Error() string {
	rtn, _ := json.Marshal(err)
	return string(rtn)
}
