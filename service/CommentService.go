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
		TableName:           aws.String(CommentTableName),
		Item:                commentItem,
		ConditionExpression: aws.String("attribute_not_exists(CommentId)"),
	})

	return err
}

func GetCommentRelatedProperties(user *model.User, comments []model.Comment) ([]model.User, []bool, error) {
	authorUsernames := make([]string, 0, len(comments))
	for _, comment := range comments {
		authorUsernames = append(authorUsernames, comment.Author)
	}

	authors, err := GetUserListByUsername(authorUsernames)
	if err != nil {
		return nil, nil, err
	}

	following, err := IsFollowing(user, authorUsernames)
	if err != nil {
		return nil, nil, err
	}

	return authors, following, nil
}

func GetComments(slug string) ([]model.Comment, error) {
	articleId, err := model.SlugToArticleId(slug)
	if err != nil {
		return nil, err
	}

	queryComments := dynamodb.QueryInput{
		TableName:                 aws.String(CommentTableName),
		IndexName:                 aws.String("CreatedAt"),
		KeyConditionExpression:    aws.String("ArticleId=:articleId"),
		ExpressionAttributeValues: Int64Key(":articleId", articleId),
		ScanIndexForward:          aws.Bool(false),
	}

	const queryInitialCapacity = 16
	items, err := QueryItems(&queryComments, 0, queryInitialCapacity)
	if err != nil {
		return nil, err
	}

	comments := make([]model.Comment, len(items))
	err = dynamodbattribute.UnmarshalListOfMaps(items, &comments)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func DeleteComment(slug string, commentId int64, username string) error {
	articleId, err := model.SlugToArticleId(slug)
	if err != nil {
		return err
	}

	key := model.CommentKey{
		ArticleId: articleId,
		CommentId: commentId,
	}

	item, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		return err
	}

	deleteComment := dynamodb.DeleteItemInput{
		TableName:                 aws.String(CommentTableName),
		Key:                       item,
		ConditionExpression:       aws.String("Author=:username"),
		ExpressionAttributeValues: StringKey(":username", username),
	}

	_, err = DynamoDB().DeleteItem(&deleteComment)

	return err
}
