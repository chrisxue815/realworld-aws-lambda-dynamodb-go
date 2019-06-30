package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/service"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/util"
)

type RequestBody struct {
	User UserRequest `json:"user"`
}

type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResponseBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Image    string `json:"image"`
	Bio      string `json:"bio"`
	Token    string `json:"token"`
}

func Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	requestBody := RequestBody{}
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	err = model.ValidatePassword(requestBody.User.Password)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	passwordHash, err := service.Scrypt(requestBody.User.Password)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	user := model.User{
		Username:     requestBody.User.Username,
		Email:        requestBody.User.Email,
		PasswordHash: passwordHash,
	}

	err = service.PutUser(user)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	token, err := service.GenerateToken(user.Username)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	responseBody := ResponseBody{
		Username: user.Username,
		Email:    user.Email,
		Image:    user.Image,
		Bio:      user.Bio,
		Token:    token,
	}

	return util.NewSuccessResponse(201, responseBody)
}

func main() {
	lambda.Start(Handle)
}
