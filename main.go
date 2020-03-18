package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"wiki/controllers"
)

func main() {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.NotFoundHandler = staticC.NotFoundHandler()
	r.Handle("/", staticC.Index).Methods("GET")
	r.Handle("/content/{category}/{topic}/post/{id}", staticC.Post).Methods("GET")
	r.Handle("/content/categories", staticC.Categories).Methods("GET")
	r.Handle("/content/{category}/", staticC.Post).Methods("GET")
	r.Handle("/user/{username}", staticC.Control).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	http.ListenAndServe(":3000", r)
}
