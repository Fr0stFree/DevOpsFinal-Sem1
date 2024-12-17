package app

import (
	"net/http"

	"project_sem/internal/db"
	"project_sem/internal/handlers"
)

func NewServerRouter(repo db.Repositorier) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v0/prices", handlers.GetProducts(repo))
	mux.HandleFunc("POST /api/v0/prices", handlers.CreateProducts(repo))
	return mux
}
