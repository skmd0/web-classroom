package controllers

import (
	"github.com/gorilla/mux"
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
		New:      views.NewView("bulma", "posts/new"),
		ShowView: views.NewView("bulma", "posts/show"),
		ps:       ps,
		r:        r,
	}
}

type Posts struct {
	New      *views.View
	ShowView *views.View
	ps       posts.PostService
	r        *mux.Router
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

// GET /posts/:id
func (p *Posts) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
	}
	post, err := p.ps.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Post not found in DB.", http.StatusNotFound)
		default:
			http.Error(w, "Whoops! Something went wrong.", http.StatusInternalServerError)
		}
		return
	}
	vd := views.Data{Yield: post}
	p.ShowView.Render(w, vd)
}
