package main

import (
	"github.com/aijie/michat/config"
	"github.com/aijie/michat/datas/repository/redis"
	"github.com/aijie/michat/server/rpc_cli"
)

func main() {
	redis.InitDb()
	rpc_cli.InitLogicIntClient(config.LogicConf.ConnRPCAddrs)

}

