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
)
