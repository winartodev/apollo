package errors

type Errors interface {
	Error() *Data
	WithReason(reason string) *Data
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
