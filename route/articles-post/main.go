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

type RequestBody struct {
	Article ArticleRequest `json:"article"`
}

type ArticleRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList"`
}

type ResponseBody struct {
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

func Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, _, err := service.GetCurrentUser(request.Headers["Authorization"])
	if err != nil {
		return util.NewUnauthorizedResponse()
	}

	requestBody := RequestBody{}
	err = json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	now := time.Now().UTC()
	nowUnixNano := now.UnixNano()
	nowStr := now.Format(model.TimestampFormat)

	article := model.Article{
		Title:       requestBody.Article.Title,
		Description: requestBody.Article.Description,
		Body:        requestBody.Article.Body,
		TagList:     requestBody.Article.TagList,
		CreatedAt:   nowUnixNano,
		UpdatedAt:   nowUnixNano,
		Author:      user.Username,
	}

	err = service.PutArticle(&article)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	responseBody := ResponseBody{
		Article: ArticleResponse{
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        requestBody.Article.TagList,
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

	return util.NewSuccessResponse(201, responseBody)
}

func main() {
	lambda.Start(Handle)
}
