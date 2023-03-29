package mq

import (
	"github.com/betashepherd/stomp/v3"
	"time"
)

var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
	//stomp.ConnOpt.Login("admin", "admin"),
	//stomp.ConnOpt.Host("/"),
	stomp.ConnOpt.HeartBeat(time.Second, time.Second),
}

func createActiveMqConn(serverAddr string, options ...func(*stomp.Conn) error) (*stomp.Conn, error) {
	conn, err := stomp.Dial("tcp", serverAddr, options...)
	if err != nil {
		println("cannot connect to server", err.Error())
		return nil, err
	}

	return conn, nil
}

var ActiveMQ *stomp.Conn

func InitActiveMQ() {
	//ActiveMQ, _ = createActiveMqConn(activeMQ.Addr, options...)
}
