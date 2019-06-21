package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

func GetItemByKey(tableName string, keyName string, keyValue string, out interface{}) error {
	input := dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(keyValue),
			},
		},
	}

	output, err := DynamoDB().GetItem(&input)
	if err != nil {
		return err
	}

	err = dynamodbattribute.UnmarshalMap(output.Item, out)
	if err != nil {
		return err
	}

	return nil
}
