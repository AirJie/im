package ws_service

import (
	"errors"
	//"github.com/gorilla/websocket"
	"log"
	"sync"
)

var (
	ErrClientNotFound = errors.New("non-existent client")
)

const (
	ListenInterval  = 5
	MessageChanSize = 20
	ClientNums      = 20
)

type WsService interface {
	Start()
	SendMessage(string, string) error
}

type WsServer struct {
	mutex   sync.Mutex
	Clients map[string]*Client
	Message chan *Message
	AddCli  chan *Client
	DelCli  chan *Client
}

type HttpServer struct {
	wsServer *WsServer
}


func NewWsServer() *WsServer {
	return &WsServer{
		Clients: make(map[string]*Client),
		Message: make(chan *Message, MessageChanSize),
		AddCli:  make(chan *Client, ClientNums),
		DelCli:  make(chan *Client, ClientNums),
	}
}

func NewHttpServer(ws *WsServer) *HttpServer {
	return &HttpServer{wsServer: ws}
}

func (ws *WsServer) addClient(c *Client) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	if _, ok := ws.Clients[c.UserId]; !ok {
		ws.Clients[c.UserId] = c
	}
	return nil
}

func (ws *WsServer) delClient(userId string) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	if _, ok := ws.Clients[userId]; ok {
		//map delete
		delete(ws.Clients, userId)
	}
	return nil
}

func (ws *WsServer) Start() {
	for {
		select {
		case msg := <-ws.Message:
			ws.SendMessage(msg.UserId, msg.Message)
		case c := <-ws.AddCli:
			ws.addClient(c)
		case c := <-ws.DelCli:
			ws.delClient(c.UserId)
		}
	}
}

func (ws *WsServer) SendMessage(userId, message string) error {
	client := ws.Clients[userId]
	if client == nil {
		log.Fatal("There is no valid Client")
		return ErrClientNotFound
	}
	go client.sendMessage(userId, message)
	return nil
}

