package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"fmt"
)

// Slack is used for send IncomingWebHook request
type Slack struct {
	Text string `json:"text"`
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

	if err != nil {
		fmt.Println(err.Error())
	}
	
	return err
}

func (slack Slack) json() (string) {
	jsonStr, _ := json.Marshal(slack)
	return string(jsonStr)
}