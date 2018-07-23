package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"fmt"
	"io/ioutil"
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
	resp, err := http.PostForm(
        webHookURL,
        url.Values{"payload": {slack.json()}},
	)

	if err != nil {
		fmt.Println(err.Error())
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	
	fmt.Println(string(body))
	
	return err
}

func (slack Slack) json() (string) {
	jsonStr, _ := json.Marshal(slack)
	return string(jsonStr)
}