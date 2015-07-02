package gcm

type Result struct {
	Error          ResultError `json:"error"`
	MessageID      string      `json:"message_id"`
	RegistrationID string      `json:"registration_id"`
}

func (r *Result) Failed() bool {
	return r.Error != ""
}

func (r *Result) Success() bool {
	return !r.Failed()
}

func (r *Result) TokenChanged() bool {
	return r.RegistrationID != ""
}

type Response struct {
	MulticastID  int64    `json:"multicast_id"`
	Success      uint     `json:"success"`
	Failure      uint     `json:"failure"`
	CanonicalIDs uint     `json:"canonical_ids"`
	Results      []Result `json:"results"`
}
