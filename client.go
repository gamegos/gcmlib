package gcmlib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	gcmEndpoint = "https://gcm-http.googleapis.com/gcm/send"
)

// The Client type encapsulates
type Client struct {
	config *Config
}

type Config struct {
	APIKey     string
	HTTPClient *http.Client
	Endpoint   string
	MaxRetries uint
}

func NewClient(config Config) *Client {
	if config.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	}

	if config.Endpoint == "" {
		config.Endpoint = gcmEndpoint
	}

	return &Client{config: &config}
}

func (c *Client) Send(message *Message) (*response, *gcmError) {
	r := uint(0)
	for {
		res, err := c.doSend(message)
		if err == nil {
			return res, nil
		}

		if !err.ShouldRetry() || c.config.MaxRetries < 1 {
			return nil, err
		}

		if r == c.config.MaxRetries {
			return nil, err
		}

		time.Sleep((1 << r) * 400 * time.Millisecond)
		r++
	}
}

func (c *Client) doSend(message *Message) (*response, *gcmError) {
	req, err := createHTTPRequest(message, c.config.Endpoint, c.config.APIKey)

	if err != nil {
		return nil, newError(ErrorUnknown, err.Error())
	}

	res, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, newError(ErrorConnection, err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, newError(ErrorUnknown, err.Error())
	}

	//log.Printf("RESPONSE: %#v\n", res)
	//log.Printf("BODY: %#v\n", string(body))

	switch {
	case res.StatusCode == 400:
		return nil, newError(ErrorBadRequest, string(body))
	case res.StatusCode == 401:
		return nil, newError(ErrorAuthentication, "")
	case res.StatusCode == 413:
		return nil, newError(ErrorRequestEntityTooLarge, "")
	case res.StatusCode >= 500:
		return nil, newError(ErrorServiceUnavailable, "")
	case res.StatusCode != 200:
		return nil, newError(ErrorUnknown, string(body))
	}

	responseObj := &response{}
	if err := json.Unmarshal(body, responseObj); err != nil {
		return nil, newError(ErrorResponseParse, err.Error())
	}

	return responseObj, nil
}

func createHTTPRequest(message *Message, endpoint string, apiKey string) (*http.Request, error) {
	body, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	req.Header.Set("Authorization", "key="+apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
