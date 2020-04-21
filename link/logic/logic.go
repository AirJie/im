package logic

import (
	"github.com/aijie/michat/config"
	"github.com/aijie/michat/datas/pb"
	"github.com/aijie/michat/server/logger"
	"google.golang.org/grpc"
	"net"
)

func StartRPCServer() {
	go func() {
		listen, err := net.Listen("tcp", config.LogicConf.RPCIntListenAddr)
		if err != nil {
			logger.Sugar.Error(err)
			return
		}
		server := grpc.NewServer()
		pb.RegisterLogicClientExtServer()
	}()
}
