package service

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
)

func GetFavoriteArticleIdsByUsername(username string, offset, limit int) ([]int64, error) {
	queryArticleIds := dynamodb.QueryInput{
		TableName:                 aws.String(FavoriteArticleTableName),
		IndexName:                 aws.String("FavoritedAt"),
		KeyConditionExpression:    aws.String("Username=:username"),
		ExpressionAttributeValues: StringKey(":username", username),
		Limit:                     aws.Int32(int32(offset + limit)),
		ScanIndexForward:          aws.Bool(false),
		ProjectionExpression:      aws.String("ArticleId"),
	}

	items, err := QueryItems(&queryArticleIds, offset, limit)
	if err != nil {
		return nil, err
	}

	favoriteArticles := make([]model.FavoriteArticle, len(items))
	err = attributevalue.UnmarshalListOfMaps(items, &favoriteArticles)
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

	keys := make([]AWSObject, 0, len(articles))
	for _, article := range articles {
		keys = append(keys, AWSObject{
			"Username":  StringValue(user.Username),
			"ArticleId": Int64Value(article.ArticleId),
		})
	}

	batchGetFavoriteArticles := dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			FavoriteArticleTableName: {
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
				err = attributevalue.UnmarshalMap(item, &favoriteArticle)
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

func SetFavoriteArticle(favoriteArticle model.FavoriteArticle) error {
	item, err := attributevalue.MarshalMap(favoriteArticle)
	if err != nil {
		return err
	}

	transactItems := make([]types.TransactWriteItem, 0, 2)

	// Favorite the article
	transactItems = append(transactItems, types.TransactWriteItem{
		Put: &types.Put{
			TableName:           aws.String(FavoriteArticleTableName),
			Item:                item,
			ConditionExpression: aws.String("attribute_not_exists(Username) AND attribute_not_exists(ArticleId)"),
		},
	})

	// Update favorites count
	transactItems = append(transactItems, types.TransactWriteItem{
		Update: &types.Update{
			TableName:                 aws.String(ArticleTableName),
			Key:                       Int64Key("ArticleId", favoriteArticle.ArticleId),
			ConditionExpression:       aws.String("attribute_exists(ArticleId)"),
			UpdateExpression:          aws.String("ADD FavoritesCount :one"),
			ExpressionAttributeValues: IntKey(":one", 1),
		},
	})

	_, err = DynamoDB().TransactWriteItems(context.TODO(), &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})

	if err != nil {
		return model.NewInputError("slug", "not found or already favorited")
	}

	return nil
}

func UnfavoriteArticle(favoriteArticle model.FavoriteArticleKey) error {
	item, err := attributevalue.MarshalMap(favoriteArticle)
	if err != nil {
		return err
	}

	transactItems := make([]types.TransactWriteItem, 0, 2)

	// Unfavorite the article
	transactItems = append(transactItems, types.TransactWriteItem{
		Delete: &types.Delete{
			TableName:           aws.String(FavoriteArticleTableName),
			Key:                 item,
			ConditionExpression: aws.String("attribute_exists(Username) AND attribute_exists(ArticleId)"),
		},
	})

	// Update favorites count
	transactItems = append(transactItems, types.TransactWriteItem{
		Update: &types.Update{
			TableName:                 aws.String(ArticleTableName),
			Key:                       Int64Key("ArticleId", favoriteArticle.ArticleId),
			ConditionExpression:       aws.String("attribute_exists(ArticleId)"),
			UpdateExpression:          aws.String("ADD FavoritesCount :minus_one"),
			ExpressionAttributeValues: IntKey(":minus_one", -1),
		},
	})

	_, err = DynamoDB().TransactWriteItems(context.TODO(), &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})

	if err != nil {
		return model.NewInputError("slug", "not found or not favorited")
	}

	return nil
}
