package model

import (
	"context"
	"errors"
	"github.com/aijie/michat/server/errorcode"
	"github.com/aijie/michat/server/logger"
	"google.golang.org/grpc/metadata"
	"strconv"
)

const (
	CtxAppId     = "app_id"
	CtxUserId    = "user_id"
	CtxDeviceId  = "device_id"
	CtxRequestId = "request_id"
	CtxToken     = "token"
)

var (
	errMetadata = errors.New("metadata error")
	errTokenNotExists = errors.New("token not exist")
	errUnauthorized   = errors.New("error unauthorized")
)

func parseMD(md metadata.MD, target string) (int64, error) {
	values, ok := md[target]
	if !ok && len(values) == 0 {
		return 0, errMetadata
	}
	value, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func GetCtxAppId(ctx context.Context) int64 {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0
	}
	appId, err := parseMD(md, CtxAppId)
	if err != nil {
		logger.Sugar.Error(err)
		return 0
	}
	return appId
}

func GetCtxInfo(ctx context.Context) (int64, int64, int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, 0, 0, errorcode.ErrUnauthorized
	}
	appId, err := parseMD(md, CtxAppId)
	userId, err := parseMD(md, CtxUserId)
	deviceId, err := parseMD(md, CtxDeviceId)
	if err != nil {
		return 0, 0, 0, errorcode.ErrUnauthorized
	}
	return appId, userId, deviceId, nil
}

func GetCtxRequestId(ctx context.Context) int64 {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0
	}
	reqIds, ok := md[CtxRequestId]
	if !ok && len(reqIds) == 0 {
		return 0
	}
	reqId, err := strconv.ParseInt(reqIds[0], 10, 64)
	if err != nil {
		return 0
	}
	return reqId
}

func GetCtxToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errTokenNotExists
	}
	tokens, ok := md[CtxToken]
	if !ok && len(tokens) == 0 {
		return "", errTokenNotExists
	}
	return tokens[0], nil
}

func ContextWithReqId(ctx context.Context, requestId int64) context.Context {
	md := metadata.Pairs(CtxRequestId, strconv.FormatInt(requestId, 10))
	return metadata.NewOutgoingContext(ctx, md)
}
