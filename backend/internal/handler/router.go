package handler

import (
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(roomH *RoomHandler, playerH *PlayerHandler, wsH *WSHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-Session-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Post("/rooms", roomH.Create)
		r.Get("/rooms/{id}", roomH.Get)
		r.Post("/players", playerH.Create)
	})

	// WebSocket
	r.Get("/ws/{roomID}", wsH.Handle)

	// Serve frontend static files (production)
	staticDir := "./static"
	if _, err := os.Stat(staticDir); err == nil {
		fileServer(r, staticDir)
	}

	return r
}

func fileServer(r chi.Router, staticDir string) {
	root := http.Dir(staticDir)
	fsHandler := http.FileServer(root)

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// If it's an API or WS path, skip
		if strings.HasPrefix(r.URL.Path, "/api") || strings.HasPrefix(r.URL.Path, "/ws") {
			http.NotFound(w, r)
			return
		}

		// Try to serve the file
		f, err := fs.Stat(os.DirFS(staticDir), strings.TrimPrefix(r.URL.Path, "/"))
		if err != nil || f.IsDir() {
			// SPA fallback: serve index.html
			http.ServeFile(w, r, staticDir+"/index.html")
			return
		}
		fsHandler.ServeHTTP(w, r)
	})
}
