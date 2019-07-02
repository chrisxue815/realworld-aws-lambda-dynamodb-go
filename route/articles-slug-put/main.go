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

	oldArticle, err := service.GetArticleBySlug(request.PathParameters["slug"])
	if err != nil {
		return util.NewErrorResponse(err)
	}

	if oldArticle.ArticleId == 0 {
		return util.NewErrorResponse(util.NewInputError("slug", "not found"))
	}

	newArticle := model.Article{
		ArticleId:      oldArticle.ArticleId,
		Title:          requestBody.Article.Title,
		Description:    requestBody.Article.Description,
		Body:           requestBody.Article.Body,
		TagList:        requestBody.Article.TagList,
		CreatedAt:      oldArticle.CreatedAt,
		UpdatedAt:      time.Now().UTC().UnixNano(),
		FavoritesCount: oldArticle.FavoritesCount,
		Author:         oldArticle.Author,
	}

	err = service.UpdateArticle(oldArticle, &newArticle)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	isFavorited, authors, following, err := service.GetArticleRelatedProperties(user, []model.Article{newArticle})
	if err != nil {
		return util.NewErrorResponse(err)
	}

	responseBody := ResponseBody{
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

	return util.NewSuccessResponse(200, responseBody)
}

func main() {
	lambda.Start(Handle)
}
