package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/service"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/util"
	"time"
)

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

	articleId, err := model.SlugToArticleId(input.PathParameters["slug"])
	if err != nil {
		return util.NewErrorResponse(err)
	}

	favoriteArticleKey := model.FavoriteArticleKey{
		Username:  user.Username,
		ArticleId: articleId,
	}

	err = service.UnfavoriteArticle(favoriteArticleKey)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	article, err := service.GetArticleByArticleId(articleId)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	isFavorited, authors, following, err := service.GetArticleRelatedProperties(user, []model.Article{article}, true)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	response := Response{
		Article: ArticleResponse{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        article.TagList,
			CreatedAt:      time.Unix(0, article.CreatedAt).Format(model.TimestampFormat),
			UpdatedAt:      time.Unix(0, article.UpdatedAt).Format(model.TimestampFormat),
			Favorited:      isFavorited[0],
			FavoritesCount: article.FavoritesCount,
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

func main() {
	lambda.Start(Handle)
}
