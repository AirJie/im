package server

import (
	"context"
	"github.com/aijie/michat/datas/pb"
	"github.com/aijie/michat/logger"
	"google.golang.org/grpc"
)
func interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	//err := invoker(ctx, method, req, reply, cc, opts...)
	//return gerrors.WrapRPCError(err)
	return nil
}

var (
	LogicIntClient pb.LogicIntClient
)

func InitLogicIntClient(add string) {
	client, err := grpc.DialContext(context.TODO(), add, grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor))
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	LogicIntClient = pb.NewLogicIntClient(client)
}

