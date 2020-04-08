package ws_service

import "golang.org/x/net/websocket"

type Link struct {
	Conn     *websocket.Conn
	AppId    int64
	UserId   int64
	DeviceId int64
	Message chan *Message
}

type Message struct {
	UserId  string `json:"user_id"`
	Message string `json:"message"`
}

func NewLink(conn *websocket.Conn, appId, userId, deviceId int64) *Link {
	return &Link{
		Conn:     conn,
		AppId:    appId,
		UserId:   userId,
		DeviceId: deviceId,
	}
}

func (lk *Link) DoConn() {
	for {
		lk.Conn.
		//select {
		//case msg := <- lk.Message:
		//	//ws.SendMessage(msg.UserId, msg.Message)
		//}
	}
}
