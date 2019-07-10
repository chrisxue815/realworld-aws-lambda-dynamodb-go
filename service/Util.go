package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"reflect"
	"strconv"
	"strings"
)

type AWSObject = map[string]*dynamodb.AttributeValue

func StringKey(name, value string) AWSObject {
	return AWSObject{
		name: StringValue(value),
	}
}

func StringValue(value string) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		S: aws.String(value),
	}
}

func IntKey(name string, value int) AWSObject {
	return AWSObject{
		name: IntValue(value),
	}
}

func IntValue(value int) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		N: aws.String(strconv.Itoa(value)),
	}
}

func Int64Key(name string, value int64) AWSObject {
	return AWSObject{
		name: Int64Value(value),
	}
}

func Int64Value(value int64) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		N: aws.String(strconv.FormatInt(value, 10)),
	}
}

func BlobValue(value []byte) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		B: value,
	}
}

func ReverseIndexInt64(values []int64) map[int64]int {
	indices := make(map[int64]int)
	for i, v := range values {
		indices[v] = i
	}
	return indices
}

func IsUpdateBuilderEmpty(update expression.UpdateBuilder) bool {
	return reflect.ValueOf(&update).Elem().FieldByName("operationList").IsNil()
}

func IsConditionalCheckFailed(err error) bool {
	aerr, ok := err.(awserr.Error)
	if !ok {
		return false
	}

	switch aerr.Code() {
	case dynamodb.ErrCodeConditionalCheckFailedException:
		return true
	case dynamodb.ErrCodeTransactionCanceledException:
		// There should be a better way to do this.
		// https://github.com/aws/aws-sdk-go/issues/2318
		// "If using Java, DynamoDB lists the cancellation reasons on the CancellationReasons
		// property. This property is not set for other languages."
		// https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#DynamoDB.TransactWriteItems
		return strings.Contains(aerr.Message(), "ConditionalCheckFailed")
	default:
		return false
	}
}
