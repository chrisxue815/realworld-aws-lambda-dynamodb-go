package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/service"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/util"
	"time"
)

type Request struct {
	Article ArticleRequest `json:"article"`
}

type ArticleRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList"`
}

type Response struct {
	Article ArticleResponse `json:"article"`
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

	request := Request{}
	err = json.Unmarshal([]byte(input.Body), &request)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	now := time.Now().UTC()
	nowUnixNano := now.UnixNano()
	nowStr := now.Format(model.TimestampFormat)

	article := model.Article{
		Title:       request.Article.Title,
		Description: request.Article.Description,
		Body:        request.Article.Body,
		TagList:     request.Article.TagList,
		CreatedAt:   nowUnixNano,
		UpdatedAt:   nowUnixNano,
		Author:      user.Username,
	}

	err = service.PutArticle(&article)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	response := Response{
		Article: ArticleResponse{
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        request.Article.TagList,
			Slug:           article.Slug,
			CreatedAt:      nowStr,
			UpdatedAt:      nowStr,
			Favorited:      false,
			FavoritesCount: 0,
			Author: AuthorResponse{
				Username:  user.Username,
				Bio:       user.Bio,
				Image:     user.Image,
				Following: false,
			},
		},
	}

	return util.NewSuccessResponse(201, response)
}

func main() {
	lambda.Start(Handle)
}
