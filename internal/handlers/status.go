package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Посмотрите вывод в консоль этих двух строк
	// fmt.Printf("Response: \n%+v\n\n", w)
	// fmt.Printf("Request: \n%+v\n\n", r)

	slog.Info("Received health check request", slog.String("method", r.Method), slog.String("url", r.URL.String()))
	render.JSON(w, r, HealthResponse{Status: "OK"})
}
