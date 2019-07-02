package service

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/util"
	"strconv"
	"strings"
)

func PutArticle(article *model.Article) error {
	err := article.Validate()
	if err != nil {
		return err
	}

	const maxAttempt = 5

	// Try to find a unique article id
	for attempt := 0; ; attempt++ {
		err := putArticleWithRandomId(article)

		if err == nil {
			return nil
		}

		if attempt >= maxAttempt {
			return err
		}

		aerr, ok := err.(awserr.Error)
		if !ok || aerr.Code() != dynamodb.ErrCodeConditionalCheckFailedException {
			return err
		}

		RenewArticleIdRand()
	}
}

func putArticleWithRandomId(article *model.Article) error {
	article.ArticleId = 1 + ArticleIdRand().Int63n(model.MaxArticleId-1) // range: [1, MaxArticleId)
	article.MakeSlug()

	articleItem, err := dynamodbattribute.MarshalMap(article)
	if err != nil {
		return err
	}

	transactItems := make([]*dynamodb.TransactWriteItem, 0, 1+2*len(article.TagList))

	// Put a new article
	transactItems = append(transactItems, &dynamodb.TransactWriteItem{
		Put: &dynamodb.Put{
			TableName:           aws.String(ArticleTableName.Get()),
			Item:                articleItem,
			ConditionExpression: aws.String("attribute_not_exists(ArticleId)"),
		},
	})

	for _, tag := range article.TagList {
		articleTag := model.ArticleTag{
			Tag:       tag,
			ArticleId: article.ArticleId,
			CreatedAt: article.CreatedAt,
		}

		item, err := dynamodbattribute.MarshalMap(articleTag)
		if err != nil {
			return err
		}

		// Link article with tag
		transactItems = append(transactItems, &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				TableName: aws.String(ArticleTagTableName.Get()),
				Item:      item,
			},
		})

		// Update article count for each tag
		transactItems = append(transactItems, &dynamodb.TransactWriteItem{
			Update: &dynamodb.Update{
				TableName:                 aws.String(TagTableName.Get()),
				Key:                       StringKey("Tag", tag),
				UpdateExpression:          aws.String("ADD ArticleCount :one"),
				ExpressionAttributeValues: IntKey(":one", 1),
			},
		})
	}

	_, err = DynamoDB().TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})

	return err
}

func GetArticles(offset, limit int, author, tag, favorited string) ([]model.Article, error) {
	if offset < 0 {
		return nil, util.NewInputError("offset", "must be non-negative")
	}

	if limit <= 0 {
		return nil, util.NewInputError("limit", "must be positive")
	}

	const maxDepth = 1000
	if offset+limit > maxDepth {
		return nil, util.NewInputError("offset + limit", fmt.Sprintf("must be smaller or equal to %d", maxDepth))
	}

	numFilters := getNumFilters(author, tag, favorited)
	if numFilters > 1 {
		return nil, util.NewInputError("author, tag, favorited", "only one of these can be specified")
	}

	if numFilters == 0 {
		return getAllArticles(offset, limit)
	}

	if author != "" {
		return getArticlesByAuthor(author, offset, limit)
	}

	if tag != "" {
		return getArticlesByTag(tag, offset, limit)
	}

	if favorited != "" {
		return getFavoriteArticlesByUsername(favorited, offset, limit)
	}

	return nil, errors.New("unreachable code")
}

func getNumFilters(author, tag, favorited string) int {
	numFilters := 0
	if author != "" {
		numFilters++
	}
	if tag != "" {
		numFilters++
	}
	if favorited != "" {
		numFilters++
	}
	return numFilters
}

func getAllArticles(offset, limit int) ([]model.Article, error) {
	queryArticles := dynamodb.QueryInput{
		TableName:                 aws.String(ArticleTableName.Get()),
		IndexName:                 aws.String("CreatedAt"),
		KeyConditionExpression:    aws.String("Dummy=:zero"),
		ExpressionAttributeValues: IntKey(":zero", 0),
		Limit:                     aws.Int64(int64(offset + limit)),
		ScanIndexForward:          aws.Bool(false),
	}

	items, err := QueryItems(&queryArticles, offset, limit)
	if err != nil {
		return nil, err
	}

	articles := make([]model.Article, len(items))
	err = dynamodbattribute.UnmarshalListOfMaps(items, &articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func getArticlesByAuthor(author string, offset, limit int) ([]model.Article, error) {
	queryArticles := dynamodb.QueryInput{
		TableName:                 aws.String(ArticleTableName.Get()),
		IndexName:                 aws.String("Author"),
		KeyConditionExpression:    aws.String("Author=:author"),
		ExpressionAttributeValues: StringKey(":author", author),
		Limit:                     aws.Int64(int64(offset + limit)),
		ScanIndexForward:          aws.Bool(false),
	}

	items, err := QueryItems(&queryArticles, offset, limit)
	if err != nil {
		return nil, err
	}

	articles := make([]model.Article, len(items))
	err = dynamodbattribute.UnmarshalListOfMaps(items, &articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func getArticlesByTag(tag string, offset, limit int) ([]model.Article, error) {
	articleIds, err := GetArticleIdsByTag(tag, offset, limit)
	if err != nil {
		return nil, err
	}

	return getArticlesByArticleIds(articleIds, limit)
}

func getFavoriteArticlesByUsername(username string, offset, limit int) ([]model.Article, error) {
	articleIds, err := GetFavoriteArticleIdsByUsername(username, offset, limit)
	if err != nil {
		return nil, err
	}

	return getArticlesByArticleIds(articleIds, limit)
}

func getArticlesByArticleIds(articleIds []int64, limit int) ([]model.Article, error) {
	if len(articleIds) == 0 {
		return make([]model.Article, 0), nil
	}

	keys := make([]AWSObject, 0, len(articleIds))
	for _, articleId := range articleIds {
		keys = append(keys, Int64Key("ArticleId", articleId))
	}

	batchGetArticles := dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			ArticleTableName.Get(): {
				Keys: keys,
			},
		},
	}

	responses, err := BatchGetItems(&batchGetArticles, limit)
	if err != nil {
		return nil, err
	}

	articles := make([]model.Article, len(articleIds))
	articleIdToIndex := ReverseIndexInt64(articleIds)

	for _, response := range responses {
		for _, items := range response {
			for _, item := range items {
				article := model.Article{}
				err = dynamodbattribute.UnmarshalMap(item, &article)
				if err != nil {
					return nil, err
				}

				index := articleIdToIndex[article.ArticleId]
				articles[index] = article
			}
		}
	}

	return articles, nil
}

func GetArticleRelatedProperties(user *model.User, articles []model.Article) ([]bool, []model.User, []bool, error) {
	isFavorited, err := IsArticleFavoritedByUser(user, articles)
	if err != nil {
		return nil, nil, nil, err
	}

	authors, err := GetArticleAuthors(articles)
	if err != nil {
		return nil, nil, nil, err
	}

	following, err := IsFollowingArticleAuthor(user, articles)
	if err != nil {
		return nil, nil, nil, err
	}

	return isFavorited, authors, following, nil
}

func GetArticleBySlug(slug string) (model.Article, error) {
	dashIndex := strings.LastIndexByte(slug, '-')
	if dashIndex == -1 {
		return model.Article{}, util.NewInputError("slug", "invalid")
	}

	articleId, err := strconv.ParseInt(slug[dashIndex+1:], 16, 64)
	if err != nil {
		return model.Article{}, util.NewInputError("slug", "invalid")
	}

	article := model.Article{}
	err = GetItemByKey(ArticleTableName.Get(), Int64Key("ArticleId", articleId), &article)
	if err != nil {
		return model.Article{}, err
	}

	return article, nil
}

func UpdateArticle(oldArticle model.Article, newArticle *model.Article) error {
	err := newArticle.Validate()
	if err != nil {
		return err
	}

	newArticle.MakeSlug()

	transactItems := make([]*dynamodb.TransactWriteItem, 0, 1+2*model.MaxNumTagsPerArticle)

	oldTagSet := util.NewStringSetFromSlice(oldArticle.TagList)
	newTagSet := util.NewStringSetFromSlice(newArticle.TagList)
	oldTags := oldTagSet.Difference(newTagSet)
	newTags := newTagSet.Difference(oldTagSet)

	expr, err := buildArticleUpdateExpression(oldArticle, *newArticle, len(oldTags) != 0 || len(newTags) != 0)
	if err != nil {
		return err
	}

	// No field changed
	if expr.Update() == nil {
		return nil
	}

	// Update article
	transactItems = append(transactItems, &dynamodb.TransactWriteItem{
		Update: &dynamodb.Update{
			TableName:                 aws.String(ArticleTableName.Get()),
			Key:                       Int64Key("ArticleId", oldArticle.ArticleId),
			ConditionExpression:       aws.String("attribute_exists(ArticleId)"),
			UpdateExpression:          expr.Update(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
		},
	})

	for tag := range oldTags {
		// Unlink article from tag
		transactItems = append(transactItems, &dynamodb.TransactWriteItem{
			Delete: &dynamodb.Delete{
				TableName: aws.String(ArticleTagTableName.Get()),
				Key: AWSObject{
					"Tag":       StringValue(tag),
					"ArticleId": Int64Value(oldArticle.ArticleId),
				},
			},
		})

		// Update article count for each tag
		transactItems = append(transactItems, &dynamodb.TransactWriteItem{
			Update: &dynamodb.Update{
				TableName:                 aws.String(TagTableName.Get()),
				Key:                       StringKey("Tag", tag),
				UpdateExpression:          aws.String("ADD ArticleCount :minus_one"),
				ExpressionAttributeValues: IntKey(":minus_one", -1),
			},
		})
	}

	for tag := range newTags {
		articleTag := model.ArticleTag{
			Tag:       tag,
			ArticleId: oldArticle.ArticleId,
			CreatedAt: oldArticle.CreatedAt,
		}

		item, err := dynamodbattribute.MarshalMap(articleTag)
		if err != nil {
			return err
		}

		// Link article with tag
		transactItems = append(transactItems, &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				TableName: aws.String(ArticleTagTableName.Get()),
				Item:      item,
			},
		})

		// Update article count for each tag
		transactItems = append(transactItems, &dynamodb.TransactWriteItem{
			Update: &dynamodb.Update{
				TableName:                 aws.String(TagTableName.Get()),
				Key:                       StringKey("Tag", tag),
				UpdateExpression:          aws.String("ADD ArticleCount :one"),
				ExpressionAttributeValues: IntKey(":one", 1),
			},
		})
	}

	_, err = DynamoDB().TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})
	if err != nil {
		return err
	}

	return nil
}

func buildArticleUpdateExpression(oldArticle model.Article, newArticle model.Article, updateTagList bool) (expression.Expression, error) {
	update := expression.UpdateBuilder{}

	if oldArticle.Slug != newArticle.Slug {
		update = update.Set(expression.Name("Slug"), expression.Value(newArticle.Slug))
	}

	if oldArticle.Title != newArticle.Title {
		update = update.Set(expression.Name("Title"), expression.Value(newArticle.Title))
	}

	if oldArticle.Description != newArticle.Description {
		update = update.Set(expression.Name("Description"), expression.Value(newArticle.Description))
	}

	if oldArticle.Body != newArticle.Body {
		update = update.Set(expression.Name("Body"), expression.Value(newArticle.Body))
	}

	if updateTagList {
		update = update.Set(expression.Name("TagList"), expression.Value(newArticle.TagList))
	}

	if oldArticle.UpdatedAt != newArticle.UpdatedAt {
		update = update.Set(expression.Name("UpdatedAt"), expression.Value(newArticle.UpdatedAt))
	}

	if IsUpdateBuilderEmpty(update) {
		return expression.Expression{}, nil
	}

	builder := expression.NewBuilder().WithUpdate(update)
	return builder.Build()
}
