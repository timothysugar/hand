package templates

import (
	"embed"
	"html/template"
	"io"
	"log"
)

//go:embed views/*.html
var tmplFS embed.FS

type Template struct {
	templates *template.Template
}

func New() *Template {
	templates := template.Must(template.New("").ParseFS(tmplFS, "views/*.html"))
	log.Printf("Found templates%s", templates.DefinedTemplates())
	return &Template{
		templates: templates,
	}
}

// https://stackoverflow.com/a/69244593
func (t *Template) Render(w io.Writer, name string, data any) error {
	tmpl := template.Must(t.templates.Clone())
	tmpl, err := tmpl.ParseFS(tmplFS, "views/"+name)
	if err != nil {
		return err
	}
	log.Println("Rendering template: ", name, " with data: ", data)
	return tmpl.ExecuteTemplate(w, name, data)
}
