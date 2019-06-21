package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
)

func CreateUser(user model.User) error {
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

	_, err = DynamoDB().TransactWriteItems(&transaction)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(email string) (model.User, error) {
	username, err := GetUsernameByEmail(email)
	if err != nil {
		return model.User{}, err
	}

	return GetUserByUsername(username)
}

func GetUsernameByEmail(email string) (string, error) {
	emailUser := model.EmailUser{}
	err := GetItemByKey(model.EmailUserTableName(), "Email", email, &emailUser)
	if err != nil {
		return "", err
	}

	return emailUser.Username, nil
}

func GetUserByUsername(username string) (model.User, error) {
	user := model.User{}
	err := GetItemByKey(model.UserTableName(), "Username", username, &user)
	return user, err
}

func GetCurrentUser(auth string) (model.User, string, error) {
	username, token, err := VerifyAuthorization(auth)
	if err != nil {
		return model.User{}, "", err
	}

	user, err := GetUserByUsername(username)
	if err != nil {
		return model.User{}, "", err
	}

	return user, token, err
}
