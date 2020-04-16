package ws_conn

import (
	"context"
	"github.com/aijie/michat/datas/pb"
	"github.com/aijie/michat/link/ws_conn"
	"github.com/aijie/michat/server/logger"
	"google.golang.org/grpc"
	"net"
)

type WsConnServer struct {}


func (s *WsConnServer)DeliverMessage(ctx context.Context, in *pb.DeliverMessageReq) (*pb.DeliverMessageResp, error){
	logger.Logger.Info("deliver message")
	out := new(pb.DeliverMessageResp)
	return out, ws_conn.DeliverMessage(ctx, in)
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//do auth
	return handler(ctx, req)
}

func StartRPCServer() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(UnaryServerInterceptor))
	var connServer WsConnServer
	pb.RegisterConnInitServer(server, &connServer)
	server.Serve(listener)
}
