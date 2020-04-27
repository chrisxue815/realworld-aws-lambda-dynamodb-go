package service

import (
	"fmt"
	"os"
	"sync"
)

var UserTableName = NewTableName("user")
var EmailUserTableName = NewTableName("email-user")
var FollowTableName = NewTableName("follow")
var ArticleTableName = NewTableName("article")
var ArticleTagTableName = NewTableName("article-tag")
var TagTableName = NewTableName("tag")
var FavoriteArticleTableName = NewTableName("favorite-article")
var CommentTableName = NewTableName("comment")

type TableName struct {
	suffix   string
	fullName string
	once     sync.Once
}

func NewTableName(suffix string) TableName {
	return TableName{
		suffix: suffix,
	}
}

func (t *TableName) Get() string {
	t.once.Do(func() {
		t.fullName = fmt.Sprintf("realworld-%s-%s", os.Getenv("STAGE"), t.suffix)
	})
	return t.fullName
}
