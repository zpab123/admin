////////////////////////////////////////////////////////////////////////////////
// mqtt 公用模块

package mqtt

import (
	"net"
	"sync"
	"sync/atomic"

	proto "github.com/huin/mqtt" // mqtt 协议
)

////////////////////////////////////////////////////////////////////////////////
// 私有对象 -- 服务器状态

// 服务器状态
type stats struct {
	recv       int64 // 接收消息状态
	sent       int64 // 发送消息状态
	clients    int64 // 是否处于连接状态 1= 连接 -1 = 断开
	clientsMax int64
	lastmsgs   int64
}

// 设置为：接收消息状态
func (s *stats) messageRecv() {
	atomic.AddInt64(&s.recv, 1)
}

// 设置为：发送消息状态
func (s *stats) messageSend() {
	atomic.AddInt64(&s.sent, 1)
}

// 设置为： 连接状态
func (s *stats) clientConnect() {
	atomic.AddInt64(&s.clients, 1)
}

// 设置为： 断开状态
func (s *stats) clientDisconnect() {
	atomic.AddInt64(&s.clients, -1)
}

////////////////////////////////////////////////////////////////////////////////
// 私有对象 -- 连接对象

// 接收通道
type receipt chan struct{}

type job struct {
	m proto.Message
	r receipt
}

////////////////////////////////////////////////////////////////////////////////
// 私有对象 -- 连接对象

// An IncomingConn represents a connection into a Server.
type incomingConn struct {
	svr      *Server
	conn     net.Conn
	jobs     chan job
	clientid string
	Done     chan struct{}
}

////////////////////////////////////////////////////////////////////////////////
// 私有对象 -- 消息对象： 订阅

type subscriptions struct {
	workers   int
	posts     chan post
	mu        sync.Mutex // guards access to fields below
	subs      map[string][]*incomingConn
	wildcards []wild
	retain    map[string]retain
	stats     *stats
}
