package main

import (
	"encoding/json"
	"net/http"
    "net/url"
)

// Slack is used for send IncomingWebHook request
type Slack struct {
	Text string
}

// NewSlack creates slack object
func NewSlack(text string) *Slack {
	slack := new(Slack)
	slack.Text = text
	return slack
}

// Send is used for send IncomingWebHook request
func (slack Slack) Send(webHookURL string) (error) {
	_, err := http.PostForm(
        webHookURL,
        url.Values{"payload": {slack.json()}},
	)
	
	return err
}

func (slack Slack) json() (string) {
	jsonStr, _ := json.Marshal(slack)
	return string(jsonStr)
}