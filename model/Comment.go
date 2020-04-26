package model

const MaxCommentId = 0x1000000 // exclusive

type CommentKey struct {
	ArticleId int64
	CommentId int64
}

type Comment struct {
	CommentKey
	CreatedAt int64
	UpdatedAt int64
	Body      string
	Author    string
}

func (comment *Comment) Validate() error {
	if comment.Body == "" {
		return NewInputError("body", "can't be blank")
	}

	return nil
}
