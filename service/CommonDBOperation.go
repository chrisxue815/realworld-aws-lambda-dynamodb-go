package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/util"
)

func GetItemByKey(tableName string, key map[string]*dynamodb.AttributeValue, out interface{}) error {
	input := dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
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

func QueryItems(queryInput *dynamodb.QueryInput, offset, cap int) ([]map[string]*dynamodb.AttributeValue, error) {
	items := make([]map[string]*dynamodb.AttributeValue, 0, cap)
	resultIndex := 0

	err := DynamoDB().QueryPages(queryInput, func(page *dynamodb.QueryOutput, lastPage bool) bool {
		pageCount := len(page.Items)

		if resultIndex+pageCount > offset {
			start := util.MaxInt(0, offset-resultIndex)
			for i := start; i < pageCount; i++ {
				items = append(items, page.Items[i])
			}
		}

		resultIndex += pageCount
		return true
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}

func BatchGetItems(batchGetInput *dynamodb.BatchGetItemInput, cap int) ([]map[string][]map[string]*dynamodb.AttributeValue, error) {
	responses := make([]map[string][]map[string]*dynamodb.AttributeValue, 0, cap)

	err := DynamoDB().BatchGetItemPages(batchGetInput, func(page *dynamodb.BatchGetItemOutput, lastPage bool) bool {
		responses = append(responses, page.Responses)
		return true
	})

	if err != nil {
		return nil, err
	}

	return responses, nil
}
