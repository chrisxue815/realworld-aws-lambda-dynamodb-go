package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/service"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/util"
	"strconv"
	"time"
)

type Response struct {
	Articles      []ArticleResponse `json:"articles"`
	ArticlesCount int               `json:"articlesCount"`
}

type ArticleResponse struct {
	Slug           string         `json:"slug"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Body           string         `json:"body"`
	TagList        []string       `json:"tagList"`
	CreatedAt      string         `json:"createdAt"`
	UpdatedAt      string         `json:"updatedAt"`
	Favorited      bool           `json:"favorited"`
	FavoritesCount int64          `json:"favoritesCount"`
	Author         AuthorResponse `json:"author"`
}

type AuthorResponse struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

func Handle(input events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, _, err := service.GetCurrentUser(input.Headers["Authorization"])
	if err != nil {
		return util.NewUnauthorizedResponse()
	}

	offset, err := strconv.Atoi(input.QueryStringParameters["offset"])
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(input.QueryStringParameters["limit"])
	if err != nil {
		limit = 20
	}

	articles, err := service.GetFeed(user.Username, offset, limit)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	isFavorited, authors, _, err := service.GetArticleRelatedProperties(user, articles, false)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	articleResponses := make([]ArticleResponse, 0, len(articles))

	for i, article := range articles {
		articleResponses = append(articleResponses, ArticleResponse{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        article.TagList,
			CreatedAt:      time.Unix(0, article.CreatedAt).Format(model.TimestampFormat),
			UpdatedAt:      time.Unix(0, article.UpdatedAt).Format(model.TimestampFormat),
			Favorited:      isFavorited[i],
			FavoritesCount: article.FavoritesCount,
			Author: AuthorResponse{
				Username:  authors[i].Username,
				Bio:       authors[i].Bio,
				Image:     authors[i].Image,
				Following: true,
			},
		})
	}

	response := Response{
		Articles:      articleResponses,
		ArticlesCount: len(articleResponses),
	}

	return util.NewSuccessResponse(200, response)
}

func main() {
	lambda.Start(Handle)
}
