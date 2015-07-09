package gcmlib

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

type Options struct {
	APIKey     string
	HTTPClient *http.Client
	Endpoint   *url.URL
}

func NewClient(apiKey string) *Client {
	return NewClientWithOptions(&Options{APIKey: apiKey})
}

func NewClientWithOptions(options *Options) *Client {
	var httpClient *http.Client
	var endpoint string

	if options.HTTPClient != nil {
		httpClient = options.HTTPClient
	} else {
		httpClient = http.DefaultClient
	}

	if options.Endpoint != nil {
		endpoint = options.Endpoint.String()
	} else {
		endpoint = gcmEndpoint
	}

	return &Client{apiKey: options.APIKey, httpClient: httpClient, endpoint: endpoint}
}

func (c *Client) Send(message *Message) (*response, *gcmError) {
	req, err := createHTTPRequest(message, c.endpoint, c.apiKey)

	if err != nil {
		return nil, newError(ErrorUnknown, err.Error())
	}

	res, err := c.httpClient.Do(req)
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
