package errors

import "net/http"

const (
	APOLLO001 = "APOLLO_001"
	APOLLO002 = "APOLLO_002"
	APOLLO005 = "APOLLO_005"
	APOLLO006 = "APOLLO_006"
	APOLLO007 = "APOLLO_007"
	APOLLO008 = "APOLLO_008"
)

var (
	DataAlreadyExistsErr      = New(http.StatusOK, APOLLO001, "Data already exists").WithReason
	MissingRequestBodyErr     = New(http.StatusBadRequest, APOLLO002, "Request mismatch ")
	InternalServerErr         = New(http.StatusInternalServerError, APOLLO005, "We could not process your request due to malformed request, please check again").WithReason
	AuthorizationErr          = New(http.StatusUnauthorized, APOLLO005, "Authorization required").WithReason
	FailedParseRequestBodyErr = New(http.StatusUnauthorized, APOLLO006, "Failed parse request body").WithReason
	BadRequestWithReasonErr   = New(http.StatusBadRequest, APOLLO007, "Bad request").WithReason
	DecodeQueryParamErr       = New(http.StatusBadRequest, APOLLO008, "Bad request").WithReason("failed to decode query parameters")
)

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
