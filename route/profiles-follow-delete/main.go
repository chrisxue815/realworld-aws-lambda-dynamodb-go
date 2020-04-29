package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/service"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/util"
)

type Response struct {
	Profile ProfileResponse `json:"profile"`
}

type ProfileResponse struct {
	Username  string `json:"username"`
	Image     string `json:"image"`
	Bio       string `json:"bio"`
	Following bool   `json:"following"`
}

func Handle(input events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, _, err := service.GetCurrentUser(input.Headers["Authorization"])
	if err != nil {
		return util.NewUnauthorizedResponse()
	}

	publisher, err := service.GetUserByUsername(input.PathParameters["username"])
	if err != nil {
		return util.NewErrorResponse(err)
	}

	err = service.Unfollow(user.Username, publisher.Username)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	response := Response{
		Profile: ProfileResponse{
			Username:  publisher.Username,
			Image:     publisher.Image,
			Bio:       publisher.Bio,
			Following: false,
		},
	}

	return util.NewSuccessResponse(200, response)
}

func main() {
	lambda.Start(Handle)
}
