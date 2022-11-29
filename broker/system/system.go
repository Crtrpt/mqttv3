package system

// 服务器的当前统计信息
// 通常发送到 $SYS 主题
type Info struct {
	Version             string `json:"version"`              // 当前服务器的版本
	Started             int64  `json:"started"`              // 启动时间.
	Uptime              int64  `json:"uptime"`               // 在线时长.
	BytesRecv           int64  `json:"bytes_recv"`           // 总共接收到多少包.
	BytesSent           int64  `json:"bytes_sent"`           // 发送了多少包.
	ClientsConnected    int64  `json:"clients_connected"`    // 当前客户端连接数
	ClientsDisconnected int64  `json:"clients_disconnected"` // 断开客户端总数.
	ClientsMax          int64  `json:"clients_max"`          // 客户端最大并发数.
	ClientsTotal        int64  `json:"clients_total"`        // 连接数+断开连接的数量
	ConnectionsTotal    int64  `json:"connections_total"`    // 已经连接的客户端数量
	MessagesRecv        int64  `json:"messages_recv"`        // 收到多少条消息.
	MessagesSent        int64  `json:"messages_sent"`        // 发送多少条消息.
	PublishDropped      int64  `json:"publish_dropped"`      // 丢弃多少条消息.
	PublishRecv         int64  `json:"publish_recv"`         // 发布了多少条消息.
	PublishSent         int64  `json:"publish_sent"`         // 发布的消息多少条被接收了.
	Retained            int64  `json:"retained"`             // 当前保留了多少条消息.
	Inflight            int64  `json:"inflight"`             // 多少条消息正在处理.
	Subscriptions       int64  `json:"subscriptions"`        // 过滤的订阅总数.
}
