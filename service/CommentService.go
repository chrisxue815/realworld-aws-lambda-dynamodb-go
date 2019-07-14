package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
)

func PutComment(comment *model.Comment) error {
	err := comment.Validate()
	if err != nil {
		return err
	}

	const maxAttempt = 5

	// Try to find a unique comment id
	for attempt := 0; ; attempt++ {
		err := putCommentWithRandomId(comment)

		if err == nil {
			return nil
		}

		if attempt >= maxAttempt {
			return err
		}

		if !IsConditionalCheckFailed(err) {
			return err
		}

		CommentIdRand.RenewSeed()
	}
}

func putCommentWithRandomId(comment *model.Comment) error {
	comment.CommentId = 1 + CommentIdRand.Get().Int63n(model.MaxCommentId-1) // range: [1, MaxCommentId)

	commentItem, err := dynamodbattribute.MarshalMap(comment)
	if err != nil {
		return err
	}

	// Put a new article
	_, err = DynamoDB().PutItem(&dynamodb.PutItemInput{
		TableName:           aws.String(CommentTableName.Get()),
		Item:                commentItem,
		ConditionExpression: aws.String("attribute_not_exists(CommentId)"),
	})

	return err
}
