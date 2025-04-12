package errors

import "fmt"

type Errors interface {
	Error() *Data
	WithReason(reason string) *Data
	ToString() string
	ToError() error
}

type Data struct {
	StatusCode int
	ErrorCode  string `json:"error_code"`
	Message    string `json:"message"`
	Reason     string `json:"reason"`
}

func New(statusCode int, errorCode string, message string) Errors {
	return &Data{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

func (e *Data) WithReason(reason string) *Data {
	e.Reason = reason
	return e
}

func (e *Data) Error() *Data {
	return e
}

func (e *Data) ToString() string {
	if e.Reason != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Reason)
	}

	return e.Message
}

func (e *Data) ToError() error {
	return fmt.Errorf(e.ToString())
}
