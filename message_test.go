package gcmlib

import (
	"reflect"
	"testing"
)

var messageValidateTests = []struct {
	msg *Message
	err error
}{
	{&Message{}, errNoRegID},
	{&Message{To: "foo", RegistrationIDs: []RegistrationID{"id0"}}, errBothToAndRegID},
	{&Message{RegistrationIDs: make([]RegistrationID, 1001)}, errExceedMaxRegIDs},
	{&Message{To: "foo", TTL: maxTTL + 1}, errInvalidTTL},
	{&Message{To: "foo", Priority: maxPriority + 1}, errInvalidPriority},
	{&Message{To: "foo", Data: map[string]string{"from": "bar"}}, errReservedDataKey},
	{&Message{To: "foo", Data: map[string]string{"google.key": "bar"}}, errReservedDataKeyPrefix},

	{&Message{To: "foo", Data: map[string]string{"foo": "bar"}}, nil},
}

func TestValidate(t *testing.T) {
	for _, tt := range messageValidateTests {
		have, want := tt.msg.Validate(), tt.err
		if !reflect.DeepEqual(have, want) {
			t.Errorf("Message.Validate(%q): have: %#v, want: %#v\n", tt.msg, have, want)
		}
	}
}
