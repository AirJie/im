package ws_service

import (
	"github.com/aijie/michat/datas/pb"
	"github.com/aijie/michat/logger"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/status"

	//"golang.org/x/net/websocket"
	"github.com/gorilla/websocket"
	"log"
)

type Linker struct {
	Conn     *websocket.Conn
	AppId    int64
	UserId   int64
	DeviceId int64
	Message  chan *Message
}

type messageType int

const (
	UnknownType messageType = iota
	SignInType
	SyncType
	HeartbeatType
	MessageType
)

func NewLink(conn *websocket.Conn, appId, userId, deviceId int64) *Linker {
	return &Linker{
		Conn:     conn,
		AppId:    appId,
		UserId:   userId,
		DeviceId: deviceId,
	}
}

func (r *Linker) DoConn() {
	for {
		_, data, err := r.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		r.HandleMessage(data)
	}
}

func (r *Linker) Heartbeat(input pb.Input) {
	r.Output(pb.SessionType_Heartbeat, 0, nil, nil)
}

func (r *Linker) HandleMessage(data []byte) {
	input := pb.Input{}
	proto.Unmarshal(data, &input)
	switch input.Type {
	case pb.SessionType_Unknown:
		log.Fatal("Get unknow type of message")
	case pb.SessionType_Heartbeat:
		r.Heartbeat(input)
	case pb.SessionType_Sync:
	case pb.SessionType_SignIn:
	case pb.SessionType_MessageStream:
	default:
		logger.Sugar.Error("input type error")
	}
}

func (r *Linker) Output(ptype pb.SessionType, requestId int64, err error, message proto.Message) {
	var out = pb.Output{
		Type: ptype,
		Id:   requestId,
	}
	if err != nil {
		status, _ := status.FromError(err)
		out.Errcode = int32(status.Code())
		out.Message = status.Message()
	}

	if message != nil {
		out.Data, err = proto.Marshal(message)
		if err != nil {
			logger.Sugar.Error(err)
			return
		}
	}
	bytesData, err := proto.Marshal(message)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	err = r.Conn.WriteMessage(websocket.BinaryMessage, bytesData)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
}

