package model

import (
	"github.com/chrisxue815/realworld-aws-lambda-dynamodb-go/util"
)

const MaxCommentId = 0x1000000 // exclusive

type Comment struct {
	ArticleId int64
	CommentId int64
	CreatedAt int64
	UpdatedAt int64
	Body      string
	Author    string
}

func (comment *Comment) Validate() error {
	if comment.Body == "" {
		return util.NewInputError("body", "can't be blank")
	}

	return nil
}
