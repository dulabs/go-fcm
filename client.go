// Golang client library for Firebase Cloud Messaging.
// It uses Firebase Cloud Messaging HTTP protocol: https://firebase.google.com/docs/cloud-messaging/http-server-ref.
package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	// DefaultTimeout duration in second
	DefaultTimeout time.Duration = 30 * time.Second

	API_SEND_MESSAGE = "https://fcm.googleapis.com/fcm/send"

	API_APP_INSTANCE = "https://iid.googleapis.com/iid/info/%s?details=true"

	API_USER_LOOKUP = "https://identitytoolkit.googleapis.com/v1/accounts:lookup?key=%s"

	API_TOPIC_BATCH_ADD = "https://iid.googleapis.com/iid/v1:batchAdd"

	API_TOPIC_BATCH_REMOVE = "https://iid.googleapis.com/iid/v1:batchRemove"
)

var (
	// ErrInvalidAPIKey occurs if API key is not set.
	ErrInvalidAPIKey = errors.New("client API Key is invalid")

	ErrInvalidServerKey = errors.New("client Server Key is invalid")
)

func NewClient(apiKey, serverKey string, opts ...Option) (*Client, error) {

	if apiKey == "" {
		return nil, ErrInvalidAPIKey
	}

	if serverKey == "" {
		return nil, ErrInvalidServerKey
	}

	client := &Client{
		ApiKey:    apiKey,
		ServerKey: serverKey,
		timeout:   DefaultTimeout,
	}

	for _, o := range opts {
		if err := o(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

type Client struct {
	// API key
	ApiKey string
	// Server key
	// Firebase console > Project Settings > Cloud Messaging > Server key
	ServerKey string

	timeout time.Duration
}

// Get information about app instances
//
// https://developers.google.com/instance-id/reference/server#create_a_relation_mapping_for_an_app_instance
func (c *Client) GetDeviceInfo(token string) (*DeviceInfoResponse, error) {

	endpoint := fmt.Sprintf(API_APP_INSTANCE, token)
	resp, err := c.request(http.MethodGet, endpoint, nil)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	result := DeviceInfoResponse{}
	// If you have an io.Reader, use Decoder. Otherwise use Unmarshal.
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Get user info
//
// https://firebase.google.com/docs/reference/rest/auth
func (c *Client) GetUser(token string) (*User, error) {

	endpoint := fmt.Sprintf(API_USER_LOOKUP, c.ApiKey)

	params := map[string]string{
		"idToken": token,
	}

	body, _ := json.Marshal(params)
	resp, err := c.request(http.MethodPost, endpoint, body)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	result := LookupUserResponse{}

	// If you have an io.Reader, use Decoder. Otherwise use Unmarshal.
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil && len(result.Users) == 0 {
		return nil, err
	}

	return &result.Users[0], nil
}

// Send sends a Message to Firebase Cloud Messaging.
//
// The Message must specify exactly one of Token, Topic and Condition fields.
// FCM will customize the message for each target platform based on the arguments
// specified in the Message.
// Use Legacy HTTP Server Protocol
//
// https://firebase.google.com/docs/cloud-messaging/http-server-ref
func (c *Client) Send(message *Message) (*MessageResponse, error) {

	// validate
	if err := message.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	resp, err := c.request(http.MethodPost, API_SEND_MESSAGE, body)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	result := MessageResponse{}

	// If you have an io.Reader, use Decoder. Otherwise use Unmarshal.
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// SubscribeToTopic subscribes a list of registration tokens to a topic.
// The tokens list must not be empty, and have at most 1000 tokens.
func (c *Client) SubscribeToTopic(tokens []string, topic string) error {
	return c.subscribe(API_TOPIC_BATCH_ADD, tokens, topic)
}

// UnsubscribeFromTopic unsubscribes a list of registration tokens from a topic.
// The tokens list must not be empty, and have at most 1000 tokens.
func (c *Client) UnsubscribeFromTopic(tokens []string, topic string) error {
	return c.subscribe(API_TOPIC_BATCH_REMOVE, tokens, topic)
}

///////////////////////////////////////////////////////////////////

func (c *Client) request(method string, endpoint string, body []byte) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.requestWithContext(ctx, method, endpoint, body)
}

func (c *Client) requestWithContext(ctx context.Context, method string, endpoint string, body []byte) (*http.Response, error) {
	var req *http.Request
	var resp *http.Response
	var err error

	req, err = http.NewRequest(method, endpoint, bytes.NewBuffer(body))
	if err == nil {

		req = req.WithContext(ctx)
		// Add headers
		req.Header.Add("Authorization", fmt.Sprintf("key=%s", c.ServerKey))
		req.Header.Add("Content-Type", "application/json")

		// Execute
		client := &http.Client{}
		resp, err = client.Do(req)
	}

	return resp, err
}

func (c *Client) subscribe(endpoint string, tokens []string, topic string) error {
	data := TopicSubscribe{
		RegistrationTokens: tokens,
		To:                 topic,
	}

	body, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	resp, err := c.request(http.MethodPost, endpoint, body)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return nil
}
