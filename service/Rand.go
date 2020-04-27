package service

import (
	"math/rand"
	"time"
)

var ArticleIdRand = NewRand()
var CommentIdRand = NewRand()

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
