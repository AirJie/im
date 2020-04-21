package api

import (
	"context"
	"errors"
	"github.com/aijie/michat/datas/model"
	"github.com/aijie/michat/datas/pb"
	"github.com/aijie/michat/server/errorcode"
	"github.com/aijie/michat/server/logger"
	"github.com/aijie/michat/ws"
	"github.com/go-zoo/bone"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

var (
	errUnauthorizedAccess = errors.New("missing or invalid credentials provided")
	errModelData          = errors.New("missing or invalid model data")
	errMalformedSubtopic  = errors.New("malformed subtopic")
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	logic pb.LogicIntClient
)

func MakeHandler(svc ws.Service) {
	r := bone.New()
	r.HandleFunc("/ws", imWsHandler)
}

func imWsHandler(w http.ResponseWriter, r *http.Request) {
	//1. do auth
	_, err := authorize(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	ctxType := contextType(r)
}

type subscription struct {
	conn  *websocket.Conn
}

func authorize(r *http.Request) (subscription, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		auths := bone.GetQuery(r, "Authorization")
		if len(auths) == 0 {
			return subscription{}, errUnauthorizedAccess
		}
		auth = auths[0]
	}
	token := r.Header.Get(model.CtxToken)
	appId, _ := strconv.ParseInt(r.Header.Get(model.CtxAppId), 10, 64)
	userId, _ := strconv.ParseInt(r.Header.Get(model.CtxUserId), 10, 64)
	deviceId, _ := strconv.ParseInt(r.Header.Get(model.CtxDeviceId), 10, 64)
	reqId, _ := strconv.ParseInt(r.Header.Get(model.CtxRequestId), 10, 64)
	if appId == 0 || userId == 0 || deviceId == 0 || token == "" || reqId == 0 {
		logger.Sugar.Error(errModelData)
		return subscription{}, errModelData
	}
	sign := pb.SignInReq{
		AppId: appId,
		UserId: userId,
		DeviceId: deviceId,
		Token: token,
	}
	_, err := logic.SignIn(model.ContextWithReqId(context.TODO(), reqId), &sign)
	if err != nil {
		e, ok := status.FromError(err)
		if ok && e.Code() == codes.PermissionDenied {
			return subscription{}, errorcode.ErrUnauthorized
		}
		return subscription{}, err
	}
	return subscription{}, nil
}

func contextType(r *http.Request) string {
	ct := r.Header.Get("Context-type")
	if ct == "" {
		vals := bone.GetQuery(r, "context-type")
		if len(vals) == 0 {
			return ""
		}
		return vals[0]
	}
	return ct
}
