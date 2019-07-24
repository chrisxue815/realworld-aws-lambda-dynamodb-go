package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/service"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/util"
)

type ResponseBody struct {
	Profile ProfileResponse `json:"profile"`
}

type ProfileResponse struct {
	Username  string `json:"username"`
	Image     string `json:"image"`
	Bio       string `json:"bio"`
	Following bool   `json:"following"`
}

func Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user, _, _ := service.GetCurrentUser(request.Headers["Authorization"])

	publisher, err := service.GetUserByUsername(request.PathParameters["username"])
	if err != nil {
		return util.NewErrorResponse(err)
	}

	following, err := service.IsFollowing(user, []string{publisher.Username})
	if err != nil {
		return util.NewErrorResponse(err)
	}

	responseBody := ResponseBody{
		Profile: ProfileResponse{
			Username:  publisher.Username,
			Image:     publisher.Image,
			Bio:       publisher.Bio,
			Following: following[0],
		},
	}

	return util.NewSuccessResponse(200, responseBody)
}

func main() {
	lambda.Start(Handle)
}
