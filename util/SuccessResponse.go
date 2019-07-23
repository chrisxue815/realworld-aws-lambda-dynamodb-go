package util

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

func CORSHeaders() map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Credentials": "true",
	}
}

func NewSuccessResponse(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    CORSHeaders(),
	}

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return NewErrorResponse(err)
		}

		response.Body = string(jsonBody)
	}

	return response, nil
}
