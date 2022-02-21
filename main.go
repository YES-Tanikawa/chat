package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}

func (t *templateHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func () {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
	t := &templateHandler{filename: "chat.html"}
	test := &templateHandler{filename: "test.html"}
	//ルート
	http.HandleFunc("/", t.ServerHTTP)
	http.HandleFunc("/test", test.ServerHTTP)
	//Webサーバーを起動
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
