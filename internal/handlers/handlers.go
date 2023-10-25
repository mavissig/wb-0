package handlers

import (
	"encoding/json"
	"io"
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

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("read upload: ", err.Error())
	}

	err = json.Unmarshal(bytes, &order)
	if err != nil {
		log.Println(err)
	}
}
