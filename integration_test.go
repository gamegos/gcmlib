// +build integration

package gcm

import (
	"flag"
	"testing"
)

var apiKey = flag.String("key", "", "GCM API KEY")
var regID = flag.String("regid", "", "A valid registration id")
var changedRegID = flag.String("cregid", "", "A changed registration id")
var dryRun = flag.Bool("dry", true, "Dry run")

type badRequestTestCase struct {
	message *Message
	err     *gcmError
}

var badRequestTestCases = []badRequestTestCase{
	{
		&Message{
			RegistrationIDs: make([]RegistrationID, 1001),
		},
		&gcmError{Type: BadRequestError, Message: "Number of messages on bulk (1001) exceeds maximum allowed (1000)\n", ShouldRetry: false},
	},
	{
		&Message{
			To:              "xx",
			RegistrationIDs: []RegistrationID{"id0"},
		},
		&gcmError{Type: BadRequestError, Message: "Must use either \"registration_ids\" field or \"to\", not both\n", ShouldRetry: false},
	},

	// reserved "from" keyword
	{
		&Message{
			To:   "xx",
			Data: map[string]string{"from": "reserved test"},
		},
		&gcmError{Type: BadRequestError, Message: "\"data\" key \"from\" is a reserved keyword\n", ShouldRetry: false},
	},

	// TTL
	{
		&Message{
			To:  "JohnDoe",
			TTL: maxTTL + 1,
		},
		&gcmError{Type: BadRequestError, Message: "Invalid value (2419201) for \"time_to_live\": must be between 0 and 2419200\n", ShouldRetry: false},
	},

	// Missing registration_ids
	{
		&Message{},
		&gcmError{Type: BadRequestError, Message: "Missing \"registration_ids\" field\n", ShouldRetry: false},
	},

	// message too long
	/*
		{
			&Message{
				To:           "xx",
				Notification: &Notification{Body: strings.Repeat("a", 1024*1024)},
			},
			&gcmError{Type: RequestEntityTooLargeError, ShouldRetry: false},
		},
	*/
}

func TestBadRequests(t *testing.T) {
	client := NewClient(*apiKey)

	for _, tc := range badRequestTestCases {
		res, err := client.Send(tc.message)

		if res != nil {
			t.Errorf("Response: expected: %#v, actual: %#v", nil, res)
		}

		if err.Type != tc.err.Type {
			t.Errorf("err.Type: expected: %#v, actual: %#v", tc.err.Type, err.Type)
		}

		if err.Message != tc.err.Message {
			t.Errorf("err.Message: expected: %#v, actual: %#v", tc.err.Message, err.Message)
		}

		if err.ShouldRetry != tc.err.ShouldRetry {
			t.Errorf("err.ShouldRetry: expected: %#v, actual: %#v", tc.err.ShouldRetry, err.ShouldRetry)
		}
	}
}

func TestAuthenticationError(t *testing.T) {
	client := NewClient("invalid-api-key")

	res, err := client.Send(&Message{})
	expectedErr := &gcmError{Type: AuthenticationError}

	if res != nil {
		t.Errorf("Response: expected: %#v, actual: %#v", nil, res)
	}

	if err.Type != expectedErr.Type {
		t.Errorf("err.Type: expected: %#v, actual: %#v", expectedErr.Type, err.Type)
	}

	if err.Message != expectedErr.Message {
		t.Errorf("err.Message: expected: %#v, actual: %#v", expectedErr.Message, err.Message)
	}

	if err.ShouldRetry != expectedErr.ShouldRetry {
		t.Errorf("err.ShouldRetry: expected: %#v, actual: %#v", expectedErr.ShouldRetry, err.ShouldRetry)
	}
}

func TestSuccess(t *testing.T) {
	if *regID == "" {
		t.Skip("skipping success test since no 'regid' parameter provided")
	}

	client := NewClient(*apiKey)
	msg := &Message{
		To:           RegistrationID(*regID),
		DryRun:       *dryRun,
		Notification: &Notification{Title: "gcm integration test message"},
	}
	t.Logf("Sending gcm message to: '%.40s...'", *regID)

	res, err := client.Send(msg)

	if err != nil {
		t.Errorf("Couldn't send gcm message: %s", err)
		return
	}

	if len(res.Results) != 1 {
		t.Errorf("No results returned")
		return
	}
	if !res.Results[0].Success() {
		t.Errorf("Particular message delivery problem: %s", res.Results[0].Error)
	}
}

func TestChangedRegistrationID(t *testing.T) {
	if *changedRegID == "" {
		t.Skip("skipping 'changed registration id' test since no 'cregid' parameter provided")
	}

	client := NewClient(*apiKey)
	msg := &Message{
		To:           RegistrationID(*changedRegID),
		DryRun:       *dryRun,
		Notification: &Notification{Title: "gcm integration test message"},
	}
	t.Logf("Sending gcm message to: %.80s...", *changedRegID)

	res, err := client.Send(msg)

	if err != nil {
		t.Errorf("Couldn't send gcm message: %s", err)
		return
	}

	if len(res.Results) != 1 {
		t.Errorf("No results returned")
		return
	}
	if !res.Results[0].Success() {
		t.Errorf("Particular message delivery problem: %s", res.Results[0].Error)
	}

	if !res.Results[0].TokenChanged() {
		t.Errorf("Provided registration id is already canonical")
	}

}
