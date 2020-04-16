package ws_conn

import (
	"context"
	"github.com/aijie/michat/datas/pb"
	"github.com/aijie/michat/server/grpclib"
	"github.com/aijie/michat/server/logger"
	"github.com/aijie/michat/server/rpc_cli"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/status"

	"github.com/gorilla/websocket"
	"log"
)

type Linker struct {
	Conn     *websocket.Conn
	AppId    int64
	UserId   int64
	DeviceId int64
	Message  chan *pb.Message
}

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
		logger.Sugar.Error("Get unknow type of message")
	case pb.SessionType_Heartbeat:
		r.Heartbeat(input)
	case pb.SessionType_Sync:
		r.Sync(input)
	case pb.SessionType_SignIn:
		logger.Logger.Info("sign in")
	case pb.SessionType_MessageStream:
		r.MessageACK(input)
	default:
		logger.Sugar.Error("input type error")
	}
}

func (r *Linker) Sync(input pb.Input) {
	sync := pb.SyncInput{}
	err := proto.Unmarshal(input.Data, &sync)
	if err != nil {
		logger.Sugar.Error(err)
		r.Release()
		return
	}
	resp, err := rpc_cli.LogicIntClient.Sync(grpclib.ContextWithRequestId(context.TODO(), input.RequestId), &pb.SyncReq{
		AppId:    r.AppId,
		UserId:   r.UserId,
		DeviceId: r.DeviceId,
		Seq:      sync.Seq,
	})
	var message proto.Message
	if err == nil {
		message = &pb.SyncOutput{Messages: resp.Messages}
	}
	r.Output(pb.SessionType_Sync, input.RequestId, err, message)
}

func (r *Linker) MessageACK(input pb.Input) {
	var ack pb.MessageACK
	err := proto.Unmarshal(input.Data, &ack)
	if err != nil {
		logger.Sugar.Error(err)
		r.Release()
		return
	}
	rpc_cli.LogicIntClient.MessageACK(grpclib.ContextWithRequestId(context.TODO(), input.RequestId), &pb.MessageACKReq{
		AppId:     r.AppId,
		UserId:    r.UserId,
		DeviceId:  r.UserId,
		DeviceAck: ack.DeviceAck,
		Timestamp: ack.ReceiveTime,
	})
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

func (r *Linker) Release() {
	err := r.Conn.Close()
	if err != nil {
		logger.Sugar.Error(err)
	}
	//offline
	rpc_cli.LogicIntClient.Offline(context.TODO(), &pb.OfflineReq{
		AppId:    r.AppId,
		UserId:   r.UserId,
		DeviceId: r.DeviceId,
	})
}
