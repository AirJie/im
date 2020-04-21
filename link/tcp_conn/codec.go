package tcp_conn

import (
	"encoding/binary"
	"errors"
	"github.com/aijie/michat/server/logger"
	"go.uber.org/zap"
	"net"
	"sync"
	"time"
)

var (
	ErrOverReadLength = errors.New("reading length error")
)

type Codec struct {
	f       *CodecFactory
	Conn    net.Conn
	ReadBuf buffer
}

type CodecFactory struct {
	LenLen             int
	ReadContextMaxLen  int
	WriteContextMaxLen int
	ReadBufferPool     sync.Pool
	WriteBufferPool    sync.Pool
}

func NewCodecFactory(len, readContextMaxLen, writeContextMaxLen int) CodecFactory {
	return CodecFactory{
		LenLen:             len,
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

func (f *CodecFactory)GetCodec(conn net.Conn) *Codec {
	return &Codec{
		f: f,
		Conn: conn,
		ReadBuf: newBuffer(f.ReadBufferPool.Get().([]byte)),
	}
}

func (c *Codec) Read() (int, error) {
	return c.ReadBuf.loadFromReader(c.Conn)
}

func (c *Codec) Encode(bytes []byte, duration time.Duration) error {
	var buffer []byte
	if len(bytes) <= c.f.WriteContextMaxLen {
		cache := c.f.WriteBufferPool.Get().([]byte)
		buffer = cache[0 : c.f.LenLen+len(bytes)]
		defer c.f.WriteBufferPool.Put(cache)
	} else {
		buffer = make([]byte, c.f.LenLen+len(bytes))
	}
	//put len of context
	binary.BigEndian.PutUint16(buffer[0:c.f.LenLen], uint16(len(bytes)))
	copy(buffer[c.f.LenLen:], bytes)
	c.Conn.SetReadDeadline(time.Now().Add(duration))
	_, err := c.Conn.Write(buffer)
	if err != nil {
		return err
	}
	return nil
}

func (c *Codec) Decode() ([]byte, bool, error) {
	lenValue, err := c.ReadBuf.seek(0, c.f.LenLen)
	if err != nil {
		return nil, false, err
	}
	len := int(binary.BigEndian.Uint16(lenValue))
	if len > c.f.ReadContextMaxLen {
		logger.Logger.Error("decode reading error:", zap.Int("len:", len))
		return nil, false, ErrOverReadLength
	}
	data, err := c.ReadBuf.read(c.f.LenLen, len)
	if err != nil {
		return nil, false, err
	}
	return data, true, nil
}

func (c *Codec)Release() error {
	if err := c.Conn.Close(); err != nil {
		logger.Sugar.Error(err)
		return err
	}
	c.f.ReadBufferPool.Put(c.ReadBuf.buf)
	return nil
}
