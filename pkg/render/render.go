package render

import (
	"bwa/pkg/config"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

var functions = template.FuncMap{}

func NewTemplate(a *config.AppConfig) {
	app = a
}

func RenderTemplates(w http.ResponseWriter, t string) {

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		newTemplCache, err := CreateTemplateCache()
		if err != nil {
			log.Fatal(err)
		}
		tc = newTemplCache
	}

	tmpl, ok := tc[t]

	if !ok {
		log.Fatal("file Not found in the cache")
	}

	buf := new(bytes.Buffer)

	_ = tmpl.Execute(buf, nil)

	_, err := buf.WriteTo(w)

	if err != nil {
		fmt.Println(err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
