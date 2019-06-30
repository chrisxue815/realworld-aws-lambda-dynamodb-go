package model

import "github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/util"

const TimestampFormat = "2006-01-02T15:04:05.000Z"
const MaxArticleId = 0x1000000 // exclusive

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

func (a *Article) Validate() error {
	if a.Title == "" {
		return util.NewInputError("title", "must not be empty")
	}

	if a.TagList == nil {
		a.TagList = make([]string, 0)
	} else if len(a.TagList) > 5 {
		return util.NewInputError("tagList", "cannot add more than 5 tags per article")
	}

	return nil
}
