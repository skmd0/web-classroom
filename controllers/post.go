package controllers

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"wiki/context"
	"wiki/models"
	"wiki/models/keywords"
	"wiki/models/posts"
	"wiki/views"
)

func NewPosts(ps posts.PostService, r *mux.Router) *Posts {
	return &Posts{
		New:           views.NewView("bulma", "posts/new"),
		ShowView:      views.NewView("bulma", "posts/show"),
		EditView:      views.NewView("bulma", "posts/edit"),
		PostIndexView: views.NewView("bulma", "posts/index"),
		HomepageView:  views.NewView("bulma", "index"),
		ps:            ps,
		r:             r,
	}
}

type Posts struct {
	New           *views.View
	ShowView      *views.View
	EditView      *views.View
	PostIndexView *views.View
	HomepageView  *views.View
	ps            posts.PostService
	r             *mux.Router
}

type NewPostForm struct {
	Title    string `schema:"title"`
	Content  string `schema:"content"`
	Keywords []KeywordForm
}

type KeywordForm struct {
	Title      string
	Definition string
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
		p.New.Render(w, r, vd)
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

	var parsedKeywords []keywords.Keyword
	for _, key := range form.Keywords {
		if key.Title == "" || key.Definition == "" {
			continue
		}
		key := keywords.Keyword{
			Title:       key.Title,
			Description: key.Definition,
		}
		parsedKeywords = append(parsedKeywords, key)
	}
	err := p.ps.CreatePost(post, &parsedKeywords)
	if err != nil {
		vd.SetAlert(err)
		p.New.Render(w, r, vd)
		return
	}

	url, err := p.r.Get("show_post").URL("id", strconv.Itoa(int(post.ID)))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// /GET /
func (p *Posts) Homepage(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	user := context.User(r.Context())
	if user != nil {
		postsDB, err := p.ps.ByUserIdWithLimit(user.ID, 10)
		if err == nil {
			vd.Yield = postsDB
		}
	}
	p.HomepageView.Render(w, r, vd)
}

// /GET /posts
func (p *Posts) PostIndex(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())

	postsDB, err := p.ps.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	vd := views.Data{Yield: postsDB}
	p.PostIndexView.Render(w, r, vd)
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

	keys, err := p.ps.GetKeywords(post.ID)

	// generate breadcrumbs for navbar
	var vd views.Data
	pages := make([]views.Page, 0)
	homeLink := views.Page{
		Title: "Home",
		URL:   "/",
	}
	postsLink := views.Page{
		Title: "Posts",
		URL:   "/posts",
	}
	currentPage := views.Page{
		Title: post.Title,
		URL:   r.URL.Path,
	}
	pages = append(pages, homeLink, postsLink, currentPage)
	vd.Breadcrumbs = views.Breadcrumbs{
		Pages:       pages,
		LastPageKey: post.Title,
	}

	// generate HTML from markdown
	parsedContent := strings.ReplaceAll(post.Content, "\n", `

`)
	md := []byte(parsedContent)
	html := markdown.ToHTML(md, nil, nil)
	post.ContentHTML = template.HTML(html)

	content := struct {
		Post *posts.Post
		Keys *[]keywords.Keyword
	}{
		Post: post,
		Keys: keys,
	}
	vd.Yield = content
	p.ShowView.Render(w, r, vd)
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
	p.EditView.Render(w, r, vd)
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
		p.EditView.Render(w, r, vd)
		return
	}
	post.Title = form.Title
	post.Content = form.Content
	err = p.ps.Update(post)
	if err != nil {
		log.Println(err)
		vd.SetAlert(err)
		p.EditView.Render(w, r, vd)
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
		p.EditView.Render(w, r, vd)
		return
	}
	http.Redirect(w, r, "/posts", http.StatusFound)
}
