package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
)

func IsFollowingArticleAuthor(user *model.User, articles []model.Article) ([]bool, error) {
	if user == nil || len(articles) == 0 {
		return make([]bool, len(articles)), nil
	}

	authors := make(map[string]bool)
	for _, article := range articles {
		authors[article.Author] = true
	}

	keys := make([]map[string]*dynamodb.AttributeValue, 0, len(authors))
	for author := range authors {
		keys = append(keys, map[string]*dynamodb.AttributeValue{
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

	responses, err := BatchGetItems(&batchGetFollows, len(authors))
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

	following := make([]bool, 0, len(articles))
	for _, article := range articles {
		following = append(following, followingUser[article.Author])
	}

	return following, nil
}
