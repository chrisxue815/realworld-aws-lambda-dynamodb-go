package service

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
		RequestItems: map[string]types.KeysAndAttributes{
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
				err = attributevalue.UnmarshalMap(item, &follow)
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

	item, err := attributevalue.MarshalMap(follow)
	if err != nil {
		return err
	}

	putFollow := dynamodb.PutItemInput{
		TableName: aws.String(FollowTableName),
		Item:      item,
	}

	_, err = DynamoDB().PutItem(context.TODO(), &putFollow)
	if err != nil {
		fmt.Println(string(debug.Stack()))
		return err
	}

	return nil
}

func Unfollow(follower string, publisher string) error {
	follow := model.Follow{
		Follower:  follower,
		Publisher: publisher,
	}

	item, err := attributevalue.MarshalMap(follow)
	if err != nil {
		return err
	}

	deleteFollow := dynamodb.DeleteItemInput{
		TableName: aws.String(FollowTableName),
		Key:       item,
	}

	_, err = DynamoDB().DeleteItem(context.TODO(), &deleteFollow)
	if err != nil {
		fmt.Println(string(debug.Stack()))
		return err
	}

	return nil
}
