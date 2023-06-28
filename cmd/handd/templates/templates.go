package templates

import (
	"embed"
	"html/template"
	"io"
)

//go:embed views/*.html
var tmplFS embed.FS

type Template struct {
	templates *template.Template
}

func New() *Template {
	templates := template.Must(template.New("").ParseFS(tmplFS, "views/*.html"))
	return &Template{
		templates: templates,
	}
}

// https://stackoverflow.com/a/69244593
func (t *Template) Render(w io.Writer, name string, data interface{}) error {
	tmpl := template.Must(t.templates.Clone())
	t1, err := tmpl.New("").ParseFS(tmplFS, "views/"+name)
	if err != nil {
		return err
	}
	return t1.ExecuteTemplate(w, name, data)
}

func (t *Template) RenderPartial(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}