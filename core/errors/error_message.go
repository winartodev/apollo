package errors

import "net/http"

const (
	ReasonInvalidParamID = "param %s must greater than 0"
)

var (
	// Client Errors (4xx)
	MissingRequestBodyErr       = New(http.StatusBadRequest, APOLLO_002, "Request Body Is Required")
	FailedParseRequestBodyErr   = New(http.StatusUnauthorized, APOLLO_005, "Invalid Request Format")
	BadRequestErr               = New(http.StatusBadRequest, APOLLO_006, "Bad request")
	InvalidSlugErr              = New(http.StatusBadRequest, APOLLO_009, "Invalid Identifier Format")
	AuthorizationErr            = New(http.StatusUnauthorized, APOLLO_004, "Authentication required")
	InvalidUserID               = New(http.StatusUnauthorized, APOLLO_008, "Invalid User Credentials")
	UserApplicationHasNotAccess = New(http.StatusUnauthorized, USR_APP_001, "Access not granted for this operation")
	DataNotFoundErr             = New(http.StatusNotFound, APOLLO_007, "Requested Resource Not Found")
	InvalidContextParam         = New(http.StatusBadRequest, APOLLO_010, "Invalid Context Param")

	// Success with no action (2xx)
	DataAlreadyExistsErr   = New(http.StatusOK, APOLLO_001, "Resource Already Exists").WithReason
	ApplicationInactiveErr = New(http.StatusOK, APP_001, "Application temporarily unavailable")

	// Server Errors (5xx)
	InternalServerErr = New(http.StatusInternalServerError, APOLLO_003, "We Encountered An Unexpected Error").WithReason
)
