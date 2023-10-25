package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"wb/internal/model"
)

func Registry() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)
}

func index(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("./web/index.html")
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write(content)

	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		return
	}

	err = json.NewDecoder(file).Decode(&order)
	if err != nil {
		log.Println(err)
		return
	}
}
