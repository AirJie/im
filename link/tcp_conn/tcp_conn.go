package tcp_conn

import "net"

type TCPServer struct {
	Address            string
	AcceptGoroutineNum int
}

type TcpConn struct {
	Condec *Codec
}

func (c *TcpConn)DoConn() {
	//defer panicRecover
	for {

	}
}

