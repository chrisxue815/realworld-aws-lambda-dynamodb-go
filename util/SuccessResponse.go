package util

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

func NewSuccessResponse(v interface{}) (events.APIGatewayProxyResponse, error) {
	responseJSON, err := json.Marshal(v)
	if err != nil {
		return NewErrorResponse(err)
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseJSON),
	}

	return response, nil
}
