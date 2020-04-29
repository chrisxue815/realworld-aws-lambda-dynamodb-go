package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/service"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/util"
	"strconv"
)

func Handle(input events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, _, err := service.GetCurrentUser(input.Headers["Authorization"])
	if err != nil {
		return util.NewUnauthorizedResponse()
	}

	commentId, err := strconv.ParseInt(input.PathParameters["id"], 10, 64)
	if err != nil {
		return util.NewErrorResponse(model.NewInputError("id", "invalid"))
	}

	err = service.DeleteComment(input.PathParameters["slug"], commentId, user.Username)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	return util.NewSuccessResponse(200, nil)
}

func main() {
	lambda.Start(Handle)
}
