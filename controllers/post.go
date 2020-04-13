package controllers

import (
	"fmt"
	"log"
	"net/http"
	"wiki/context"
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

type NewPostForm struct {
	Title   string `schema:"title"`
	Content string `schema:"content"`
}

// Create is used to process the signup form when a user
// submits it. This is used to create a new user account.
//
// POST /posts
func (p *Posts) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form NewPostForm
	if err := ParseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		p.New.Render(w, vd)
		return
	}

	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	post := &posts.Post{
		UserID:  user.ID,
		Title:   form.Title,
		Content: form.Content,
	}

	err := p.ps.Create(post)
	if err != nil {
		vd.SetAlert(err)
		p.New.Render(w, vd)
		return
	}

	fmt.Fprintln(w, post)
}
