package service

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/util"
)

func GetItemByKey(tableName string, key AWSObject, out interface{}) (bool, error) {
	input := dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}

	output, err := DynamoDB().GetItem(context.TODO(), &input)
	if err != nil {
		return false, err
	}

	if len(output.Item) == 0 {
		return false, nil
	}

	err = attributevalue.UnmarshalMap(output.Item, out)
	if err != nil {
		return false, err
	}

	return true, nil
}

func QueryItems(queryInput *dynamodb.QueryInput, offset, cap int) ([]AWSObject, error) {
	items := make([]AWSObject, 0, cap)
	resultIndex := 0

	paginator := dynamodb.NewQueryPaginator(DynamoDB(), queryInput)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		pageCount := len(page.Items)
		if resultIndex+pageCount > offset {
			start := util.MaxInt(0, offset-resultIndex)
			for i := start; i < pageCount; i++ {
				items = append(items, page.Items[i])
			}
		}
		resultIndex += pageCount
	}

	return items, nil
}

func BatchGetItems(batchGetInput *dynamodb.BatchGetItemInput, cap int) ([]map[string][]AWSObject, error) {
	responses := make([]map[string][]AWSObject, 0, cap)

	input := batchGetInput
	for {
		output, err := DynamoDB().BatchGetItem(context.TODO(), input)
		if err != nil {
			return nil, err
		}
		responses = append(responses, output.Responses)
		if len(output.UnprocessedKeys) == 0 {
			break
		}
		input = &dynamodb.BatchGetItemInput{
			RequestItems: output.UnprocessedKeys,
		}
	}

	return responses, nil
}
