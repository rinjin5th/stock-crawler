package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// NewConf creates aws.Config for connecting AWS
func NewConf() (*aws.Config) {
	creds := credentials.NewStaticCredentials(AwsAccessKey, AwsSecretKey, "")
	return &aws.Config{
		Credentials: creds,
		Region: aws.String(AwsRegion),
	}
}

// NewSession creates session.Session for connecting AWS
func NewSession() (*session.Session) {
	return session.Must(session.NewSession(NewConf()))
}