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
	Comment CommentRequest `json:"comment"`
}

type CommentRequest struct {
	Body string `json:"body"`
}

type ResponseBody struct {
	Comment CommentResponse `json:"comment"`
}

type CommentResponse struct {
	Id        int64          `json:"id"`
	CreatedAt string         `json:"createdAt"`
	UpdatedAt string         `json:"updatedAt"`
	Body      string         `json:"body"`
	Author    AuthorResponse `json:"author"`
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

	// Make sure article exists, at least at this point
	article, err := service.GetArticleBySlug(request.PathParameters["slug"])
	if err != nil {
		return util.NewErrorResponse(err)
	}

	now := time.Now().UTC()
	nowUnixNano := now.UnixNano()
	nowStr := now.Format(model.TimestampFormat)

	comment := model.Comment{
		CommentKey: model.CommentKey{
			ArticleId: article.ArticleId,
		},
		CreatedAt: nowUnixNano,
		UpdatedAt: nowUnixNano,
		Body:      requestBody.Comment.Body,
		Author:    user.Username,
	}

	err = service.PutComment(&comment)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	responseBody := ResponseBody{
		Comment: CommentResponse{
			Id:        comment.CommentId,
			Body:      comment.Body,
			CreatedAt: nowStr,
			UpdatedAt: nowStr,
			Author: AuthorResponse{
				Username:  user.Username,
				Bio:       user.Bio,
				Image:     user.Image,
				Following: false,
			},
		},
	}

	return util.NewSuccessResponse(200, responseBody)
}

func main() {
	lambda.Start(Handle)
}
