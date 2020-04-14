package controllers

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"wiki/context"
	"wiki/models"
	"wiki/models/posts"
	"wiki/views"
)

func NewPosts(ps posts.PostService, r *mux.Router) *Posts {
	return &Posts{
		New:           views.NewView("bulma", "posts/new"),
		ShowView:      views.NewView("bulma", "posts/show"),
		EditView:      views.NewView("bulma", "posts/edit"),
		PostIndexView: views.NewView("bulma", "posts/index"),
		ps:            ps,
		r:             r,
	}
}

type Posts struct {
	New           *views.View
	ShowView      *views.View
	EditView      *views.View
	PostIndexView *views.View
	ps            posts.PostService
	r             *mux.Router
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

	url, err := p.r.Get("show_post").URL("id", strconv.Itoa(int(post.ID)))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// /GET /posts
func (p *Posts) PostIndex(w http.ResponseWriter, r *http.Request) {
	p.PostIndexView.Render(w, nil)
}

func (p *Posts) postByID(w http.ResponseWriter, r *http.Request) (*posts.Post, error) {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}
	post, err := p.ps.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Post not found in DB.", http.StatusNotFound)
		default:
			http.Error(w, "Whoops! Something went wrong.", http.StatusInternalServerError)
		}
		return nil, err
	}
	return post, nil
}

// GET /posts/:id
func (p *Posts) Show(w http.ResponseWriter, r *http.Request) {
	post, err := p.postByID(w, r)
	if err != nil {
		// the postByID method already handled the rendering of the error
		return
	}
	md := []byte(post.Content)
	html := markdown.ToHTML(md, nil, nil)
	post.ContentHTML = template.HTML(html)
	vd := views.Data{Yield: post}
	p.ShowView.Render(w, vd)
}

// Edit
//
// /GET /post/:id/edit
func (p *Posts) Edit(w http.ResponseWriter, r *http.Request) {
	post, err := p.postByID(w, r)
	if err != nil {
		// the postByID method already handled the rendering of the error
		return
	}
	user := context.User(r.Context())
	if post.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	vd := views.Data{Yield: post}
	p.EditView.Render(w, vd)
}

// Update
//
// /POST /post/:id/update
func (p *Posts) Update(w http.ResponseWriter, r *http.Request) {
	post, err := p.postByID(w, r)
	if err != nil {
		// the postByID method already handled the rendering of the error
		return
	}
	user := context.User(r.Context())
	if post.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	var form NewPostForm
	vd := views.Data{Yield: post}
	if err := ParseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		p.EditView.Render(w, vd)
		return
	}
	post.Title = form.Title
	post.Content = form.Content
	err = p.ps.Update(post)
	if err != nil {
		log.Println(err)
		vd.SetAlert(err)
		p.EditView.Render(w, vd)
		return
	}

	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return
	}
	urlPath := fmt.Sprintf("/post/%d", id)
	http.Redirect(w, r, urlPath, http.StatusFound)
}

// Delete
//
// /POST /post/:id/delete
func (p *Posts) Delete(w http.ResponseWriter, r *http.Request) {
	post, err := p.postByID(w, r)
	if err != nil {
		// the postByID method already handled the rendering of the error
		return
	}
	user := context.User(r.Context())
	if post.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	vd := views.Data{Yield: post}
	err = p.ps.Delete(post.ID)
	if err != nil {
		log.Println(err)
		vd.SetAlert(err)
		p.EditView.Render(w, vd)
		return
	}
	vd.AlertSuccess("Post deleted!")
	p.PostIndexView.Render(w, vd)
}
