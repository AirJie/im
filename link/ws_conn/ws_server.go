package ws_conn

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aijie/michat/datas/model"
	"github.com/aijie/michat/datas/pb"
	"github.com/aijie/michat/server/errorcode"
	"github.com/aijie/michat/server/logger"
	"github.com/aijie/michat/server/rpc_cli"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 65536,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	appId, _ := strconv.ParseInt(r.Header.Get(model.CtxAppId), 10, 64)
	userId, _ := strconv.ParseInt(r.Header.Get(model.CtxUserId), 10, 64)
	deviceId, _ := strconv.ParseInt(r.Header.Get(model.CtxDeviceId), 10, 64)
	token := r.Header.Get(model.CtxToken)
	requestId, _ := strconv.ParseInt(r.Header.Get(model.CtxRequestId), 10, 64)
	if appId == 0 || userId == 0 || deviceId == 0 || token == "" || requestId == 0 {
		s, _ := status.FromError(errorcode.ErrUnauthorized)
		bytes, err := json.Marshal(s)
		if err != nil {
			logger.Sugar.Error(err)
			return
		}
		w.Write(bytes)
	}

	sign := pb.SignInReq{
		AppId:    appId,
		UserId:   userId,
		DeviceId: deviceId,
		Token:    token,
	}
	_, err := rpc_cli.LogicIntClient.SignIn(model.ContextWithReqId(context.TODO(), requestId), &sign)
	s, _ := status.FromError(err)
	if s.Code() == codes.OK {
		bytes, err := proto.Marshal(s.Proto())
		if err != nil {
			logger.Sugar.Error(err)
			return
		}
		w.Write(bytes)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	if link := load(deviceId); link != nil {
		link.DeviceId = -1
	}
	ctx := NewLink(conn, appId, userId, deviceId)
	store(deviceId, ctx)
	ctx.DoConn()
}

func StartWsServer(addr string) {
	http.HandleFunc("/ws", wsHandler)
	logger.Sugar.Info("ws server start")
	fmt.Println(http.ListenAndServe(addr, nil))
}
