package util

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

func NewSuccessResponse(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return NewErrorResponse(err)
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(jsonBody),
	}

	return response, nil
}
