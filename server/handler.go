package server

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	UserId = "user_id"
)

type Handler struct {
	ws         *WsServer
	httpServer *HttpServer
}

func NewHandler(ws *WsServer, hs *HttpServer) *Handler{
	return &Handler{ws, hs}
}

func (h *Handler) MessageHandler(w http.ResponseWriter, r *http.Request) {
	m := &Message{}
	err := json.NewDecoder(r.Body).Decode(m)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid body")
	}
	log.Println("New Message from:", m.UserId)
	err = h.ws.SendMessage(m.UserId, m.Message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}

func (h *Handler) CometHandler(conn *websocket.Conn) {
	userId := conn.Request().URL.Query().Get(UserId)
	defer conn.Close()
	if "" != strings.TrimSpace(userId) {
		millis := time.Now().UnixNano() / 1e6
		c := &Client{userId, millis, conn}
		c.Listen()
	}
}
