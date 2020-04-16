package ws_conn

import (
	"context"
	"github.com/aijie/michat/datas/pb"
	"github.com/aijie/michat/logger"
)

func DeliverMessage(ctx context.Context, req *pb.DeliverMessageReq) error {
	conn := load(req.DeviceId)
	if conn == nil {
		logger.Sugar.Warn("conn id not found")
		return nil
	}
	conn.Output(pb.SessionType_MessageStream, grpclib.)
}
