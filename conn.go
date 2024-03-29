package irc

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

const (
	DEBUG = true
)

type Conn struct {
	serverConn net.Conn
	info       user
	recv       <-chan Message
	sending    *sync.Mutex
}

func Connect(server string, info user) (c Conn, err error) {
	if DEBUG {
		fmt.Printf("Connecting to %s.\n", server)
	}
	stream, err := net.Dial("tcp", server)
	if DEBUG {
		fmt.Printf("Connected.\n")
	}
	if err != nil {
		return c, err
	}
	go func() {
		c.Send(info.PassMessage())
		c.Send(info.NickMessage())
		c.Send(info.UserMessage())
	}()
	send := make(chan Message)
	recv := make(chan Message, 5)
	c = c.init(stream, info, recv)
	go handle(stream, recv, send)
	return c, nil
}

func (c Conn) init(conn net.Conn, info user, recv chan Message) Conn {
	c.serverConn = conn
	c.info = info
	c.recv = recv
	var lock sync.Mutex
	c.sending = &lock
	return c
}

func (c Conn) Recv() <-chan Message {
	return c.recv
}

func (c Conn) Send(msg Message) {
	c.sending.Lock()
	defer c.sending.Unlock()
	if DEBUG {
		fmt.Printf("->"+msg.Tmpl(), msg.Data()...)
	}
	_, err := fmt.Fprintf(c.serverConn, msg.Tmpl(), msg.Data()...)
	if err != nil {
		panic(err)
	}
}

func handle(conn net.Conn, recv, send chan Message) {
	io := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	go process(io, recv)
}

func (c Conn) Write(data []byte) (int, error) {
	c.sending.Lock()
	defer c.sending.Unlock()
	return c.serverConn.Write(data)
}
