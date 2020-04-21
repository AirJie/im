package main

import (
	"fmt"
	"github.com/aijie/michat/datas/pb"
	"github.com/aijie/michat/link/tcp_conn"
	"github.com/aijie/michat/server/logger"
	"github.com/aijie/michat/utils/aes"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

const (
	LenLen               = 2
	tcpCliReadCtxMaxLen  = 65536
	tcpCliWriteCtxMaxLen = 1024
)

var (
	codecFactory = tcp_conn.NewCodecFactory(LenLen, tcpCliReadCtxMaxLen, tcpCliWriteCtxMaxLen)
)

func main() {
	for i := 0; i < 1000; i++ {
		client := TcpClient{
			AppId:    1,
			UserId:   int64(i),
			DeviceId: int64(i),
			Seq:      0,
			codec:    nil,
		}
		//client.S
	}
}

type TcpClient struct {
	AppId    int64
	UserId   int64
	DeviceId int64
	Seq      int64
	codec    *tcp_conn.Codec
}

func (c *TcpClient) Start() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	c.codec = codecFactory.GetCodec(conn)

}

func (c *TcpClient) SignIn() {
	token, err := aes.GetToken(c.AppId, c.UserId, c.DeviceId, time.Now().Add(24*30*time.Hour).Unix(), aes.PublicKey)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	signIn := pb.SignInReq{
		UserId:   c.UserId,
		AppId:    c.AppId,
		DeviceId: c.DeviceId,
		Token:    token,
	}
	c.Output(pb.SessionType_SignIn, time.Now().UnixNano(), &signIn)
}

func (c *TcpClient) SyncTrigger() {
	c.Output(pb.SessionType_Sync, time.Now().UnixNano(), &pb.SyncInput{Seq: c.Seq})
}

func (c *TcpClient) Heartbeat() {
	ticker := time.NewTicker(time.Minute * 5)
	for range ticker.C {
		c.Output(pb.SessionType_Heartbeat, time.Now().UnixNano(), nil)
	}
}

func (c *TcpClient) Output(pt pb.SessionType, requestId int64, message proto.Message) {
	var input = pb.Input{
		Type:      pt,
		RequestId: requestId,
	}

	if message != nil {
		bytes, err := proto.Marshal(message)
		if err != nil {
			fmt.Println(err)
			return
		}
		input.Data = bytes
	}

	inputByf, err := proto.Marshal(&input)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = c.codec.Encode(inputByf, time.Second)
	if err != nil {
		fmt.Println(err)
	}
}
