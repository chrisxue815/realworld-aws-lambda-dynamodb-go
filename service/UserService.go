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

func UpdateUser(oldUser model.User, newUser model.User) error {
	emailUser := model.EmailUser{
		Email:    newUser.Email,
		Username: newUser.Username,
	}

	emailUserItem, err := dynamodbattribute.MarshalMap(emailUser)
	if err != nil {
		return err
	}

	transaction := dynamodb.TransactWriteItemsInput{}

	if oldUser.Email != newUser.Email {
		transaction.TransactItems = []*dynamodb.TransactWriteItem{
			{
				Put: &dynamodb.Put{
					TableName:           aws.String(model.EmailUserTableName()),
					Item:                emailUserItem,
					ConditionExpression: aws.String("attribute_not_exists(Email)"),
				},
			},
			{
				Delete: &dynamodb.Delete{
					TableName: aws.String(model.EmailUserTableName()),
					Key: map[string]*dynamodb.AttributeValue{
						"Email": {
							S: aws.String(oldUser.Email),
						},
					},
					ConditionExpression: aws.String("attribute_exists(Email)"),
				},
			},
		}
	}

	transaction.TransactItems = append(transaction.TransactItems, &dynamodb.TransactWriteItem{
		Update: &dynamodb.Update{
			TableName: aws.String(model.UserTableName()),
			Key: map[string]*dynamodb.AttributeValue{
				"Username": {
					S: aws.String(oldUser.Username),
				},
			},
			ConditionExpression: aws.String("attribute_exists(Username)"),
			ExpressionAttributeNames: map[string]*string{
				"#EMAIL":    aws.String("Email"),
				"#IMAGE":    aws.String("Image"),
				"#BIO":      aws.String("Bio"),
				"#PASSWORD": aws.String("Password"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":email": {
					S: aws.String(newUser.Email),
				},
				":image": {
					S: aws.String(newUser.Image),
				},
				":bio": {
					S: aws.String(newUser.Bio),
				},
				":password": {
					B: newUser.Password,
				},
			},
			UpdateExpression: aws.String("SET #EMAIL=:email, #IMAGE=:image, #BIO=:bio, #PASSWORD=:password"),
		},
	})

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
