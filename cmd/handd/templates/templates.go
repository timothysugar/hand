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
    tmpl, err := tmpl.ParseFS(tmplFS, "views/"+name)
    if err != nil {
        return err
    }
    return tmpl.ExecuteTemplate(w, name, data)
}
