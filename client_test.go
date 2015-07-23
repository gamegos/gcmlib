package gcmlib

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	apiKey := "bar"
	customHTTPClient := &http.Client{}
	customURL := "https://example.org"

	var clientTests = []struct {
		opts Config
		out  *Client
	}{
		{
			Config{APIKey: apiKey},
			&Client{&Config{APIKey: apiKey, HTTPClient: http.DefaultClient, Endpoint: gcmEndpoint}},
		},
		{
			Config{HTTPClient: customHTTPClient},
			&Client{&Config{APIKey: "", HTTPClient: customHTTPClient, Endpoint: gcmEndpoint}},
		},
		{
			Config{Endpoint: customURL},
			&Client{&Config{APIKey: "", HTTPClient: http.DefaultClient, Endpoint: customURL}},
		},
		{
			Config{APIKey: apiKey, Endpoint: customURL, HTTPClient: customHTTPClient},
			&Client{&Config{APIKey: apiKey, HTTPClient: customHTTPClient, Endpoint: customURL}},
		},
	}

	for _, tt := range clientTests {
		have, want := NewClient(tt.opts), tt.out
		if !reflect.DeepEqual(have, want) {
			t.Errorf("NewClient(%q): have: %#v, want: %#v\n", tt.opts, have, want)
		}
	}
}

func TestCreateHTTPRequest(t *testing.T) {
	msg, endpoint, apiKey := &Message{}, gcmEndpoint, "foo"

	req, err := createHTTPRequest(msg, endpoint, apiKey)

	if err != nil {
		t.Fatalf("createHTTPRequest: must succeed")
	}
	if req.URL.String() != endpoint {
		t.Errorf("createHTTPRequest: endpoint: have: %s, want: %s", req.URL.String(), endpoint)
	}

	if req.Method != "POST" {
		t.Errorf("createHTTPRequest: method: have: %s, want: %s", req.Method, "POST")
	}

	authHeader := "key=" + apiKey
	if req.Header.Get("Authorization") != authHeader {
		t.Errorf("createHTTPRequest: Authorization: have: %s, want: %s", req.Header.Get("Authorization"), apiKey)
	}
}

type mockResponse struct {
	statusCode int
	body       string
	headers    map[string]string
}

func mockHTTPServerAndClient(responses ...*mockResponse) (*httptest.Server, *http.Client) {
	i := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := responses[i]
		i++
		w.WriteHeader(res.statusCode)
		for k, v := range res.headers {
			w.Header().Set(k, v)
		}

		fmt.Fprintf(w, "%s", res.body)
	})

	server := httptest.NewServer(handler)

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	httpClient := &http.Client{Transport: transport}

	return server, httpClient
}

var sendTests = []struct {
	response     *mockResponse
	wantResponse *response
	wantError    *gcmError
}{
	// Errors
	{&mockResponse{400, "error msg", nil}, nil, newError(ErrorBadRequest, "error msg")},
	{&mockResponse{401, "error authentication", nil}, nil, newError(ErrorAuthentication, "")},
	{&mockResponse{413, "", nil}, nil, newError(ErrorRequestEntityTooLarge, "")},
	{&mockResponse{500, "", nil}, nil, newError(ErrorServiceUnavailable, "")},
	{&mockResponse{10, "error msg", nil}, nil, newError(ErrorUnknown, "error msg")},
	{&mockResponse{200, "{", nil}, nil, newError(ErrorResponseParse, "unexpected end of JSON input")}, // invalid json

	// Success
	{
		&mockResponse{200, "{\"multicast_id\":1,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"x\"}]}", nil},
		&response{
			MulticastID: 1,
			Success:     1,
			Results:     []result{result{MessageID: "x"}},
		},
		nil,
	},
}

func TestSend(t *testing.T) {
	for _, tt := range sendTests {
		server, httpClient := mockHTTPServerAndClient(tt.response)

		client := NewClient(Config{HTTPClient: httpClient, Endpoint: server.URL})

		haveRes, haveErr := client.Send(&Message{})
		server.Close()

		if !reflect.DeepEqual(haveRes, tt.wantResponse) {
			t.Errorf("Send/response: have: %#v, want: %#v\n", haveRes, tt.wantResponse)
		}

		if !reflect.DeepEqual(haveErr, tt.wantError) {
			t.Errorf("Send/error: have: %#v, want: %#v\n", haveErr, tt.wantError)
		}
	}
}

var resp500 = &mockResponse{500, "Internal Server Error", nil}
var resp200 = &mockResponse{200, "{\"multicast_id\":1,\"success\":1,\"failure\":0,\"canonical_ids\":0,\"results\":[{\"message_id\":\"x\"}]}", nil}

var sendRetryTests = []struct {
	response     []*mockResponse
	wantResponse *response
	wantError    *gcmError
	clientConfig Config
}{
	// retry & fail
	{
		[]*mockResponse{resp500, resp500, resp500, resp500},
		nil,
		newError(ErrorServiceUnavailable, ""),
		Config{MaxRetries: 3},
	},

	// retry & success
	{
		[]*mockResponse{resp500, resp500, resp500, resp500, resp200},
		&response{
			MulticastID: 1,
			Success:     1,
			Results:     []result{result{MessageID: "x"}},
		},
		nil,
		Config{MaxRetries: 4},
	},
}

func TestRetry(t *testing.T) {
	for _, tt := range sendRetryTests {
		server, httpClient := mockHTTPServerAndClient(tt.response...)

		opts := tt.clientConfig

		opts.Endpoint = server.URL
		opts.HTTPClient = httpClient

		client := NewClient(opts)

		haveRes, haveErr := client.Send(&Message{})
		server.Close()

		if !reflect.DeepEqual(haveRes, tt.wantResponse) {
			t.Errorf("Send/response: have: %#v, want: %#v\n", haveRes, tt.wantResponse)
		}

		if !reflect.DeepEqual(haveErr, tt.wantError) {
			t.Errorf("Send/error: have: %#v, want: %#v\n", haveErr, tt.wantError)
		}
	}

}
