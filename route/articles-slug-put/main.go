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

	oldArticle, err := service.GetArticleBySlug(input.PathParameters["slug"])
	if err != nil {
		return util.NewErrorResponse(err)
	}

	newArticle := createNewArticle(request, oldArticle)

	err = service.UpdateArticle(oldArticle, &newArticle)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	isFavorited, authors, following, err := service.GetArticleRelatedProperties(user, []model.Article{newArticle}, true)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	response := Response{
		Article: ArticleResponse{
			Slug:           newArticle.Slug,
			Title:          newArticle.Title,
			Description:    newArticle.Description,
			Body:           newArticle.Body,
			TagList:        newArticle.TagList,
			CreatedAt:      time.Unix(0, newArticle.CreatedAt).Format(model.TimestampFormat),
			UpdatedAt:      time.Unix(0, newArticle.UpdatedAt).Format(model.TimestampFormat),
			Favorited:      isFavorited[0],
			FavoritesCount: newArticle.FavoritesCount,
			Author: AuthorResponse{
				Username:  authors[0].Username,
				Bio:       authors[0].Bio,
				Image:     authors[0].Image,
				Following: following[0],
			},
		},
	}

	return util.NewSuccessResponse(200, response)
}

func createNewArticle(request Request, oldArticle model.Article) model.Article {
	newArticle := model.Article{
		ArticleId:      oldArticle.ArticleId,
		Title:          request.Article.Title,
		Description:    request.Article.Description,
		Body:           request.Article.Body,
		TagList:        request.Article.TagList,
		CreatedAt:      oldArticle.CreatedAt,
		UpdatedAt:      time.Now().UTC().UnixNano(),
		FavoritesCount: oldArticle.FavoritesCount,
		Author:         oldArticle.Author,
	}

	if newArticle.Title == "" {
		newArticle.Title = oldArticle.Title
	}

	if newArticle.Description == "" {
		newArticle.Description = oldArticle.Description
	}

	if newArticle.Body == "" {
		newArticle.Body = oldArticle.Body
	}

	if newArticle.TagList == nil {
		newArticle.TagList = oldArticle.TagList
	}

	return newArticle
}

func main() {
	lambda.Start(Handle)
}
