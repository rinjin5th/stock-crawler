package main

import (
	"os"
)

var (
	// AwsAccessKey is AWS_ACCESS_KEY in environment variable
	AwsAccessKey = os.Getenv("ACCESS_KEY")
	// AwsSecretKey is AWS_SECRET_KEY in environment variable
	AwsSecretKey = os.Getenv("SECRET_KEY")
)

const (
	// AwsRegion is the default connection destination
	AwsRegion    = "us-east-1"
	// LowerLimit is threshold to trigger an alert
	LowerLimit = -9
	// UpperLimit is threshold to trigger an notify
	UpperLimit = 6
	// SlackWebHookURL is alert send destination
	SlackWebHookURL = "https://hooks.slack.com/services/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"
)
