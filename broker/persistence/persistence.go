package persistence

import (
	"github.com/crtrpt/mqtt/broker/system"
)

const (

	// KSubscription is the key for subscription data.
	KSubscription = "sub"

	// KServerInfo is the key for server info data.
	KServerInfo = "srv"

	// KRetained is the key for retained messages data.
	KRetained = "ret"

	// KInflight is the key for inflight messages data.
	KInflight = "ifm"

	// KClient is the key for client data.
	KClient = "cl"
)

// 持久化接口
type Store interface {
	//打开新的存储
	Open() error
	//关闭存储
	Close()
	//读取全部订阅
	ReadSubscriptions() (v []Subscription, err error)
	//写入订阅
	WriteSubscription(v Subscription) error
	//删除订阅
	DeleteSubscription(id string) error

	//读取客户端列表
	ReadClients() (v []Client, err error)
	//存储一个客户端
	WriteClient(v Client) error
	DeleteClient(id string) error

	// 读取飞行中的消息
	ReadInflight() (v []Message, err error)
	//把写入的消息存储起来
	WriteInflight(v Message) error
	//删除处理中的信息
	DeleteInflight(id string) error
	//设置飞行时长
	SetInflightTTL(seconds int64)
	//清除过期消息
	ClearExpiredInflight(expiry int64) error

	//读取broker信息
	ReadServerInfo() (v ServerInfo, err error)
	//写入broker信息
	WriteServerInfo(v ServerInfo) error

	//读取保留信息
	ReadRetained() (v []Message, err error)
	//写入保留信息
	WriteRetained(v Message) error
	//删除保留信息
	DeleteRetained(id string) error
}

// 服务器的统计信息
type ServerInfo struct {
	system.Info        // 系统信息
	ID          string // 存储key.
}

// Subscription contains the details of a topic filter subscription.
type Subscription struct {
	ID     string // 存储key.
	T      string // 存储类型
	Client string // the id of the client who the subscription belongs to.
	Filter string // the topic filter being subscribed to.
	QoS    byte   // 订阅的qos质量
}

// 保留消息或者正在处理的消息
type Message struct {
	Payload     []byte      // 消息体
	FixedHeader FixedHeader // 消息头
	T           string      // 持久化数据类型
	ID          string      // 存储key
	Client      string      // 表示这条消息是谁发送的.
	TopicName   string      // 消息要发送到的主题
	Created     int64       // 消息创建时间.
	Sent        int64       // 最后发送时间
	Resends     int         // 重试次数
	PacketID    uint16      // 消息id.
}

// mqtt固定头消息
type FixedHeader struct {
	Remaining int  // 消息体长度.
	Type      byte // 消息类型(PUBLISH, SUBSCRIBE, etc) from bits 7 - 4 (byte 1).
	Qos       byte // 服务质量.
	Dup       bool // 是否重复发送.
	Retain    bool // 消息是否保留.
}

// 客户端可以存储的消息
type Client struct {
	LWT      LWT    // 客户端的遗言.
	Username []byte // 用户名.
	ID       string // 持久化的id.
	ClientID string // client id.
	T        string // 存储的类型.
	Listener string // tcp / websocket ?
}

// 客户端的遗言数据
type LWT struct {
	Message []byte // 客户端断开的时候发送的消息体.
	Topic   string // 客户端断开的时候发送的目标topic
	Qos     byte   // 服务质量
	Retain  bool   // 消息是否保留
}
