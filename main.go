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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user,
		password, dbname)
	services, err := util.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	//services.DestructiveReset()
	services.AutoMigrate()

	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	postsC := controllers.NewPosts(services.Post, r)

	r.NotFoundHandler = staticC.NotFoundHandler()
	r.Handle("/", staticC.Index).Methods("GET")
	r.Handle("/content/categories", staticC.Categories).Methods("GET")
	r.Handle("/user/{username}", staticC.Control).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("GET")
	r.HandleFunc("/login", usersC.LoginUser).Methods("POST")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/cookie", usersC.CookieTest).Methods("GET")

	// Posts
	requireUserMw := middleware.RequireUser{UserService: services.User}
	r.Handle("/post/new", requireUserMw.Apply(postsC.New)).Methods("GET")
	r.HandleFunc("/posts", requireUserMw.ApplyFn(postsC.Create)).Methods("POST")
	r.HandleFunc("/post/{id:[0-9]+}", postsC.Show).Methods("GET").Name("show_post")
	r.HandleFunc("/post/{id:[0-9]+}/edit", requireUserMw.ApplyFn(postsC.Edit)).Methods("GET")
	r.HandleFunc("/post/{id:[0-9]+}/update", requireUserMw.ApplyFn(postsC.Update)).Methods("POST")

	fmt.Println("Running the server on :3000")
	http.ListenAndServe(":3000", r)
}
