package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

const (
	maxConnectionNum = 100000
)

type Server struct {
	ip       string
	port     int
	listener net.Listener
}

func NewServer(ip string, port int) *Server {
	return &Server{
		ip:   ip,
		port: port,
	}
}

func (this *Server) StartServer() {
	fmt.Println("Init Server\n")
	this.listen()
}

func (this *Server) handleClient(clientConn *ClientConn, pool *ClientConnPool) {
	clientConn.Run(pool)
}

func (this *Server) listen() bool {

	var listen_addr = this.ip + ":" + strconv.Itoa(int(this.port))
	var err error

	this.listener, err = net.Listen("tcp", listen_addr)

	if err != nil {
		fmt.Println("Init socket fail!\n", err.Error())
		return false
	}

	pool := NewClientPool(maxConnectionNum)
	go pool.CheckConnTimeOut()

	for {

		client, err := this.listener.Accept()

		if err != nil {
			return false
		}

		clientConn, err := pool.GetConn()

		if err != nil {
			fmt.Println("Get Client Obj from Connection Poll fail!!!\n")
			return false
		}

		log.Printf("one new connection:%s\n", client.RemoteAddr().String())

		clientConn.conn = client

		pool.SetConnTomap(clientConn)
		this.handleClient(clientConn, pool)
	}
}
