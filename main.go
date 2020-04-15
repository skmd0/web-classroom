package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"net/http"
	"wiki/controllers"
	"wiki/middleware"
	"wiki/util"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "future"
	dbname   = "postgres"
)

func main() {
	// DATABASE SETUP
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user,
		password, dbname)
	services, err := util.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	//services.DestructiveReset()
	services.AutoMigrate()

	// CONTROLLERS
	r := mux.NewRouter()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	postsC := controllers.NewPosts(services.Post, r)

	// MIDDLEWARE
	// applyUserMw is required pretty much everywhere because
	// navbar needs user info to render correct elements
	applyUserMw := middleware.User{UserService: services.User}

	// requireUserMw redirects to login page if user is not logged in.
	// When the user is logged in, it applies the user data to the context.
	// That's why it embeds applyUserMw, which sets the user.
	requireUserMw := middleware.RequireUser{User: applyUserMw}

	// ROUTING
	r.NotFoundHandler = staticC.NotFoundHandler()
	r.Handle("/", applyUserMw.Apply(staticC.Index)).Methods("GET")

	// User
	r.HandleFunc("/login", applyUserMw.ApplyFn(usersC.Login)).Methods("GET")
	r.HandleFunc("/login", applyUserMw.ApplyFn(usersC.LoginUser)).Methods("POST")
	r.HandleFunc("/signup", applyUserMw.ApplyFn(usersC.New)).Methods("GET")
	r.HandleFunc("/signup", applyUserMw.ApplyFn(usersC.Create)).Methods("POST")

	// Posts
	r.Handle("/post/new", requireUserMw.Apply(postsC.New)).Methods("GET")
	r.HandleFunc("/posts", requireUserMw.ApplyFn(postsC.Create)).Methods("POST")
	r.HandleFunc("/posts", requireUserMw.ApplyFn(postsC.PostIndex)).Methods("GET")
	r.HandleFunc("/post/{id:[0-9]+}", applyUserMw.ApplyFn(postsC.Show)).Methods("GET").Name("show_post")
	r.HandleFunc("/post/{id:[0-9]+}/edit", requireUserMw.ApplyFn(postsC.Edit)).Methods("GET")
	r.HandleFunc("/post/{id:[0-9]+}/update", requireUserMw.ApplyFn(postsC.Update)).Methods("POST")
	r.HandleFunc("/post/{id:[0-9]+}/delete", requireUserMw.ApplyFn(postsC.Delete)).Methods("POST")

	fmt.Println("Running the server on :3000")
	http.ListenAndServe(":3000", r)
}
