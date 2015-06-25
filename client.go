package gcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	gcmEndpoint = "https://gcm-http.googleapis.com/gcm/send"
)

// Request level error codes
const (
	BadRequest          = "BadRequest"
	AuthenticationError = "AuthenticationError"
	InternalServerError = "InternalServerError"
	UnknownError        = "UnknownError"
)

// Message level error codes
const (
	MissingRegistrationError       = "MissingRegistration"
	InvalidRegistrationError       = "InvalidRegistration"
	NotRegisteredError             = "NotRegistered"
	MessageTooBigError             = "MessageTooBig"
	InvalidDataKeyError            = "InvalidDataKey"
	InvalidTTLError                = "InvalidTtl"
	DeviceMessageRateExceededError = "DeviceMessageRateExceeded"
	TopicsMessageRateExceededError = "TopicsMessageRateExceeded"
	MismatchSenderIDError          = "MismatchSenderId"
)

func NewClient(apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{apiKey: apiKey, httpClient: httpClient, Endpoint: gcmEndpoint}
}

type Client struct {
	apiKey     string
	httpClient *http.Client
	Endpoint   string
}

type GCMError struct {
	HTTPCode    int
	Type        string
	Message     string
	ShouldRetry bool
	//@TODO consider retry-after header
}

func (e *GCMError) Error() string {
	return fmt.Sprintf("[%s] %d: %s", e.Type, e.HTTPCode, e.Message)
}

func (c *Client) Send(message *Message) (*Response, error) {
	req, err := c.createHTTPRequest(message)

	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	//log.Printf("RESPONSE: %#v\n", res)
	//log.Printf("BODY: %#v\n", string(body))

	switch {
	case res.StatusCode == 400:
		return nil, &GCMError{HTTPCode: res.StatusCode, Type: InternalServerError, Message: string(body)}
	case res.StatusCode == 401:
		return nil, &GCMError{HTTPCode: res.StatusCode, Type: AuthenticationError}
	case res.StatusCode >= 500:
		return nil, &GCMError{HTTPCode: res.StatusCode, Type: InternalServerError, ShouldRetry: true}
	case res.StatusCode != 200:
		return nil, &GCMError{HTTPCode: res.StatusCode, Type: UnknownError, ShouldRetry: false, Message: string(body)}
	}

	responseObj := &Response{}
	if err := json.Unmarshal(body, responseObj); err != nil {
		return nil, err
	}

	return responseObj, err
}

func (c *Client) createHTTPRequest(message *Message) (*http.Request, error) {
	body, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "key="+c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	defer req.Body.Close()

	return req, err
}
