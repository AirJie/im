package tcp_conn

import (
	"context"
	"github.com/aijie/michat/config"
	"github.com/aijie/michat/datas/model"
	"github.com/aijie/michat/datas/pb"
	"github.com/aijie/michat/server/logger"
	"github.com/aijie/michat/server/rpc_cli"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

const PreConn = -1 // 设备第二次重连时，标记设备的上一条连接

type handler struct{}

var Handler = new(handler)

// Handler 处理客户端的上行包
func (h *handler) Handler(ctx *TcpConn, bytes []byte) {
	var input pb.Input
	err := proto.Unmarshal(bytes, &input)
	if err != nil {
		logger.Logger.Error("unmarshal error", zap.Error(err))
		ctx.Release()
		return
	}

	// 对未登录的用户进行拦截
	if input.Type != pb.SessionType_SignIn && ctx.IsSignIn == false {
		// 应该告诉用户没有登录
		ctx.Release()
		return
	}

	switch input.Type {
	case pb.SessionType_SignIn:
		h.SignIn(ctx, input)
	case pb.SessionType_Sync:
		h.Sync(ctx, input)
	case pb.SessionType_Heartbeat:
		h.Heartbeat(ctx, input)
	case pb.SessionType_MessageStream:
		h.MessageACK(ctx, input)
	default:
		logger.Logger.Error("handler switch other")
	}
	return
}

// SignIn 登录
func (*handler) SignIn(ctx *TcpConn, input pb.Input) {
	var signIn pb.SignInput
	err := proto.Unmarshal(input.Data, &signIn)
	if err != nil {
		logger.Sugar.Error(err)
		ctx.Release()
		return
	}

	_, err = rpc_cli.LogicIntClient.SignIn(model.ContextWithRequestId(context.TODO(), input.RequestId), &pb.SignInReq{
		AppId:    signIn.AppId,
		UserId:   signIn.UserId,
		DeviceId: signIn.DeviceId,
		Token:    signIn.Token,
		ConnAddr: config.ConnConf.LocalAddr,
	})

	ctx.Output(pb.SessionType_SignIn, input.RequestId, err, nil)
	if err != nil {
		ctx.Release()
		return
	}

	ctx.AppId = signIn.AppId
	ctx.UserId = signIn.UserId
	ctx.DeviceId = signIn.DeviceId
	ctx.IsSignIn = true

	// 断开这个设备之前的连接
	preCtx := load(ctx.DeviceId)
	if preCtx != nil {
		preCtx.DeviceId = PreConn
	}

	store(ctx.DeviceId, ctx)
}

// Sync 消息同步
func (*handler) Sync(ctx *TcpConn, input pb.Input) {
	var sync pb.SyncInput
	err := proto.Unmarshal(input.Data, &sync)
	if err != nil {
		logger.Sugar.Error(err)
		ctx.Release()
		return
	}

	resp, err := rpc_cli.LogicIntClient.Sync(model.ContextWithRequestId(context.TODO(), input.RequestId), &pb.SyncReq{
		AppId:    ctx.AppId,
		UserId:   ctx.UserId,
		DeviceId: ctx.DeviceId,
		Seq:      sync.Seq,
	})

	var message proto.Message
	if err == nil {
		message = &pb.SyncOutput{Messages: resp.Messages}
	}
	ctx.Output(pb.SessionType_Sync, input.RequestId, err, message)
}

// Heartbeat 心跳
func (*handler) Heartbeat(ctx *TcpConn, input pb.Input) {
	ctx.Output(pb.SessionType_Heartbeat, input.RequestId, nil, nil)
	logger.Sugar.Infow("heartbeat", "device_id", ctx.DeviceId, "user_id", ctx.UserId)
}

// MessageACK 消息收到回执
func (*handler) MessageACK(ctx *TcpConn, input pb.Input) {
	var messageACK pb.MessageACK
	err := proto.Unmarshal(input.Data, &messageACK)
	if err != nil {
		logger.Sugar.Error(err)
		ctx.Release()
		return
	}

	_, _ = rpc_cli.LogicIntClient.MessageACK(model.ContextWithRequestId(context.TODO(), input.RequestId), &pb.MessageACKReq{
		AppId:       ctx.AppId,
		UserId:      ctx.UserId,
		DeviceId:    ctx.DeviceId,
		DeviceAck:   messageACK.DeviceAck,
		//Rec: messageACK.ReceiveTime,
	})
}

// Offline 设备离线
func (*handler) Offline(ctx *TcpConn) {
	_, _ = rpc_cli.LogicIntClient.Offline(context.TODO(), &pb.OfflineReq{
		AppId:    ctx.AppId,
		UserId:   ctx.UserId,
		DeviceId: ctx.DeviceId,
	})
}
