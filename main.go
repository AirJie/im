package main

import (
	"github.com/aijie/michat/ws_service"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)


func main() {
	ws := ws_service.NewWsServer()
	htps := ws_service.NewHttpServer(ws)
	h := ws_service.NewHandler(ws, htps)

	go ws.Start()
	r := mux.NewRouter()
	r.HandleFunc("/message", h.MessageHandler)
	r.Handle("/com", websocket.Handler(h.CometHandler))
	r.Headers("Content-type", "application-json; charset=utf-8")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}
