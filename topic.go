package fcm

type TopicSubscribe struct {
	To                 string   `bson:"to,omitempty" json:"to,omitempty"`
	RegistrationTokens []string `bson:"registration_tokens,omitempty" json:"registration_tokens,omitempty"`
}

// Topic message HTTP response body (JSON)
type TopicMessageResponse struct {
	MessageId int64  `bson:"message_id,omitempty" json:"message_id,omitempty"`
	Error     string `bson:"error,omitempty" json:"error,omitempty"`
}
