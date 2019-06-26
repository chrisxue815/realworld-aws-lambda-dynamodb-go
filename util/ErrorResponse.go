package util

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

type InputError map[string][]string

func (e InputError) Error() string {
	js, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}

	return string(js)
}

func NewInputError(inputName, message string) InputError {
	return InputError{
		inputName: {message},
	}
}

type ErrorResponseBody struct {
	Errors InputError `json:"errors"`
}

func NewErrorResponse(err error) (events.APIGatewayProxyResponse, error) {
	inputError, ok := err.(InputError)
	if !ok {
		// Internal server error
		return events.APIGatewayProxyResponse{}, err
	}

	body := ErrorResponseBody{
		Errors: inputError,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 422,
		Body:       string(jsonBody),
	}
	return response, nil
}

func NewUnauthorizedResponse() (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		StatusCode: 401,
	}
	return response, nil
}
