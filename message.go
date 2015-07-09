package gcmlib

import (
	"errors"
	"fmt"
	"strings"
)

const (
	maxRegistrationIDs = 1000
	maxTTL             = 4 * 7 * 24 * 60 * 60
	maxPriority        = 10
)

type RegistrationID string

type Message struct {
	// Targets
	To              RegistrationID   `json:"to,omitempty"`
	RegistrationIDs []RegistrationID `json:"registration_ids"`

	// Options
	CollapseKey           string `json:"collapse_key,omitempty"`
	Priority              uint8  `json:"priority,omitempty"`
	ContentAvailable      bool   `json:"content_available,omitempty"`
	DelayWhileIdle        bool   `json:"delay_while_idle,omitempty"`
	TTL                   uint   `json:"time_to_live,omitempty"`
	RestrictedPackageName string `json:"restricted_package_name,omitempty"`
	DryRun                bool   `json:"dry_run,omitempty"`

	// Payload
	Notification *Notification     `json:"notification,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
}

type Notification struct {
	Title        string `json:"title,omitempty"`
	Body         string `json:"body,omitempty"`
	Icon         string `json:"icon,omitempty"`
	Sound        string `json:"sound,omitempty"`
	Badge        string `json:"badge,omitempty"`
	Tag          string `json:"tag,omitempty"`
	Color        string `json:"color,omitempty"`
	ClickAction  string `json:"click_action,omitempty"`
	BodyLocKey   string `json:"body_loc_key,omitempty"`
	BodyLocArgs  string `json:"body_loc_args,omitempty"`
	TitleLocArgs string `json:"title_loc_args,omitempty"`
	TitleLocKey  string `json:"title_loc_key,omitempty"`
}

var reservedDataKeys = []string{"from", "collapse_key"}

var (
	errNoRegID               = errors.New("requires at least one registration id")
	errBothToAndRegID        = errors.New("cannot set both 'to' and 'registration_ids' field")
	errInvalidPriority       = errors.New("priority value must be between 0 - 10")
	errExceedMaxRegIDs       = fmt.Errorf("registration_ids field cannot contain more than %d registration id", maxRegistrationIDs)
	errInvalidTTL            = fmt.Errorf("time_to_live value must be between 0 - %d", maxTTL)
	errReservedDataKey       = fmt.Errorf("data must not contain a reserved key. Reserved keys: %s", strings.Join(reservedDataKeys, ", "))
	errReservedDataKeyPrefix = errors.New("data must not contain keys beginning with gcm or google")
)

func (m *Message) Validate() error {
	if (m.To == "") && (len(m.RegistrationIDs) == 0) {
		return errNoRegID
	}

	if (m.To != "") && (len(m.RegistrationIDs) > 0) {
		return errBothToAndRegID
	}

	// check len(RegistrationIDs) > 1000
	if len(m.RegistrationIDs) > maxRegistrationIDs {
		return errExceedMaxRegIDs
	}

	// check TTL
	if m.TTL > maxTTL {
		return errInvalidTTL
	}

	// check priority
	if m.Priority > maxPriority {
		return errInvalidPriority
	}

	// check data keys...
	for key := range m.Data {
		for _, rkey := range reservedDataKeys {
			if key == rkey {
				return errReservedDataKey
			}
		}

		if strings.HasPrefix(key, "google") || strings.HasPrefix(key, "gcm") {
			return errReservedDataKeyPrefix
		}
	}
	//

	return nil
}
