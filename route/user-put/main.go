package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/service"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/util"
)

type Request struct {
	User UserRequest `json:"user"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Image    string `json:"image"`
	Bio      string `json:"bio"`
}

type Response struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Image    string `json:"image"`
	Bio      string `json:"bio"`
	Token    string `json:"token"`
}

func Handle(input events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	oldUser, token, err := service.GetCurrentUser(input.Headers["Authorization"])
	if err != nil {
		return util.NewUnauthorizedResponse()
	}

	request := Request{}
	err = json.Unmarshal([]byte(input.Body), &request)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	err = model.ValidatePassword(request.User.Password)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	passwordHash, err := model.Scrypt(request.User.Password)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	newUser := model.User{
		Username:     oldUser.Username,
		Email:        request.User.Email,
		PasswordHash: passwordHash,
		Image:        request.User.Image,
		Bio:          request.User.Bio,
	}

	err = service.UpdateUser(*oldUser, newUser)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	response := Response{
		User: UserResponse{
			Username: newUser.Username,
			Email:    newUser.Email,
			Image:    newUser.Image,
			Bio:      newUser.Bio,
			Token:    token,
		},
	}

	return util.NewSuccessResponse(200, response)
}

func main() {
	lambda.Start(Handle)
}
