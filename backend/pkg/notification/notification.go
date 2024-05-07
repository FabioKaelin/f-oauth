package notification

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/fabiokaelin/f-oauth/config"
)

// Config is the struct for a notification config
type Config struct {
	Title    string // required
	Message  string // message of notification
	Type     string // to manage notification on phone (vibrating, ignore, sound)
	Action   string // url which will be opened when user clicks on notification
	ImageURL string // image like in youtube videos
	// function to send notification
	// Send func() error
}

// Send sends a notification
func (n *Config) Send() error {
	//make a http get request to https://wirepusher.com/send

	if n.Title == "" {
		return errors.New("title is required")
	}
	if config.NotificationID == "" {
		return errors.New("notificationId is required (set in environment variable)")
	}

	url1 := "https://wirepusher.com/send?id=" + config.NotificationID + "&title=" + n.Title

	if n.Message != "" {
		url1 += "&message=" + n.Message
	}
	if n.Type != "" {
		url1 += "&type=" + n.Type
	}
	if n.Action != "" {
		url1 += "&action=" + n.Action
	}
	if n.ImageURL != "" {
		url1 += "&image=" + n.ImageURL
	}

	// encrypt url

	newURL := url.QueryEscape(url1)

	_, err := http.NewRequest("GET", newURL, nil)
	if err != nil {
		return err
	}

	return nil
}
