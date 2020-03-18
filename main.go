package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"wiki/controllers"
	"wiki/views"
)

var (
	indexView      *views.View
	categoriesView *views.View
	topicView      *views.View
	postView       *views.View
	controlView    *views.View
	loginView      *views.View
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	panicIfErr(indexView.Render(w, nil))
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	panicIfErr(postView.Render(w, nil))
}

func categories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	panicIfErr(categoriesView.Render(w, nil))
}

func topics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	panicIfErr(topicView.Render(w, nil))
}

func control(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	panicIfErr(controlView.Render(w, nil))
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	panicIfErr(loginView.Render(w, nil))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>404</h1>")
}

func main() {
	indexView = views.NewView("bulma", "views/index.gohtml")
	categoriesView = views.NewView("bulma", "views/categories.gohtml")
	topicView = views.NewView("bulma", "views/topic.gohtml")
	postView = views.NewView("bulma", "views/post.gohtml")
	controlView = views.NewView("bulma", "views/control.gohtml")
	loginView = views.NewView("bulma", "views/login.gohtml")
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFound)
	r.HandleFunc("/", index)
	r.HandleFunc("/content/{category}/{topic}/post/{id}", post)
	r.HandleFunc("/content/categories", categories)
	r.HandleFunc("/content/{category}/", topics)
	r.HandleFunc("/user/{username}", control)
	r.HandleFunc("/login", login)
	r.HandleFunc("/signup", usersC.New)
	http.ListenAndServe(":3000", r)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
