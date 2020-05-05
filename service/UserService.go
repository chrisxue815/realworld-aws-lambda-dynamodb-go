package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/model"
)

func PutUser(user model.User) error {
	err := user.Validate()
	if err != nil {
		return err
	}

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

	// Put a new user, make sure username and email are unique
	transaction := dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Put: &dynamodb.Put{
					TableName:           aws.String(UserTableName),
					Item:                userItem,
					ConditionExpression: aws.String("attribute_not_exists(Username)"),
				},
			},
			{
				Put: &dynamodb.Put{
					TableName:           aws.String(EmailUserTableName),
					Item:                emailUserItem,
					ConditionExpression: aws.String("attribute_not_exists(Email)"),
				},
			},
		},
	}

	_, err = DynamoDB().TransactWriteItems(&transaction)
	if err != nil {
		// TODO: distinguish:
		// NewInputError("username", "has already been taken")
		// NewInputError("email", "has already been taken")
		return err
	}

	return nil
}

func UpdateUser(oldUser model.User, newUser model.User) error {
	err := newUser.Validate()
	if err != nil {
		return err
	}

	transactItems := make([]*dynamodb.TransactWriteItem, 0, 3)

	if oldUser.Email != newUser.Email {
		newEmailUser := model.EmailUser{
			Email:    newUser.Email,
			Username: newUser.Username,
		}

		newEmailUserItem, err := dynamodbattribute.MarshalMap(newEmailUser)
		if err != nil {
			return err
		}

		// Link user with the new email
		transactItems = append(transactItems, &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				TableName:           aws.String(EmailUserTableName),
				Item:                newEmailUserItem,
				ConditionExpression: aws.String("attribute_not_exists(Email)"),
			},
		})

		// Unlink user from the old email
		transactItems = append(transactItems, &dynamodb.TransactWriteItem{
			Delete: &dynamodb.Delete{
				TableName:           aws.String(EmailUserTableName),
				Key:                 StringKey("Email", oldUser.Email),
				ConditionExpression: aws.String("attribute_exists(Email)"),
			},
		})
	}

	newUserItem, err := dynamodbattribute.MarshalMap(newUser)
	if err != nil {
		return err
	}

	// Update user info
	transactItems = append(transactItems, &dynamodb.TransactWriteItem{
		Put: &dynamodb.Put{
			TableName:                 aws.String(UserTableName),
			Item:                      newUserItem,
			ConditionExpression:       aws.String("Email = :email"),
			ExpressionAttributeValues: StringKey(":email", oldUser.Email),
		},
	})

	_, err = DynamoDB().TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})
	if err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(email string) (model.User, error) {
	if email == "" {
		return model.User{}, model.NewInputError("email", "can't be blank")
	}

	username, err := GetUsernameByEmail(email)
	if err != nil {
		return model.User{}, err
	}

	return GetUserByUsername(username)
}

func GetUsernameByEmail(email string) (string, error) {
	emailUser := model.EmailUser{}
	found, err := GetItemByKey(EmailUserTableName, StringKey("Email", email), &emailUser)

	if err != nil {
		return "", err
	}

	if !found {
		return "", model.NewInputError("email", "not found")
	}

	return emailUser.Username, nil
}

func GetUserByUsername(username string) (model.User, error) {
	if username == "" {
		return model.User{}, model.NewInputError("username", "can't be blank")
	}

	user := model.User{}
	found, err := GetItemByKey(UserTableName, StringKey("Username", username), &user)

	if err != nil {
		return model.User{}, err
	}

	if !found {
		return model.User{}, model.NewInputError("username", "not found")
	}

	return user, err
}

func GetCurrentUser(auth string) (*model.User, string, error) {
	username, token, err := model.VerifyAuthorization(auth)
	if err != nil {
		return nil, "", err
	}

	user, err := GetUserByUsername(username)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func GetUserListByUsername(usernames []string) ([]model.User, error) {
	if len(usernames) == 0 {
		return make([]model.User, 0), nil
	}

	usernameSet := make(map[string]bool)
	for _, username := range usernames {
		usernameSet[username] = true
	}

	keys := make([]AWSObject, 0, len(usernameSet))
	for username := range usernameSet {
		keys = append(keys, StringKey("Username", username))
	}

	batchGetUsers := dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			UserTableName: {
				Keys: keys,
			},
		},
	}

	responses, err := BatchGetItems(&batchGetUsers, len(usernames))
	if err != nil {
		return nil, err
	}

	usersByUsername := make(map[string]model.User)

	for _, response := range responses {
		for _, items := range response {
			for _, item := range items {
				user := model.User{}
				err = dynamodbattribute.UnmarshalMap(item, &user)
				if err != nil {
					return nil, err
				}

				usersByUsername[user.Username] = user
			}
		}
	}

	users := make([]model.User, 0, len(usernames))
	for _, username := range usernames {
		users = append(users, usersByUsername[username])
	}

	return users, nil
}
