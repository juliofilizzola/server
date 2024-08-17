package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/juliofilizzola/server/internal/store/pgstore"
	"github.com/juliofilizzola/server/internal/utils"
)

type ApiHandler struct {
	Queries     *pgstore.Queries
	Router      *chi.Mux
	Up          websocket.Upgrader
	Subscribers map[string]map[*websocket.Conn]context.CancelFunc
	Mutex       *sync.Mutex
}

func (a ApiHandler) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	type createRoomRequest struct {
		Theme string `json:"theme"`
		Name  string `json:"name"`
	}

	var request createRoomRequest

	if err := utils.ParseJson(r, &request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := pgstore.CreateRoomParams{
		Theme: request.Theme,
		Name:  request.Name,
	}

	room, err := a.Queries.CreateRoom(context.Background(), params)

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

	utils.WriteJsonResponse(w, http.StatusCreated, data)
}

func (a ApiHandler) HandleGetRoom(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")
	roomID, err := uuid.Parse(ID)
	if err != nil {
		http.Error(w, "invalid room ID", http.StatusBadRequest)
		return
	}

	room, err := a.Queries.GetRoomByID(context.Background(), roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "room not found", http.StatusNotFound)
			return
		}
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(room)

	utils.WriteJsonResponse(w, http.StatusOK, data)
}

func (a ApiHandler) HandlerGetRooms(w http.ResponseWriter, _ *http.Request) {
	params := pgstore.ListRoomsParams{
		Limit:  10,
		Offset: 0,
	}
	rooms, err := a.Queries.ListRooms(context.Background(), params)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(rooms)

	utils.WriteJsonResponse(w, http.StatusOK, data)
}

func (a ApiHandler) HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	idMessage := chi.URLParam(r, "id")
	messageID, err := uuid.Parse(idMessage)
	if err != nil {
		http.Error(w, "invalid message ID", http.StatusBadRequest)
		return
	}

	message, err := a.Queries.GetMessage(context.Background(), messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "message not found", http.StatusNotFound)
			return
		}
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(message)

	utils.WriteJsonResponse(w, http.StatusOK, data)
}

func (a ApiHandler) HandleReactToMessage(w http.ResponseWriter, r *http.Request) {
	idMessage := chi.URLParam(r, "id")
	messageID, err := uuid.Parse(idMessage)
	if err != nil {
		http.Error(w, "invalid message ID", http.StatusBadRequest)
		return
	}

	_, err = a.Queries.GetMessage(context.Background(), messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "message not found", http.StatusNotFound)
			return
		}
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	reaction, err := a.Queries.AddReactionFromMessage(context.Background(), messageID)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(reaction)

	utils.WriteJsonResponse(w, http.StatusOK, data)
}

func (a ApiHandler) HandleRemoveReaction(w http.ResponseWriter, r *http.Request) {
	idMessage := chi.URLParam(r, "id")
	messageID, err := uuid.Parse(idMessage)
	if err != nil {
		http.Error(w, "invalid message ID", http.StatusBadRequest)
		return
	}

	_, err = a.Queries.GetMessage(context.Background(), messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "message not found", http.StatusNotFound)
			return
		}
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	reaction, err := a.Queries.RemoveReactionFromMessage(context.Background(), messageID)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(reaction)

	utils.WriteJsonResponse(w, http.StatusOK, data)
}

func (a ApiHandler) HandleMarkAsAnswered(w http.ResponseWriter, r *http.Request) {
	idMessage := chi.URLParam(r, "id")
	messageID, err := uuid.Parse(idMessage)
	if err != nil {
		http.Error(w, "invalid message ID", http.StatusBadRequest)
		return
	}

	_, err = a.Queries.GetMessage(context.Background(), messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "message not found", http.StatusNotFound)
			return
		}
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	resolved, err := a.Queries.AnswerMessage(context.Background(), messageID)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(resolved)

	utils.WriteJsonResponse(w, http.StatusOK, data)
}

func (a ApiHandler) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	rawRoomId := chi.URLParam(r, "room_id")

	roomID, err := uuid.Parse(rawRoomId)

	if err != nil {
		http.Error(w, "invalid room", http.StatusBadRequest)
		return
	}

	_, err = a.Queries.GetRoomByID(context.Background(), roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "room not found", http.StatusNotFound)
			return
		}
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	conn, err := a.Up.Upgrade(w, r, nil)
	if err != nil {
		slog.Warn("error upgrading connection", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			slog.Error("error closing connection", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
	}(conn)

	ctx, cancel := context.WithCancel(r.Context())
	a.Mutex.Lock()

	if _, ok := a.Subscribers[rawRoomId]; ok {
		a.Subscribers[rawRoomId][conn] = cancel
	} else {
		slog.Info("new subscribers", rawRoomId)
		a.Subscribers[rawRoomId] = make(map[*websocket.Conn]context.CancelFunc)
		a.Subscribers[rawRoomId][conn] = cancel
	}

	a.Mutex.Unlock()

	<-ctx.Done()
	a.Mutex.Lock()
	delete(a.Subscribers[rawRoomId], conn)
	a.Mutex.Unlock()
}
