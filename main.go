package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHeader struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
	fmt.Println("start study-go-chat")
	r := NewRoom()
	http.Handle("/", &templateHeader{filename: "chat.html"})
	http.Handle("/room", r)

	// run chatroom
	go r.run()

	// run webserver
	fmt.Println("run server http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
