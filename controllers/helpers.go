package controllers

import (
	"github.com/gorilla/schema"
	"net/http"
	"net/url"
)

func ParseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	// r.PostForm contains only POST form data
	return ParseValues(r.PostForm, dst)
}

func ParseURLParams(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	// r.Form contains URL param values
	return ParseValues(r.Form, dst)
}

func ParseValues(values url.Values, dst interface{}) error {
	dec := schema.NewDecoder()
	if err := dec.Decode(dst, values); err != nil {
		return err
	}
	return nil
}
