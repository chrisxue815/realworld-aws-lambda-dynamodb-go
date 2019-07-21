package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
)

func IsFollowing(user *model.User, usernames []string) ([]bool, error) {
	if user == nil || len(usernames) == 0 {
		return make([]bool, len(usernames)), nil
	}

	usernameSet := make(map[string]bool)
	for _, username := range usernames {
		usernameSet[username] = true
	}

	keys := make([]AWSObject, 0, len(usernameSet))
	for author := range usernameSet {
		keys = append(keys, AWSObject{
			"Follower":  StringValue(user.Username),
			"Publisher": StringValue(author),
		})
	}

	batchGetFollows := dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			FollowTableName.Get(): {
				Keys:                 keys,
				ProjectionExpression: aws.String("Publisher"),
			},
		},
	}

	responses, err := BatchGetItems(&batchGetFollows, len(usernameSet))
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

	following := make([]bool, 0, len(usernames))
	for _, username := range usernames {
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
		TableName: aws.String(FollowTableName.Get()),
		Item:      item,
	}

	_, err = DynamoDB().PutItem(&putFollow)

	return err
}
