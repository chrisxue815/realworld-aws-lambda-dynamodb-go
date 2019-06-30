package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
)

func GetArticleIdsByTag(tag string, offset, limit int) ([]int64, error) {
	queryArticleIds := dynamodb.QueryInput{
		TableName:                 aws.String(ArticleTagTableName.Get()),
		IndexName:                 aws.String("CreatedAt"),
		KeyConditionExpression:    aws.String("Tag=:tag"),
		ExpressionAttributeValues: StringKey(":tag", tag),
		Limit:                     aws.Int64(int64(offset + limit)),
		ScanIndexForward:          aws.Bool(false),
		ProjectionExpression:      aws.String("ArticleId"),
	}

	items, err := QueryItems(&queryArticleIds, offset, limit)
	if err != nil {
		return nil, err
	}

	articleTags := make([]model.ArticleTag, len(items))
	err = dynamodbattribute.UnmarshalListOfMaps(items, &articleTags)
	if err != nil {
		return nil, err
	}

	articleIds := make([]int64, 0, len(items))

	for _, articleTag := range articleTags {
		articleIds = append(articleIds, articleTag.ArticleId)
	}

	return articleIds, nil
}
