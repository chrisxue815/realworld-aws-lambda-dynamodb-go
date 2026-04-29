package service

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"reflect"
	"strconv"
)

type AWSObject = map[string]types.AttributeValue

func StringKey(name, value string) AWSObject {
	return AWSObject{
		name: StringValue(value),
	}
}

func StringValue(value string) types.AttributeValue {
	return &types.AttributeValueMemberS{Value: value}
}

func IntKey(name string, value int) AWSObject {
	return AWSObject{
		name: IntValue(value),
	}
}

func IntValue(value int) types.AttributeValue {
	return &types.AttributeValueMemberN{Value: strconv.Itoa(value)}
}

func Int64Key(name string, value int64) AWSObject {
	return AWSObject{
		name: Int64Value(value),
	}
}

func Int64Value(value int64) types.AttributeValue {
	return &types.AttributeValueMemberN{Value: strconv.FormatInt(value, 10)}
}

func BlobValue(value []byte) types.AttributeValue {
	return &types.AttributeValueMemberB{Value: value}
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
	var condCheckFailed *types.ConditionalCheckFailedException
	if errors.As(err, &condCheckFailed) {
		return true
	}

	var txnCancelled *types.TransactionCanceledException
	if errors.As(err, &txnCancelled) {
		for _, reason := range txnCancelled.CancellationReasons {
			if aws.ToString(reason.Code) == "ConditionalCheckFailed" {
				return true
			}
		}
	}

	return false
}
