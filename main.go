package main

import (
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
	"time"
)

type TemplateHandler struct {
	Once     sync.Once
	Filename string
	Templ    *template.Template
}

func (t *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.Once.Do(func() {
		t.Templ = template.Must(template.ParseFiles(filepath.Join("template", t.Filename)))
	})

	t.Templ.Execute(w, r)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	room := NewRoom()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", &TemplateHandler{Filename: "index.html"})
	http.Handle("/chat", &TemplateHandler{Filename: "chat.html"})

	http.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
		roomName := r.URL.Query().Get("room")
		if roomName == "" {
			http.Error(w, "Room name required", http.StatusBadRequest)
			return
		}
		realRoom := getRoom(roomName)
		realRoom.ServeHTTP(w, r)
	})

	go room.Run()

	log.Println("Starting server on port:", ":8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
