package gcm

import "fmt"

// Error
type gcmError struct {
	code    gcmErrorCode
	message string
}

func (e *gcmError) Error() string {
	return fmt.Sprintf("%s: %s", e.code.String(), e.message)
}

func (e *gcmError) ShouldRetry() bool {
	return e.code == ErrorConnection || e.code == ErrorServiceUnavailable
}

func (e *gcmError) Code() gcmErrorCode {
	return e.code
}

func newError(errorCode gcmErrorCode, message string) *gcmError {
	return &gcmError{code: errorCode, message: message}
}

type gcmErrorCode int

// Errors
const (
	ErrorUnknown               gcmErrorCode = 1
	ErrorBadRequest            gcmErrorCode = 400
	ErrorAuthentication        gcmErrorCode = 401
	ErrorRequestEntityTooLarge gcmErrorCode = 413
	ErrorServiceUnavailable    gcmErrorCode = 500
	ErrorResponseParse         gcmErrorCode = 1000
	ErrorConnection            gcmErrorCode = 1001
)

func (t gcmErrorCode) String() string {
	switch t {
	case ErrorBadRequest:
		return "bad request"
	case ErrorAuthentication:
		return "authentication error"
	case ErrorRequestEntityTooLarge:
		return "request entity too large"
	case ErrorServiceUnavailable:
		return "gcm service unavailable"
	case ErrorResponseParse:
		return "response body is not well-formed json"
	case ErrorConnection:
		return "connection error"
	default:
		return "unknown error"
	}
}

type resultError string

// Result-level error codes can be found in response.Results for each individual
// recipient.
// See https://developers.google.com/cloud-messaging/server-ref#error-codes
const (
	ResultErrorMissingRegistration       resultError = "MissingRegistration"
	ResultErrorInvalidRegistration       resultError = "InvalidRegistration"
	ResultErrorNotRegistered             resultError = "NotRegistered"
	ResultErrorMessageTooBig             resultError = "MessageTooBig"
	ResultErrorInvalidDataKey            resultError = "InvalidDataKey"
	ResultErrorInvalidTTL                resultError = "InvalidTtl"
	ResultErrorDeviceMessageRateExceeded resultError = "DeviceMessageRateExceeded"
	ResultErrorTopicsMessageRateExceeded resultError = "TopicsMessageRateExceeded"
	ResultErrorMismatchSenderID          resultError = "MismatchSenderId"
	ResultErrorInvalidPackageName        resultError = "InvalidPackageName"
	ResultErrorInternalServerError       resultError = "InternalServerError"
	ResultErrorUnavailable               resultError = "Unavailable"
)
