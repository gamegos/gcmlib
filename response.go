package gcm

type result struct {
	Error          resultError `json:"error"`
	MessageID      string      `json:"message_id"`
	RegistrationID string      `json:"registration_id"`
}

func (r *result) Failed() bool {
	return r.Error != ""
}

func (r *result) Success() bool {
	return !r.Failed()
}

func (r *result) TokenChanged() bool {
	return r.RegistrationID != ""
}

type response struct {
	MulticastID  int64    `json:"multicast_id"`
	Success      uint     `json:"success"`
	Failure      uint     `json:"failure"`
	CanonicalIDs uint     `json:"canonical_ids"`
	Results      []result `json:"results"`
}
