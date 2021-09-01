package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/go-yaml/yaml"
	"go.uber.org/zap"
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
	configYaml, err := ioutil.ReadFile("./zap_config.yaml")
	if err != nil {
		panic(err)
	}
	var zapConfig zap.Config
	if err := yaml.Unmarshal(configYaml, &zapConfig); err != nil {
		panic(err)
	}
	logger, _ := zapConfig.Build()

	logger.Info("start study-go-chat")

	r := NewRoom()
	http.Handle("/", &templateHeader{filename: "chat.html"})
	http.Handle("/room", r)

	// run chatroom
	go r.run()

	// run webserver
	url := "http://localhost:8080"
	logger.Info("run server", zap.String("url", url))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
