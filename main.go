package main

import (
	"github.com/aijie/michat/server"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)


func main() {
	ws := server.NewWsServer()
	htps := server.NewHttpServer(ws)
	h := server.NewHandler(ws, htps)

	go ws.Start()
	r := mux.NewRouter()
	r.HandleFunc("/message", h.MessageHandler)
	r.Handle("/com", websocket.Handler(h.CometHandler))
	r.Headers("Content-type", "application-json; charset=utf-8")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}
