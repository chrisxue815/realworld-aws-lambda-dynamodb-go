package util

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

type ErrorResponseBody struct {
	Errors ErrorList `json:"errors"`
}

type ErrorList struct {
	Body []string `json:"body"`
}

func NewErrorResponse(err error) (events.APIGatewayProxyResponse, error) {
	js, err := NewErrorResponseJSON(err.Error())
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 422,
		Body:       js,
	}
	return response, nil
}

func NewErrorResponseJSON(msg string) (string, error) {
	body := NewErrorResponseBody(msg)
	js, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return string(js), nil
}

func NewErrorResponseBody(msg string) ErrorResponseBody {
	return ErrorResponseBody{
		Errors: ErrorList{
			Body: []string{
				msg,
			},
		},
	}
}
