package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
)

func IsFollowing(follower *model.User, publishers []string) ([]bool, error) {
	if follower == nil || len(publishers) == 0 {
		return make([]bool, len(publishers)), nil
	}

	publisherSet := make(map[string]bool)
	for _, publisher := range publishers {
		publisherSet[publisher] = true
	}

	keys := make([]AWSObject, 0, len(publisherSet))
	for publisher := range publisherSet {
		keys = append(keys, AWSObject{
			"Follower":  StringValue(follower.Username),
			"Publisher": StringValue(publisher),
		})
	}

	batchGetFollows := dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			FollowTableName: {
				Keys:                 keys,
				ProjectionExpression: aws.String("Publisher"),
			},
		},
	}

	responses, err := BatchGetItems(&batchGetFollows, len(publisherSet))
	if err != nil {
		return nil, err
	}

	followingUser := make(map[string]bool)

	for _, response := range responses {
		for _, items := range response {
			for _, item := range items {
				follow := model.Follow{}
				err = dynamodbattribute.UnmarshalMap(item, &follow)
				if err != nil {
					return nil, err
				}

				followingUser[follow.Publisher] = true
			}
		}
	}

	following := make([]bool, 0, len(publishers))
	for _, username := range publishers {
		following = append(following, followingUser[username])
	}

	return following, nil
}

func Follow(follower string, publisher string) error {
	follow := model.Follow{
		Follower:  follower,
		Publisher: publisher,
	}

	item, err := dynamodbattribute.MarshalMap(follow)
	if err != nil {
		return err
	}

	putFollow := dynamodb.PutItemInput{
		TableName: aws.String(FollowTableName),
		Item:      item,
	}

	_, err = DynamoDB().PutItem(&putFollow)

	return err
}

func Unfollow(follower string, publisher string) error {
	follow := model.Follow{
		Follower:  follower,
		Publisher: publisher,
	}

	item, err := dynamodbattribute.MarshalMap(follow)
	if err != nil {
		return err
	}

	deleteFollow := dynamodb.DeleteItemInput{
		TableName: aws.String(FollowTableName),
		Key:       item,
	}

	_, err = DynamoDB().DeleteItem(&deleteFollow)

	return err
}
