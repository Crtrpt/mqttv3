package listeners

import (
	"crypto/tls"
	"net"
	"sync"
	"sync/atomic"

	"github.com/crtrpt/mqtt/broker/listeners/auth"
	"github.com/crtrpt/mqtt/broker/system"
)

type MQTTv5 struct {
	sync.RWMutex
	id       string       // the internal id of the listener.
	protocol string       // the MQTTv5 protocol to use.
	address  string       // the network address to bind to.
	listen   net.Listener // a net.Listener which will listen for new clients.
	config   *Config      // configuration values for the listener.
	end      uint32       // ensure the close methods are only called once.
}

func NewMQTTv5(id, address string) *MQTTv5 {
	return &MQTTv5{
		id:       id,
		protocol: "tcp",
		address:  address,
		config: &Config{
			Auth: new(auth.Allow),
			TLS:  new(TLS),
		},
	}
}

func (l *MQTTv5) SetConfig(config *Config) {
	l.Lock()
	if config != nil {
		l.config = config
		if l.config.Auth == nil {
			l.config.Auth = new(auth.Disallow)
		}
	}

	l.Unlock()
}

// ID returns the id of the listener.
func (l *MQTTv5) ID() string {
	l.RLock()
	id := l.id
	l.RUnlock()
	return id
}

// 监听网络
func (l *MQTTv5) Listen(s *system.Info) error {
	var err error

	if l.config.TLS != nil && len(l.config.TLS.Certificate) > 0 && len(l.config.TLS.PrivateKey) > 0 {
		var cert tls.Certificate
		cert, err = tls.X509KeyPair(l.config.TLS.Certificate, l.config.TLS.PrivateKey)
		if err != nil {
			return err
		}

		l.listen, err = tls.Listen(l.protocol, l.address, &tls.Config{
			Certificates: []tls.Certificate{cert},
		})
	} else if l.config.TLSConfig != nil {
		l.listen, err = tls.Listen(l.protocol, l.address, l.config.TLSConfig)
	} else {
		l.listen, err = net.Listen(l.protocol, l.address)
	}

	if err != nil {
		return err
	}

	return nil
}

// 等待连接建立
func (l *MQTTv5) Serve(establish EstablishFunc) {
	for {
		if atomic.LoadUint32(&l.end) == 1 {
			return
		}

		conn, err := l.listen.Accept()
		if err != nil {
			return
		}

		if atomic.LoadUint32(&l.end) == 0 {
			go func() {
				_ = establish(l.id, conn, l.config.Auth)
			}()
		}
	}
}

// 关闭服务器和所有的客户端连接
func (l *MQTTv5) Close(closeClients CloseFunc) {
	l.Lock()
	defer l.Unlock()

	if atomic.CompareAndSwapUint32(&l.end, 0, 1) {
		closeClients(l.id)
	}

	if l.listen != nil {
		err := l.listen.Close()
		if err != nil {
			return
		}
	}
}
