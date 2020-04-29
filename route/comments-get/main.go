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
	Comments []CommentResponse `json:"comments"`
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

func Handle(input events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, _, _ := service.GetCurrentUser(input.Headers["Authorization"])

	comments, err := service.GetComments(input.PathParameters["slug"])
	if err != nil {
		return util.NewErrorResponse(err)
	}

	authors, following, err := service.GetCommentRelatedProperties(user, comments)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	commentResponses := make([]CommentResponse, 0, len(comments))

	for i, comment := range comments {
		commentResponses = append(commentResponses, CommentResponse{
			Id:        comment.CommentId,
			Body:      comment.Body,
			CreatedAt: time.Unix(0, comment.CreatedAt).Format(model.TimestampFormat),
			UpdatedAt: time.Unix(0, comment.UpdatedAt).Format(model.TimestampFormat),
			Author: AuthorResponse{
				Username:  authors[i].Username,
				Bio:       authors[i].Bio,
				Image:     authors[i].Image,
				Following: following[i],
			},
		})
	}

	response := Response{
		Comments: commentResponses,
	}

	return util.NewSuccessResponse(200, response)
}

func main() {
	lambda.Start(Handle)
}
