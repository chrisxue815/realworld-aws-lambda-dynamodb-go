package service

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"math/rand"
	"os"
	"sync"
	"time"
)

var once sync.Once
var svc *dynamodb.DynamoDB

var UserTableName = NewTableName("user")
var EmailUserTableName = NewTableName("email-user")
var FollowTableName = NewTableName("follow")
var ArticleTableName = NewTableName("article")
var ArticleTagTableName = NewTableName("article-tag")
var TagTableName = NewTableName("tag")
var FavoriteArticleTableName = NewTableName("favorite-article")
var CommentTableName = NewTableName("comment")

var ArticleIdRand = NewRand()
var CommentIdRand = NewRand()

func initializeSingletons() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc = dynamodb.New(sess)
}

func DynamoDB() *dynamodb.DynamoDB {
	once.Do(initializeSingletons)
	return svc
}

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

type Rand struct {
	random *rand.Rand
}

func NewRand() Rand {
	r := Rand{}
	r.RenewSeed()
	return r
}

func (r *Rand) RenewSeed() {
	r.random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (r *Rand) Get() *rand.Rand {
	return r.random
}
