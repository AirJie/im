package server

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"net/http"
	"strings"
	"time"
)
var (
	UserId = "user_id"
)

type Handler struct {
	ws *WsServer
}

func (h *Handler) MessageHandler(w http.ResponseWriter, r *http.Request) {
	m := &Message{}
	err := json.NewDecoder(r.Body).Decode(m)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid body")
	}
	h.ws.SendMessage(m.UserId, m.Message)
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
