package notification

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/fabiokaelin/f-oauth/config"
	"github.com/fabiokaelin/f-oauth/pkg/user"
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

	url1 := "https://wirepusher.com/send?id=" + config.NotificationID + "&title=" + url.QueryEscape(n.Title)

	if n.Message != "" {
		url1 += "&message=" + url.QueryEscape(n.Message)
	}
	if n.Type != "" {
		url1 += "&type=" + url.QueryEscape(n.Type)
	}
	if n.Action != "" {
		url1 += "&action=" + url.QueryEscape(n.Action)
	}
	if n.ImageURL != "" {
		url1 += "&image=" + url.QueryEscape(n.ImageURL)
	}

	// fmt.Println(url1)

	request, err := http.NewRequest("GET", url1, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("status code is not 200")
	}

	return nil
}

func NewUserNotification(newUser user.User) error {
	notificationConfig := Config{
		Title:   fmt.Sprintf("New User: %s", newUser.Name),
		Message: fmt.Sprintf("Provider: %s\nEmail: %s", newUser.Provider, newUser.Email),
		Type:    "newuser",
	}

	err := notificationConfig.Send()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
