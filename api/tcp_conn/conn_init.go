package tcp_conn

import (
	"context"
	"github.com/aijie/michat/datas/pb"
	"github.com/aijie/michat/logger"
	"google.golang.org/grpc"
	"net"
)

type ConnServer struct {}


func (s *ConnServer)DeliverMessage(context.Context, *pb.DeliverMessageReq) (*pb.DeliverMessageResp, error){
	return pb.DeliverMessageResp{}, ws_conn.
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, error) {
	//do auth
	return handler(ctx, req)
}

func StartRPCServer() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(UnaryServerInterceptor))
	pb.RegisterConnInitServer(server, )
	server.Serve(listener)
}
