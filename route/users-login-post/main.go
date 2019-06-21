package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/service"
)

type RequestBody struct {
	User UserRequest `json:"user"`
}

type UserRequest struct {
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
		return model.NewErrorResponse(err)
	}

	user, err := getUserByEmail(requestBody.User.Email)
	if err != nil {
		return model.NewErrorResponse(err)
	}

	password, err := service.Scrypt(requestBody.User.Password)
	if err != nil {
		return model.NewErrorResponse(err)
	}

	if !bytes.Equal(password, user.Password) {
		return model.NewErrorResponse(errors.New("wrong password"))
	}

	token, err := service.GenerateJWT(user.Username)
	if err != nil {
		return model.NewErrorResponse(err)
	}

	responseBody := ResponseBody{
		Username: user.Username,
		Email:    user.Email,
		Image:    user.Image,
		Bio:      user.Bio,
		Token:    token,
	}

	responseJSON, err := json.Marshal(responseBody)
	if err != nil {
		return model.NewErrorResponse(err)
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseJSON),
	}

	return response, nil
}

func getUserByEmail(email string) (model.User, error) {
	username, err := getUsernameByEmail(email)
	if err != nil {
		return model.User{}, err
	}

	return getUserByUsername(username)
}

func getUsernameByEmail(email string) (string, error) {
	emailUser := model.EmailUser{}
	err := service.GetItemByKey(model.EmailUserTableName(), "Email", email, &emailUser)
	if err != nil {
		return "", err
	}

	return emailUser.Username, nil
}

func getUserByUsername(username string) (model.User, error) {
	user := model.User{}
	err := service.GetItemByKey(model.UserTableName(), "Username", username, &user)
	return user, err
}

func main() {
	lambda.Start(Handle)
}
