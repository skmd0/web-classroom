package controllers

import (
	"net/http"
	"wiki/views"
)

type Static struct {
	Index      *views.View
	Categories *views.View
	Topic      *views.View
	Post       *views.View
	Control    *views.View
	NotFound   *views.View
}

func NewStatic() *Static {
	return &Static{
		Index:      views.NewView("bulma", "static/index"),
		Categories: views.NewView("bulma", "categories"),
		Topic:      views.NewView("bulma", "topic"),
		Post:       views.NewView("bulma", "post"),
		Control:    views.NewView("bulma", "control"),
		NotFound:   views.NewView("bulma", "notfound"),
	}
}

func (s *Static) NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		s.NotFound.Render(w, r)
	}
}
