package gcm

import "fmt"

// Error
type gcmError struct {
	Type        gcmErrorType
	Message     string
	ShouldRetry bool
}

type gcmErrorType int

const (
	BadRequestError            gcmErrorType = 400
	AuthenticationError        gcmErrorType = 401
	RequestEntityTooLargeError gcmErrorType = 413
	InternalServerError        gcmErrorType = 500
	ResponseParseError         gcmErrorType = 1000
	ConnectionError            gcmErrorType = 1001
)

func (t gcmErrorType) String() string {
	switch t {
	case BadRequestError:
		return "Bad request"
	case AuthenticationError:
		return "Authentication error"
	case RequestEntityTooLargeError:
		return "Request entity too large"
	case InternalServerError:
		return "GCM server error"
	case ResponseParseError:
		return "Response body is not well-formed json"
	case ConnectionError:
		return "Connection error"
	default:
		return "Unknown error"
	}
}

func (e *gcmError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Type.String(), e.Message)
}

type resultError string

// Message level error codes
const (
	MissingRegistration       resultError = "MissingRegistration"
	InvalidRegistration       resultError = "InvalidRegistration"
	NotRegistered             resultError = "NotRegistered"
	MessageTooBig             resultError = "MessageTooBig"
	InvalidDataKey            resultError = "InvalidDataKey"
	InvalidTTL                resultError = "InvalidTtl"
	DeviceMessageRateExceeded resultError = "DeviceMessageRateExceeded"
	TopicsMessageRateExceeded resultError = "TopicsMessageRateExceeded"
	MismatchSenderID          resultError = "MismatchSenderId"
)
