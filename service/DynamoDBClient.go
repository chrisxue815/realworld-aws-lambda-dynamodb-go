package service

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"sync"
)

var once sync.Once
var svc *dynamodb.DynamoDB

func initializeSingletons() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc = dynamodb.New(sess)
}

func DynamoDB() *dynamodb.DynamoDB {
	once.Do(initializeSingletons)
	return svc
}
