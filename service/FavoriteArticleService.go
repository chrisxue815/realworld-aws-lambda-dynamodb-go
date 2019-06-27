package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
)

func GetFavoriteArticleIdsByUsername(username string, offset, limit int) ([]int64, error) {
	queryArticleIds := dynamodb.QueryInput{
		TableName: aws.String(FavoriteArticleTableName.Get()),
		IndexName: aws.String("FavoritedAt"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {
				S: aws.String(username),
			},
		},
		KeyConditionExpression: aws.String("Username=:username"),
		Limit:                  aws.Int64(int64(offset + limit)),
		ScanIndexForward:       aws.Bool(false),
		ProjectionExpression:   aws.String("ArticleId"),
	}

	items, err := QueryItems(&queryArticleIds, offset, limit)
	if err != nil {
		return nil, err
	}

	favoriteArticles := make([]model.FavoriteArticle, len(items))
	err = dynamodbattribute.UnmarshalListOfMaps(items, &favoriteArticles)
	if err != nil {
		return nil, err
	}

	articleIds := make([]int64, 0, len(items))

	for _, favoriteArticle := range favoriteArticles {
		articleIds = append(articleIds, favoriteArticle.ArticleId)
	}

	return articleIds, nil
}

func IsArticleFavoritedByUser(user *model.User, articles []model.Article) ([]bool, error) {
	if user == nil || len(articles) == 0 {
		return make([]bool, len(articles)), nil
	}

	keys := make([]map[string]*dynamodb.AttributeValue, 0, len(articles))
	for _, article := range articles {
		keys = append(keys, map[string]*dynamodb.AttributeValue{
			"Username":  StringValue(user.Username),
			"ArticleId": Int64Value(article.ArticleId),
		})
	}

	batchGetFavoriteArticles := dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			FavoriteArticleTableName.Get(): {
				Keys:                 keys,
				ProjectionExpression: aws.String("ArticleId"),
			},
		},
	}

	responses, err := BatchGetItems(&batchGetFavoriteArticles, len(articles))
	if err != nil {
		return nil, err
	}

	isFavorited := make([]bool, len(articles))
	articleIdToIndex := reverseIndexArticleIds(articles)

	for _, response := range responses {
		for _, items := range response {
			for _, item := range items {
				favoriteArticle := model.FavoriteArticle{}
				err = dynamodbattribute.UnmarshalMap(item, &favoriteArticle)
				if err != nil {
					return nil, err
				}

				index := articleIdToIndex[favoriteArticle.ArticleId]
				isFavorited[index] = true
			}
		}
	}

	return isFavorited, nil
}

func reverseIndexArticleIds(articles []model.Article) map[int64]int {
	indices := make(map[int64]int)
	for i, article := range articles {
		indices[article.ArticleId] = i
	}
	return indices
}
