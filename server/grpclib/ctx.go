package grpclib

import (
	"context"
	"google.golang.org/grpc/metadata"
	"strconv"
)

const (
	CtxAppId = "app_id"
	CtxUserId = "user_id"
	CtxDeviceId = "device_id"
	CtxRequestId = "request_id"
	CtxToken = "token"
)

func GetCtxAppId(ctx context.Context) int64 {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0
	}
	appIds, ok := md[CtxAppId]
	if !ok && len(appIds) == 0{
		return 0
	}
	appId, err := strconv.ParseInt(appIds[0], 10, 64)
	if err != nil {
		return 0
	}
	return appId
}

func GetCtxInfo(ctx context.Context) (int64, int64, int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, 0, 0,
	}
	userIds, ok := md[CtxUserId]
	if !ok && len(userIds) == 0{
		return 0
	}
	userId, err := strconv.ParseInt(userIds[0], 10, 64)
	if err != nil {
		return 0
	}
	return userId
}

func GetCtxRequestId(ctx context.Context) int64 {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0
	}
	reqIds, ok := md[CtxRequestId]
	if !ok && len(reqIds) == 0{
		return 0
	}
	reqId, err := strconv.ParseInt(reqIds[0], 10, 64)
	if err != nil {
		return 0
	}
	return reqId
}

func GetCtxToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	tokens, ok := md[CtxToken]
	if !ok && len(tokens) == 0{
		return ""
	}
	return tokens[0]
}

func ContextWithRequestId(ctx context.Context, requestId int64) context.Context{
	md := metadata.Pairs(CtxRequestId, strconv.FormatInt(requestId, 10))
	return metadata.NewOutgoingContext(ctx, md)
}
