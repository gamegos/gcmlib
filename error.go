package gcm

import "fmt"

// Error
type Error struct {
	Type        ErrorCode
	Message     string
	ShouldRetry bool
}

type ErrorCode int

const (
	BadRequestError            ErrorCode = 400
	AuthenticationError        ErrorCode = 401
	RequestEntityTooLargeError ErrorCode = 413
	InternalServerError        ErrorCode = 500

	ResponseParseError ErrorCode = 1000
)

func (t ErrorCode) String() string {
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
	default:
		return "Unknown error"
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.Type.String(), e.Message)
}

// ResultError
type ResultError string

// Message level error codes
const (
	MissingRegistration       ResultError = "MissingRegistration"
	InvalidRegistration       ResultError = "InvalidRegistration"
	NotRegistered             ResultError = "NotRegistered"
	MessageTooBig             ResultError = "MessageTooBig"
	InvalidDataKey            ResultError = "InvalidDataKey"
	InvalidTTL                ResultError = "InvalidTtl"
	DeviceMessageRateExceeded ResultError = "DeviceMessageRateExceeded"
	TopicsMessageRateExceeded ResultError = "TopicsMessageRateExceeded"
	MismatchSenderID          ResultError = "MismatchSenderId"
)
