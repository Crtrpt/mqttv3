package events

import (
	"github.com/crtrpt/mqtt/broker/internal/packets"
)

// 订阅消息回调接口
type Events struct {
	OnProcessMessage // 消息发布之前
	OnMessage        // 收到已发布的消息.
	OnError          // 服务器异常.
	OnConnect        // 客户端已连接
	OnDisconnect     // 客户端断开连接.
	OnSubscribe      // 订阅主题.
	OnUnsubscribe    // 取消订阅主题
}

// 包别名
type Packet packets.Packet

// 已连接的客户端的信息
type Client struct {
	ID           string
	Remote       string
	Listener     string
	Username     []byte
	CleanSession bool
}

// Clientlike is an interface for Clients and client-like objects that
// are able to describe their client/listener IDs and remote address.
type Clientlike interface {
	Info() Client
}

// OnProcessMessage is called when a publish message is received, allowing modification
// of the packet data after ACL checking has occurred but before any data is evaluated
// for processing - e.g. for changing the Retain flag. Note, this hook is ONLY called
// by connected client publishers, it is not triggered when using the direct
// s.Publish method. The function receives the sent message and the
// data of the client who published it, and allows the packet to be modified
// before it is dispatched to subscribers. If no modification is required, return
// the original packet data. If an error occurs, the original packet will
// be dispatched as if the event hook had not been triggered.
// This function will block message dispatching until it returns. To minimise this,
// have the function open a new goroutine on the embedding side.
// The `mqtt.ErrRejectPacket` error can be returned to reject and abandon any further
// processing of the packet.
// 发布消息到达broker的时候触发
type OnProcessMessage func(Client, Packet) (Packet, error)

// OnMessage function is called when a publish message is received. Note,
// this hook is ONLY called by connected client publishers, it is not triggered when
// using the direct s.Publish method. The function receives the sent message and the
// data of the client who published it, and allows the packet to be modified
// before it is dispatched to subscribers. If no modification is required, return
// the original packet data. If an error occurs, the original packet will
// be dispatched as if the event hook had not been triggered.
// This function will block message dispatching until it returns. To minimise this,
// have the function open a new goroutine on the embedding side.
type OnMessage func(Client, Packet) (Packet, error)

// 连接的时候回调
type OnConnect func(Client, Packet)

// 断开连接的时候回调
type OnDisconnect func(Client, error)

// 发生错误的时候回调
type OnError func(Client, error)

// 订阅的时候回调
type OnSubscribe func(filter string, cl Client, qos byte)

// 取消订阅的时候回调
type OnUnsubscribe func(filter string, cl Client)
