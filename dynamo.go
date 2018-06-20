package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

// NewTable creates DynamoDB table object
func NewTable(tableName string) (dynamo.Table) {
	db := dynamo.New(session.New(), NewConf())
	return db.Table(tableName)
}