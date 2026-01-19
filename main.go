package main

import (
	"html/template"
	"net/http"
)

var tpl *template.Template

func main() {
	tpl, _ = template.ParseFiles("templates/index.html")
	sm := http.NewServeMux()
	sm.HandleFunc("/", handlerAdd)
	http.ListenAndServe(":8080", sm)
}

func handlerAdd(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}
