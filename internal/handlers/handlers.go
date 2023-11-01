package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"wb/internal/cache"
	"wb/internal/model"
	"wb/pkg/nats"
)

func Registry() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/upload", Upload)
	http.HandleFunc("/getOrder/", GetOrder)
}

func Index(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("./web/Index.html")
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write(content)

	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func Upload(w http.ResponseWriter, r *http.Request) {
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

	data, err := json.Marshal(order)
	if err != nil {
		log.Println("marshal error: ", err)
	}

	var natsSrv nats.Service
	natsSrv.Connect("produce_upload")
	natsSrv.Publish("uploadFile", data)
	natsSrv.Close()
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := strings.TrimPrefix(r.URL.Path, "/getOrder/")

	_cache := cache.GetCacheInstance()
	order := _cache[orderID]

	orderJSON, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		log.Println("marshal error: ", err)
		http.Error(w, "Ошибка при форматировании заказа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(orderJSON)
}
