package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
)

//go:embed views/pages/*.go.html views/layout.go.html
var tmplFS embed.FS

var tAll *template.Template
var tPages map[string]*template.Template = make(map[string]*template.Template)

func init() {
    // Generate a template from everything to pick up all partials in all templates
	tAll = template.Must(template.ParseFS(tmplFS, "views/pages/*.go.html", "views/layout.go.html"))
	log.Println(tAll.DefinedTemplates())

	// Walk the template filesystem
	// for each file in pages, create a template with the content of the page file named "content"
	// and with the template named the same as the page file
	fs.WalkDir(tmplFS, "views/pages", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if (d.IsDir()) { return nil }
		t := template.Must(template.New(d.Name()).ParseFS(tmplFS, "views/layout.go.html", path ))
		log.Printf("templates defined in template named %s, path: %s: %s", t.Name(), path,  t.DefinedTemplates())
		tPages[t.Name()] = t
		return nil
})
}

func RenderFragment(w io.Writer, name string, data interface{}) error {
	log.Printf("rendering %v", name)
	log.Println(tAll.DefinedTemplates())
	
	res := tAll.ExecuteTemplate(w, name, data)
	log.Printf("rendered %s", name)
	return res
}

func RenderPage(w io.Writer, pageName string, data interface{}) error {

	t, ok := tPages[pageName]
	if (!ok) {
		return fmt.Errorf("template %s not found", pageName)
	}
	log.Printf("rendering page %s, templates: %s, vm: %+v", pageName, t.DefinedTemplates(), data)
	t.Execute(w, data)

	res := t.ExecuteTemplate(w, "layout.go.html", data)
	log.Printf("rendered %s", pageName)
	return res
}