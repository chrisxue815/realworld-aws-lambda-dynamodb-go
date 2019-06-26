package service

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsUpdateBuilderEmpty(t *testing.T) {
	assert.True(t, IsUpdateBuilderEmpty(expression.UpdateBuilder{}))
	assert.False(t, IsUpdateBuilderEmpty(expression.Set(expression.Name("Name"), expression.Value("Value"))))
}
