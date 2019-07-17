package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/service"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/util"
	"strconv"
)

func Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, _, err := service.GetCurrentUser(request.Headers["Authorization"])
	if err != nil {
		return util.NewUnauthorizedResponse()
	}

	commentId, err := strconv.ParseInt(request.PathParameters["id"], 10, 64)
	if err != nil {
		return util.NewErrorResponse(util.NewInputError("id", "invalid"))
	}

	err = service.DeleteComment(request.PathParameters["slug"], commentId, user.Username)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handle)
}
