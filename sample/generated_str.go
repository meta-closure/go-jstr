package model

import (
	"github.com/lestrrat/go-jsval"
)

type User struct {
	ScreenName jsval.MaybeString `json:"screen_name,omitempty"`
	UserId     jsval.MaybeString `json:"user_id,omitempty"`
}
type Post struct {
	Title  jsval.MaybeString `json:"title,omitempty"`
	Body   jsval.MaybeString `json:"body,omitempty"`
	PostId jsval.MaybeInt    `json:"post_id,omitempty"`
}
type Posts []Post
