package views

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

var (
	LayoutDir   = "views/layouts/"
	TemplateDir = "views/"
	TemplateExt = ".gohtml"
)

func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)

	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, nil)
}

func (v *View) Render(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	switch data.(type) {
	case Data:
		// do nothing because data is already correctly set
	default:
		data = Data{Yield: data}
	}
	var buf bytes.Buffer
	if err := v.Template.ExecuteTemplate(&buf, v.Layout, data); err != nil {
		log.Println("Rendering template failed.", err)
		http.Error(w, "Something went wrong! If the problem persists please email us.",
			http.StatusInternalServerError)
		return
	}
	_, _ = io.Copy(w, &buf)
}

// layoutFiles returns a slice of string names for layout files used in rendering the website
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// addTemplatePath takes in a slice of strings
// representing file paths for templates and it prepends
// the TemplateDir directory to each string in a slice
//
// {"home"} -> {"views/home"}
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

// addTemplateExt takes in a slice of strings
// representing file paths for templates and it appends
// the TemplateExt extension to each string in the slice
//
// {"home"} -> {"home.gohtml"}
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
