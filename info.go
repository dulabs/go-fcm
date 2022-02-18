package fcm

type DeviceInfoResponse struct {
	ApplicationVersion string       `json:"applicationVersion,omitempty"`
	ConnectDate        string       `json:"connectDate,omitempty"`
	Application        string       `json:"application,omitempty"`
	Scope              string       `json:"scope,omitempty"`
	AuthorizedEntity   string       `json:"authorizedEntity,omitempty"`
	ConnectionType     string       `json:"connectionType,omitempty"`
	Platform           string       `json:"platform,omitempty"`
	Rel                Relationship `json:"rel,omitempty"`
}

type Relationship struct {
	Topics map[string]RelDate `json:"topics,omitempty"`
}

type RelDate struct {
	AddDate string `json:"addDate,omitempty"`
}
