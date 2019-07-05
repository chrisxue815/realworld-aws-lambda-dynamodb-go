package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/service"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/util"
)

func Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, _, err := service.GetCurrentUser(request.Headers["Authorization"])
	if err != nil {
		return util.NewUnauthorizedResponse()
	}

	err = service.DeleteArticle(request.PathParameters["slug"], user.Username)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handle)
}
