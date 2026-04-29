package service

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var once sync.Once
var svc *dynamodb.Client

func initializeSingletons() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	svc = dynamodb.NewFromConfig(cfg)
}

func DynamoDB() *dynamodb.Client {
	once.Do(initializeSingletons)
	return svc
}
