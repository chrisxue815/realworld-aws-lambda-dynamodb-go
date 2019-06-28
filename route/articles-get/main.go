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

type ResponseBody struct {
	Articles      []ArticleResponse `json:"articles"`
	ArticlesCount int               `json:"articlesCount"`
}

type ArticleResponse struct {
	Slug           string   `json:"slug"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Body           string   `json:"body"`
	TagList        []string `json:"tagList"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"UpdatedAt"`
	Favorited      bool     `json:"favorited"`
	FavoritesCount int64    `json:"favoritesCount"`
	Author         Author   `json:"author"`
}

type Author struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

func Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, _, _ := service.GetCurrentUser(request.Headers["Authorization"])

	offset, err := strconv.Atoi(request.QueryStringParameters["offset"])
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(request.QueryStringParameters["limit"])
	if err != nil {
		limit = 20
	}

	author := request.QueryStringParameters["author"]
	tag := request.QueryStringParameters["tag"]
	favorited := request.QueryStringParameters["favorited"]

	articles, err := service.GetArticles(offset, limit, author, tag, favorited)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	isFavorited, err := service.IsArticleFavoritedByUser(user, articles)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	authors, err := service.GetArticleAuthors(articles)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	following, err := service.IsFollowing(user, articles)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	articlesResponse := make([]ArticleResponse, 0, len(articles))

	for i, article := range articles {
		articlesResponse = append(articlesResponse, ArticleResponse{
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        article.TagList,
			Slug:           article.Slug,
			CreatedAt:      time.Unix(0, article.CreatedAt).Format(model.TimestampFormat),
			UpdatedAt:      time.Unix(0, article.UpdatedAt).Format(model.TimestampFormat),
			Favorited:      isFavorited[i],
			FavoritesCount: article.FavoritesCount,
			Author: Author{
				Username:  authors[i].Username,
				Bio:       authors[i].Bio,
				Image:     authors[i].Image,
				Following: following[i],
			},
		})
	}

	responseBody := ResponseBody{
		Articles:      articlesResponse,
		ArticlesCount: len(articlesResponse),
	}

	return util.NewSuccessResponse(200, responseBody)
}

func main() {
	lambda.Start(Handle)
}
