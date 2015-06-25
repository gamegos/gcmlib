package gcm

type Message struct {
	// Targets
	To              string   `json:"to,omitempty"`
	RegistrationIDs []string `json:"registration_ids"`

	// Options
	CollapseKey           string `json:"collapse_key,omitempty"`
	Priority              uint8  `json:"priority,omitempty"`
	ContentAvailable      bool   `json:"content_available,omitempty"`
	DelayWhileIdle        bool   `json:"delay_while_idle,omitempty"`
	TTL                   int    `json:"time_to_live,omitempty"`
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
	titleLocKey  string `json:"title_loc_key,omitempty"`
}
