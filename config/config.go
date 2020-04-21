package config

var (
	LogicConf logicConf
	ConnConf  connConf
	WSConf    wsConf
)

type logicConf struct {
	MySQL                  string
	NSQIP                  string
	RedisIp                string
	RPCIntListenAddr       string
	ClientRPCExtListenAddr string
	ServerRPCExtListenAddr string
	ConnRPCAddrs           string
}

type connConf struct {
	TCPListenAddr string
	RPCListenAddr string
	LocalAddr     string
	LogicRPCAddrs string
}

type wsConf struct {
	WSListenAddr  string
	RPCListenAddr string
	LocalAddr     string
	LogicRPCAddrs string
}

func init() {
	//env := os.Getenv("gim_env")
	initDevConfig()
}
