package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlugToArticleId(t *testing.T) {
	testCases := []struct {
		slug          string
		expected      int64
		expectedError bool
	}{
		{"how-to-train-your-dragon-74728a", 0x74728a, false},
		{"74728a", 0x74728a, false},
	}

	for _, testCase := range testCases {
		actual, err := SlugToArticleId("how-to-train-your-dragon-74728a")
		assert.Equal(t, testCase.expected, actual, "%+v", testCase)
		assert.Equal(t, testCase.expectedError, err != nil, "%+v", testCase)
	}
}
