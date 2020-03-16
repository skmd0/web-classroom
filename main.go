package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Contact</h1>")
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Home</h1>")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	http.ListenAndServe(":3000", r)
}
