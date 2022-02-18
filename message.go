package fcm

import (
	"errors"
	"strings"
)

var (
	// ErrInvalidMessage occurs if push notitication message is nil.
	ErrInvalidMessage = errors.New("message is invalid")

	// ErrInvalidTarget occurs if message topic is empty.
	ErrInvalidTarget = errors.New("topic is invalid or registration ids are not set")

	// ErrToManyRegIDs occurs when registration ids more then 1000.
	ErrToManyRegIDs = errors.New("too many registrations ids")

	// ErrInvalidTimeToLive occurs if TimeToLive more then 2419200.
	ErrInvalidTimeToLive = errors.New("messages time-to-live is invalid")
)

// Downstream HTTP messages (JSON)
//
// https://firebase.google.com/docs/cloud-messaging/http-server-ref#downstream-http-messages-json
type Message struct {

	// Targets
	To              string   `json:"to,omitempty"`
	RegistrationIds []string `json:"registration_ids,omitempty"`
	Condition       string   `json:"condition,omitempty"`

	// Options
	CollapseKey           string `json:"collapse_key,omitempty"`
	Priority              string `json:"priority,omitempty"`
	ContentAvailable      bool   `json:"content_available,omitempty"`
	MutableContent        bool   `json:"mutable_content,omitempty"`
	DelayWhileIdle        bool   `json:"delay_while_idle,omitempty"`
	TimeToLive            *uint  `json:"time_to_live,omitempty"`
	RestrictedPackageName string `json:"restricted_package_name,omitempty"`
	DryRun                bool   `json:"dry_run,omitempty"`

	// Payload
	Data         map[string]string      `json:"data,omitempty"`
	Notification *Notification          `json:"notification,omitempty"`
	Apns         map[string]interface{} `json:"apns,omitempty"`
	Webpush      map[string]interface{} `json:"webpush,omitempty"`
}

// Validate returns an error if the message is not well-formed.
func (msg *Message) Validate() error {
	if msg == nil {
		return ErrInvalidMessage
	}

	// validate target identifier: `to` or `condition`, or `registration_ids`
	opCnt := strings.Count(msg.Condition, "&&") + strings.Count(msg.Condition, "||")
	if msg.To == "" && (msg.Condition == "" || opCnt > 5) && len(msg.RegistrationIds) == 0 {
		return ErrInvalidTarget
	}

	if len(msg.RegistrationIds) > 1000 {
		return ErrToManyRegIDs
	}

	if msg.TimeToLive != nil && *msg.TimeToLive > uint(2419200) {
		return ErrInvalidTimeToLive
	}
	return nil
}

type Notification struct {
	Title        string   `json:"title,omitempty"`
	Body         string   `json:"body,omitempty"`
	Sound        string   `json:"sound,omitempty"`
	ClickAction  string   `json:"click_action,omitempty"`
	BodyLocKey   string   `json:"body_loc_key,omitempty"`
	BodyLocArgs  []string `json:"body_loc_args,omitempty"`
	TitleLocKey  string   `json:"title_loc_key,omitempty"`
	TitleLocArgs []string `json:"title_loc_args,omitempty"`

	/// iOS only
	Subtitle string `json:"subtitle,omitempty"`
	Badge    string `json:"badge,omitempty"`

	/// Android only
	AndroidChannelId string `json:"android_channel_id,omitempty"`
	Icon             string `json:"icon,omitempty"`
	Tag              string `json:"tag,omitempty"`
	Color            string `json:"color,omitempty"`
}

// Downstream HTTP message response body (JSON).
// Include topic message HTTP response body.
type MessageResponse struct {
	// Message
	MulticastId int64                   `bson:"multicast_id,omitempty" json:"multicast_id,omitempty"`
	Success     int                     `bson:"success,omitempty" json:"success,omitempty"`
	Failure     int                     `bson:"failure,omitempty" json:"failure,omitempty"`
	Results     []MessageResponseResult `bson:"results,omitempty" json:"results,omitempty"`

	// Topic
	TopicMessageResponse
}

type MessageResponseResult struct {
	MessageId      string `bson:"message_id,omitempty" json:"message_id,omitempty"`
	RegistrationId string `bson:"registration_id,omitempty" json:"registration_id,omitempty"`
	Error          string `bson:"error,omitempty" json:"error,omitempty"`
}
