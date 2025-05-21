package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func RenderTemplates(w http.ResponseWriter, t string) {
	tmpl, err := template.ParseFiles("./templates/" + t + ".tmpl")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}
