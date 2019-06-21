package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/service"
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
		return model.NewErrorResponse(err)
	}

	password, err := service.Scrypt(requestBody.User.Password)
	if err != nil {
		return model.NewErrorResponse(err)
	}

	user := model.User{
		Username: requestBody.User.Username,
		Email:    requestBody.User.Email,
		Password: password,
	}

	err = createUser(user)
	if err != nil {
		return model.NewErrorResponse(err)
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

func createUser(user model.User) error {
	userItem, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	emailUser := model.EmailUser{
		Email:    user.Email,
		Username: user.Username,
	}

	emailUserItem, err := dynamodbattribute.MarshalMap(emailUser)
	if err != nil {
		return err
	}

	transaction := dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Put: &dynamodb.Put{
					TableName:           aws.String(model.UserTableName()),
					Item:                userItem,
					ConditionExpression: aws.String("attribute_not_exists(Username)"),
				},
			},
			{
				Put: &dynamodb.Put{
					TableName:           aws.String(model.EmailUserTableName()),
					Item:                emailUserItem,
					ConditionExpression: aws.String("attribute_not_exists(Email)"),
				},
			},
		},
	}

	_, err = service.DynamoDB().TransactWriteItems(&transaction)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(Handle)
}
