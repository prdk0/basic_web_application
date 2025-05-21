package handlers

import (
	"bwa/pkg/render"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "home.page")
}

func About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "about.page")
}
