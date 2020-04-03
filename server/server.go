package server

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/websocket"
	"log"
	"sync"
	"time"
)

var (
	ErrClientNotFound = errors.New("non-existent client")
)

const (
	ListenInterval = 5
	MessageChanSize = 20
	ClientNums = 20
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

type Client struct {
	UserId    string
	Timestamp int64
	conn      *websocket.Conn
}

type Message struct {
	UserId  string `json:"user_id"`
	Message string `json:"message"`
}

func NewWsServer() *WsServer {
	return &WsServer{
		Clients: make(map[string]*Client),
		Message: make(chan *Message, MessageChanSize),
		AddCli: make(chan *Client, ClientNums),
		DelCli: make(chan *Client, ClientNums),
	}
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
		return ErrClientNotFound
	}
	go client.sendMessage(userId, message)
	return nil
}

func (c *Client) sendMessage(userId, message string) {
	if c.conn != nil {
		c.conn.Write([]byte(message))
	}
}

func (c *Client) heartbeat() error {
	millis := time.Now().UnixNano() / 1e6
	heartbeat := struct {
		Heartbeat int64 `json:heartbeat`
	}{millis}
	bytes, _ := json.Marshal(heartbeat)
	_, err := c.conn.Write(bytes)
	return err
}

func (c *Client) Listen() {
	for range time.Tick(ListenInterval * time.Second) {
		err := c.heartbeat()
		if err != nil {
			log.Println("heartbeat error")
			return
		}
	}
}
