package main

import (
	"github.com/aijie/michat/api/ws_conn"
	ws "github.com/aijie/michat/link/ws_conn"
	"github.com/aijie/michat/config"
	"github.com/aijie/michat/server/rpc_cli"
	utils "github.com/aijie/michat/utils/panic"
)

func main() {
	go func() {
		defer utils.RecoverPanic()
		ws_conn.StartRPCServer()
	}()
	rpc_cli.InitLogicIntClient(config.WSConf.LogicRPCAddrs)
	ws.StartWsServer(config.WSConf.WSListenAddr)
}
