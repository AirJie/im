package tcp_conn

import (
	"github.com/aijie/michat/server/logger"
	"net"
)

type TcpServer struct {
	Addr            string
	AcceptGoroutine int
}

func NewTcpServer(addr string, routineNum int) *TcpServer {
	return &TcpServer{addr, routineNum}
}

func (t *TcpServer) Start() {
	addr, err := net.ResolveTCPAddr("tcp", t.Addr)
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	for i := 0; i < t.AcceptGoroutine; i++ {
		go t.Accept(listener)
	}
}

func (t *TcpServer) Accept(listen *net.TCPListener) {
	conn, err := listen.AcceptTCP()
	if err != nil {
		logger.Sugar.Error("Accept ", err)
		return
	}
	conn.SetKeepAlive(true)
	tcpConn := NewTcpConn(conn)
	go tcpConn.DoConn()
}
