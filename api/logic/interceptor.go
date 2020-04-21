package logic

import (
	"context"
	"github.com/aijie/michat/datas/model"
	"github.com/aijie/michat/server/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func LogicInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//do auth
	resp, err := handler(ctx, req)
	s, _ := status.FromError(err)
	if s.Code() != 0 && s.Code() < 1000 {
		md, _ := metadata.FromIncomingContext(ctx)
		logger.Logger.Error("logic_int_interceptor", zap.String("method", info.FullMethod), zap.Any("md", md), zap.Any("req", req),
			zap.Any("resp", resp), zap.Error(err), zap.Stack("logic"))
	}
	return resp, err
}

func clientInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error){
	if info.FullMethod != "/pb.LogicClientExt/RegisterDevice" {
		appId, userId, deviceId, err := model.GetCtxInfo(ctx)
		if err != nil {
			return nil, err
		}
		token, err := model.GetCtxToken(ctx)
		if err != nil {
			return nil, err
		}
		err = servie
	}
}

func LogicClientExtInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
}
