package handlers

import (
	"net/http"
	"os"
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
	return
}
