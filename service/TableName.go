package service

import (
	"fmt"
	"os"
)

var Stage = os.Getenv("STAGE")

var UserTableName = makeTableName("user")
var EmailUserTableName = makeTableName("email-user")
var FollowTableName = makeTableName("follow")
var ArticleTableName = makeTableName("article")
var ArticleTagTableName = makeTableName("article-tag")
var TagTableName = makeTableName("tag")
var FavoriteArticleTableName = makeTableName("favorite-article")
var CommentTableName = makeTableName("comment")

func makeTableName(suffix string) string {
	return fmt.Sprintf("realworld-%s-%s", Stage, suffix)
}
