package util

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

func NewSuccessResponse(statusCode int, v interface{}) (events.APIGatewayProxyResponse, error) {
	responseJSON, err := json.Marshal(v)
	if err != nil {
		return NewErrorResponse(err)
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(responseJSON),
	}

	return response, nil
}
