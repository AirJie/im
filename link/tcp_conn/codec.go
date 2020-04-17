package tcp_conn

import (
	"net"
	"sync"
)

type Codec struct {
	f       *CodecFactory
	Conn    net.Conn
	ReadBuf buffer
}

type CodecFactory struct {
	ContextLen         int
	ReadContextMaxLen  int
	WriteContextMaxLen int
	ReadBufferPool     sync.Pool
	WriteBufferPool    sync.Pool
}

func NewCodecFactory(len, readContextMaxLen, writeContextMaxLen int) CodecFactory {
	return CodecFactory{
		ContextLen:         len,
		ReadContextMaxLen:  readContextMaxLen,
		WriteContextMaxLen: writeContextMaxLen,
		ReadBufferPool: sync.Pool{
			New: func() interface{} {
				b := make([]byte, readContextMaxLen+len)
				return b
			},
		},
		WriteBufferPool: sync.Pool{
			New: func() interface{} {
				b := make([]byte, writeContextMaxLen+len)
				return b
			},
		},
	}
}
