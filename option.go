package fcm

import (
	"errors"
	"time"
)

// Option configurates Client with defined option.
type Option func(*Client) error

// WithTimeout returns Option to configure HTTP Client timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) error {
		if d.Nanoseconds() <= 0 {
			return errors.New("invalid timeout duration")
		}
		c.timeout = d
		return nil
	}
}
