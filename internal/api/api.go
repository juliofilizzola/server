package api

import (
	"context"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
	"github.com/juliofilizzola/server/internal/handler"
	"github.com/juliofilizzola/server/internal/store/pgstore"
)

const (
	maxCorsAge = 300 // Valor máximo não ignorado por nenhum dos principais navegadores
)

func NewHandler(q *pgstore.Queries) http.Handler {
	a := handler.ApiHandler{
		Queries: q,
		Up: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Subscribers: make(map[string]map[*websocket.Conn]context.CancelFunc),
		Mutex:       &sync.Mutex{},
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           maxCorsAge, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/subscribe/{room_id}", a.HandleSubscribe)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/rooms", func(r chi.Router) {
			r.Post("/", a.HandleCreateRoom)
			r.Get("/", a.HandlerGetRooms)
		})

		r.Route("/{room_id}", func(r chi.Router) {
			r.Post("/messages", a.HandleCreateRoom)
			r.Get("/messages", a.HandleGetRoom)
			r.Route("/messages/{message_id}", func(r chi.Router) {
				r.Get("/{message_id}", a.HandleGetMessages)
				r.Patch("/reaction", a.HandleReactToMessage)
				r.Delete("/reaction", a.HandleRemoveReaction)
				r.Patch("/answered", a.HandleMarkAsAnswered)
			})
		})
	})
	a.Router = r
	return a.Router
}
