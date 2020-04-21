package tcp_conn

import (
	"github.com/aijie/michat/datas/pb"
	"github.com/aijie/michat/server/logger"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"io"
	"net"
	"strings"
	"time"
)

const (
	WriteDeadline = 10 * time.Minute
	ReadDeadline  = 10 * time.Minute

	// tcp buff const
	TypeLen = 2
	LenLen = 2
	ReadContextMaxLen = 252
	WriteContextMaxLen = 508
)

var (
	codecFactory = NewCodecFactory(LenLen, ReadContextMaxLen, WriteContextMaxLen)
)
type TcpConn struct {
	Codec    *Codec
	IsSignIn bool
	AppId    int64
	DeviceId int64
	UserId   int64
}

type Package struct {
	CodeType int
	Context  []byte
}

func NewTcpConn(conn net.Conn) *TcpConn{
	codec := codecFactory.GetCodec(conn)
	return &TcpConn{
		Codec: codec,
	}
}

func (c *TcpConn) DoConn() {
	//defer panicRecover
	for {
		err := c.Codec.Conn.SetReadDeadline(time.Now().Add(ReadDeadline))
		if err != nil {
			c.HandlerErr(err)
			return
		}

		_, err = c.Codec.Read()
		if err != nil {
			logger.Sugar.Error(err)
			return
		}
		for {
			bytes, ok, err := c.Codec.Decode()
			if err != nil {
				c.HandlerErr(err)
				return
			}
			if ok {
				Handler.Handler(c, bytes)
			}
		}
	}
}

func (c *TcpConn) HandlerErr(err error) {
	logger.Sugar.Debug("read tcp error", zap.Int64("appid", c.AppId), zap.Int64("userid", c.UserId))
	str := err.Error()
	if strings.HasSuffix(str, "use of closed net work") {
		return
	}
	c.Release()
	if err == io.EOF {
		return
	}
	if strings.HasSuffix(str, "i/o timeout") {
		return
	}
}

func (c *TcpConn) Release() {
	c.Codec.Release()
}

func (c *TcpConn) Output(ptype pb.SessionType, requestId int64, err error, message proto.Message) {

}
