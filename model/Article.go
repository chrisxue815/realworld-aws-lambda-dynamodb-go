package model

import (
	"fmt"
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/util"
	"github.com/gosimple/slug"
	"strconv"
)

const TimestampFormat = "2006-01-02T15:04:05.000Z"
const MaxArticleId = 0x1000000 // exclusive
const MaxNumTagsPerArticle = 5

type Article struct {
	ArticleId      int64
	Slug           string
	Title          string
	Description    string
	Body           string
	TagList        []string
	CreatedAt      int64
	UpdatedAt      int64
	FavoritesCount int64
	Author         string
	Dummy          byte // Always 0, used for sorting articles by CreatedAt
}

type ArticleTag struct {
	Tag       string
	ArticleId int64
	CreatedAt int64
}

type FavoriteArticle struct {
	Username   string
	ArticleId  int64
	FavoriteAt int64
}

func (article *Article) Validate() error {
	if article.Title == "" {
		return util.NewInputError("title", "can't be blank")
	}

	if article.Description == "" {
		return util.NewInputError("description", "can't be blank")
	}

	if article.Body == "" {
		return util.NewInputError("body", "can't be blank")
	}

	if article.TagList == nil {
		article.TagList = make([]string, 0)
	} else if len(article.TagList) > MaxNumTagsPerArticle {
		return util.NewInputError("tagList", fmt.Sprintf("cannot add more than %d tags per article", MaxNumTagsPerArticle))
	}

	return nil
}

func (article *Article) MakeSlug() {
	slugPrefix := slug.Make(article.Title)
	article.Slug = slugPrefix + "-" + strconv.FormatInt(article.ArticleId, 16)
}
