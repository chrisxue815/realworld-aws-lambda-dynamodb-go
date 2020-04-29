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
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
	request := Request{}
	err := json.Unmarshal([]byte(input.Body), &request)
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

	user := model.User{
		Username:     request.User.Username,
		Email:        request.User.Email,
		PasswordHash: passwordHash,
	}

	err = service.PutUser(user)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	token, err := model.GenerateToken(user.Username)
	if err != nil {
		return util.NewErrorResponse(err)
	}

	response := Response{
		User: UserResponse{
			Username: user.Username,
			Email:    user.Email,
			Image:    user.Image,
			Bio:      user.Bio,
			Token:    token,
		},
	}

	return util.NewSuccessResponse(201, response)
}

func main() {
	lambda.Start(Handle)
}
