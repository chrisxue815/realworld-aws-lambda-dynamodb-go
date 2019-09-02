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

func PassArticleByValue(article Article, goPanic bool) {
	if goPanic {
		if article.ArticleId == 0 {
			// noinline
			// https://github.com/golang/go/wiki/CompilerOptimizations#function-inlining
			panic(nil)
		}
	}
}

func PassArticleByPointer(article *Article, goPanic bool) {
	if goPanic {
		if article.ArticleId == 0 {
			// noinline
			// https://github.com/golang/go/wiki/CompilerOptimizations#function-inlining
			panic(nil)
		}
	}
}

func BenchmarkPassArticleByValue(b *testing.B) {
	article := Article{}
	for i := 0; i < b.N; i++ {
		PassArticleByValue(article, false)
	}
}

func BenchmarkPassArticleByPointer(b *testing.B) {
	article := Article{}
	for i := 0; i < b.N; i++ {
		PassArticleByPointer(&article, false)
	}
}
