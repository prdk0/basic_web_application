package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

func RenderTemplates(w http.ResponseWriter, t string) {
	templCache, err := CreateTemplateCache(w, t)
	if err != nil {
		fmt.Println(err)
		return
	}

	tmpl, ok := templCache[t]

	if !ok {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)

	_ = tmpl.Execute(buf, nil)

	_, err = buf.WriteTo(w)

	if err != nil {
		fmt.Println(err)
	}

}

func CreateTemplateCache(w http.ResponseWriter, fileName string) (map[string]*template.Template, error) {
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
