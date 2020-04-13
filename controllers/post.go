package controllers

import (
	"wiki/models/posts"
	"wiki/views"
)

func NewPosts(ps posts.PostService) *Posts {
	return &Posts{
		New: views.NewView("bulma", "posts/new"),
		ps:  ps,
	}
}

type Posts struct {
	New *views.View
	ps  posts.PostService
}
