package server

import (
	"net/http"
	"wb/internal/handlers"
)

func RunServer() {
	handlers.Registry()

	http.ListenAndServe(":80", nil)
}
