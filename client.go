package gcm

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	gcmEndpoint = "https://gcm-http.googleapis.com/gcm/send"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
	endpoint   string
}

func NewClient(apiKey string) *Client {
	return NewClientWithOptions(apiKey, nil, nil)
}

func NewClientWithOptions(apiKey string, httpClient *http.Client, endpointURL *url.URL) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	endpoint := gcmEndpoint
	if endpointURL != nil {
		endpoint = endpointURL.String()
	}

	return &Client{apiKey: apiKey, httpClient: httpClient, endpoint: endpoint}
}

func (c *Client) Send(message *Message) (*response, *gcmError) {
	req, err := c.createHTTPRequest(message)

	if err != nil {
		return nil, &gcmError{Message: err.Error()}
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &gcmError{Message: err.Error(), Type: ConnectionError, ShouldRetry: true}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, &gcmError{Message: err.Error()}
	}
	res.Body.Close()

	//log.Printf("RESPONSE: %#v\n", res)
	//log.Printf("BODY: %#v\n", string(body))

	switch {
	case res.StatusCode == 400:
		return nil, &gcmError{Type: BadRequestError, Message: string(body)}
	case res.StatusCode == 401:
		return nil, &gcmError{Type: AuthenticationError}
	case res.StatusCode == 413:
		return nil, &gcmError{Type: RequestEntityTooLargeError}
	case res.StatusCode >= 500:
		return nil, &gcmError{Type: InternalServerError, ShouldRetry: true}
	case res.StatusCode != 200:
		return nil, &gcmError{ShouldRetry: false, Message: string(body)}
	}

	responseObj := &response{}
	if err := json.Unmarshal(body, responseObj); err != nil {
		return nil, &gcmError{Type: ResponseParseError, Message: err.Error()}
	}

	return responseObj, nil
}

func (c *Client) createHTTPRequest(message *Message) (*http.Request, error) {
	body, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "key="+c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	defer req.Body.Close()

	return req, nil
}
