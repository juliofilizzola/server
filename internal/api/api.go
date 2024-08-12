package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/juliofilizzola/server/internal/store/pgstore"
)

type apiHandler struct {
	queries    *pgstore.Queries
	r          *chi.Mux
	upgrades   websocket.Upgrader
	subscriber map[string]map[*websocket.Conn]context.CancelFunc
	mutex      *sync.Mutex
}

// func (a apiHandler) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
// type _body struct {
// 	Theme string `json:"theme"`
// 	Name  string `json:"name"`
// }
//
// var body _body
// if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
// 	http.Error(w, err.Error(), http.StatusBadRequest)
// 	return
// }
//
// roomId := uuid.New()
// var sendData a.queries.CreateRoomParams
// sendData.ID = roomId
// sendData.Theme = body.Theme
// sendData.Name = body.Name
//
// room, err := a.queries.CreateRoom(context.Background(), sendData)
// if err != nil {
// 	slog.Error("error creating room", err)
// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 	return
// }
// type response struct {
// 	ID string `json:"id"`
// }
// data, _ := json.Marshal(response{
// 	ID: room.ID.String(),
// })
//
// w.Header().Set("Content-Type", "application/json")
// w.Write(data)

// }

func (a apiHandler) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	type createRoomRequest struct {
		Theme string `json:"theme"`
		Name  string `json:"name"`
	}

	var request createRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := pgstore.CreateRoomParams{

		Theme: request.Theme,
		Name:  request.Name,
	}
	room, err := a.queries.CreateRoom(context.Background(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type createRoomResponse struct {
		ID string `json:"id"`
	}

	response := createRoomResponse{
		ID: room.ID.String(),
	}

	data, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (a apiHandler) handleGetRoom(w http.ResponseWriter, r *http.Request) {

}

func (a apiHandler) handleGetMessages(w http.ResponseWriter, r *http.Request) {

}

func (a apiHandler) handleGetMessage(w http.ResponseWriter, r *http.Request) {

}

func (a apiHandler) handleReactToMessage(w http.ResponseWriter, r *http.Request) {

}

func (a apiHandler) handleRemoveReaction(w http.ResponseWriter, r *http.Request) {

}

func (a apiHandler) handleMarkAsAnswered(w http.ResponseWriter, r *http.Request) {

}

func (a apiHandler) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	rawRoomId := chi.URLParam(r, "room_id")

	roomID, err := uuid.Parse(rawRoomId)

	if err != nil {
		http.Error(w, "invalid room", http.StatusBadRequest)
		return
	}

	_, err = a.queries.GetRoom(context.Background(), roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "room not found", http.StatusNotFound)
			return
		}
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	conn, err := a.upgrades.Upgrade(w, r, nil)
	if err != nil {
		slog.Warn("error upgrading connection", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	defer func(c *websocket.Conn) {
		c.Close()
	}(conn)

	ctx, cancel := context.WithCancel(r.Context())
	a.mutex.Lock()

	if _, ok := a.subscriber[rawRoomId]; ok {
		a.subscriber[rawRoomId][conn] = cancel
	} else {
		slog.Info("new subscriber", rawRoomId)
		a.subscriber[rawRoomId] = make(map[*websocket.Conn]context.CancelFunc)
		a.subscriber[rawRoomId][conn] = cancel
	}

	a.mutex.Unlock()

	<-ctx.Done()
	a.mutex.Lock()
	delete(a.subscriber[rawRoomId], conn)
	a.mutex.Unlock()
}

func (a apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.r.ServeHTTP(w, r)
}

func NewHandler(q *pgstore.Queries) http.Handler {
	a := apiHandler{
		queries: q,
		upgrades: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		subscriber: make(map[string]map[*websocket.Conn]context.CancelFunc),
		mutex:      &sync.Mutex{},
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/subscribe/{room_id}", a.handleSubscribe)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/rooms", func(r chi.Router) {
			r.Post("/", a.handleCreateRoom)
		})

		r.Route("/{room_id}", func(r chi.Router) {
			r.Post("/messages", a.handleCreateRoom)
			r.Get("/messages", a.handleGetRoom)
			r.Route("/messages/{message_id}", func(r chi.Router) {
				r.Get("/{message_id}", a.handleGetMessage)
				r.Patch("/reaction", a.handleReactToMessage)
				r.Delete("/reaction", a.handleRemoveReaction)
				r.Patch("/answered", a.handleMarkAsAnswered)
			})
		})
	})
	a.r = r
	return a
}
