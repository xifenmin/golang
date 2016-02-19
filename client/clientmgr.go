package main

import (
	"fmt"
	"net"
	"strconv"
)

type ClientConnPoolError struct {
	desc string
}

func (e ClientConnPoolError) Error() string {
	return fmt.Sprintf("ClientConnPool error:%s", e.desc)
}

var maxClientConnerror = ClientConnPoolError{"Over maximum connection!!!"}

type ClientConn struct {
	conn net.Conn
	data chan interface{}
}

type ClientConnPool struct {
	ip    string
	port  uint
	size  uint
	count uint
	pool  []*ClientConn
	free  []*ClientConn
}

func NewClientPool(ip string, port uint, size uint) *ClientConnPool {
	return &ClientConnPool{
		ip:   ip,
		port: port,
		size: size,
	}
}

func (this *ClientConnPool) PutConn(clientConn *ClientConn) error {
	this.free = append(this.free, clientConn)
	return nil
}

func (this *ClientConnPool) SendData(data []byte) {

	for i := 0; i < len(this.pool); i++ {
		datalen := this.pool[i].SendData(data)
		if datalen <= 0 {
			fmt.Println("send fail\n")
			return
		}
	}
}

func (this *ClientConnPool) GetConn() (*ClientConn, error) {

	if this.count == this.size && len(this.free) == 0 {
		return nil, maxClientConnerror
	}

	if len(this.free) > 0 {
		conn := this.free[0]
		this.free = this.free[1:]
		return conn, nil
	}

	ipAddr := this.ip + ":" + strconv.FormatUint(uint64(this.port), 10)
	conn, err := NewClientConnection(ipAddr)
	if err != nil {
		return nil, err
	}

	this.count++

	this.pool = append(this.pool, conn)
	return conn, nil
}

func NewClientConnection(ipAddr string) (*ClientConn, error) {
	conn, err := net.Dial("tcp", ipAddr)

	if err != nil {
		fmt.Println("client create fail\n")
		return nil, err
	}
	clientConn := ClientConn{}

	if err == nil {
		clientConn.conn = conn
	}

	clientConn.data = make(chan interface{})
	return &clientConn, err
}

func (this *ClientConn) SendData(data []byte) int {
	datalen, err := this.conn.Write(data)
	if err != nil {
		fmt.Printf("send data fail:%s\n", this.conn.RemoteAddr())
		return 0
	}
	return datalen
}

func (this *ClientConn) Close() error {
	err := this.conn.Close()
	return err
}
