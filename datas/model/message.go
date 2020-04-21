package model

import (
	"github.com/aijie/michat/datas/pb"
	"time"
)

type Message struct {
	Id             int64     // 自增主键
	AppId          int64     // appId
	ObjectType     int       // 所属类型
	ObjectId       int64     // 所属类型id
	RequestId      int64     // 请求id
	SenderType     int32     // 发送者类型
	SenderId       int64     // 发送者账户id
	SenderDeviceId int64     // 发送者设备id
	ReceiverType   int32     // 接收者账户id
	ReceiverId     int64     // 接收者id,如果是单聊信息，则为user_id，如果是群组消息，则为group_id
	ToUserIds      string    // 需要@的用户id列表，多个用户用，隔开
	Type           int       // 消息类型
	Content        string    // 消息内容
	Seq            int64     // 消息同步序列
	SendTime       time.Time // 消息发送时间
	Status         int32     // 创建时间
}

type SendMessage struct {
	ReceiverType pb.ReceiveType `json:"receiver_type"`
	ReceiverId   int64          `json:"receiver_id"`
	ToUserIds    []int64        `json:"to_user_ids"`
	MessageId    string         `json:"message_id"`
	SendTime     int64          `json:"send_time"`
	MessageBody  struct {
		MessageType    int    `json:"message_type"`
		MessageContent string `json:"-"`
	} `json:"message_body"`
	PbBody *pb.MessageBody `json:"-"`
}
