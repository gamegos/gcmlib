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
	BadRequestError     ErrorCode = 400
	AuthenticationError ErrorCode = 401
	InternalServerError ErrorCode = 500

	ResponseParseError ErrorCode = 1000
)

func (t ErrorCode) String() string {
	switch t {
	case BadRequestError:
		return "Bad request"
	case AuthenticationError:
		return "Authentication error"
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

// MessageError
type MessageError string

// Message level error codes
const (
	MissingRegistration       MessageError = "MissingRegistration"
	InvalidRegistration       MessageError = "InvalidRegistration"
	NotRegistered             MessageError = "NotRegistered"
	MessageTooBig             MessageError = "MessageTooBig"
	InvalidDataKey            MessageError = "InvalidDataKey"
	InvalidTTL                MessageError = "InvalidTtl"
	DeviceMessageRateExceeded MessageError = "DeviceMessageRateExceeded"
	TopicsMessageRateExceeded MessageError = "TopicsMessageRateExceeded"
	MismatchSenderID          MessageError = "MismatchSenderId"
)
