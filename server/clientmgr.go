package main

import (
	//"bytes"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const TIME_OUT = 30

type ClientConnPoolError struct {
	desc string
}

func (e ClientConnPoolError) Error() string {
	return fmt.Sprintf("ClientConnPool error:%s", e.desc)
}

var maxClientConnerror = ClientConnPoolError{"Over maximum connection!!!"}

type ClientConn struct {
	conn     net.Conn
	lastTime int64
	rwMutex  sync.RWMutex
	data     chan interface{}
}

type ClientConnPool struct {
	size  int64
	count int64
	mutex sync.Mutex
	pool  map[string]*ClientConn
	free  []*ClientConn
}

func NewClientPool(size int64) *ClientConnPool {

	return &ClientConnPool{
		size: size,
		pool: make(map[string]*ClientConn),
	}
}

func (this *ClientConnPool) CheckConnTimeOut() {

	for {
		this.mutex.Lock()
		for key, clientConn := range this.pool {
			now := time.Now().Unix()
			this.mutex.Lock()
			if now-clientConn.lastTime > TIME_OUT {
				log.Printf("Client Conn TimeOut:%s,Close!!!", this.pool[key].conn.RemoteAddr().String())
				clientConn.Close()
				delete(this.pool, key)
			}
		}
		this.mutex.Unlock()
		time.Sleep(time.Second * 1)
	}
}

func (this *ClientConnPool) PutConn(clientConn *ClientConn) error {
	this.mutex.Lock()
	this.free = append(this.free, clientConn)
	this.mutex.Unlock()
	return nil
}

func (this *ClientConnPool) SetConnTomap(clientConn *ClientConn) {
	this.pool[clientConn.conn.RemoteAddr().String()] = clientConn
	clientConn.lastTime = time.Now().Unix()
}

func (this *ClientConnPool) GetConn() (*ClientConn, error) {

	if this.count == this.size && len(this.free) == 0 {
		return nil, maxClientConnerror
	}

	if len(this.free) > 0 {
		this.mutex.Lock()
		conn := this.free[0]
		this.free = this.free[1:]
		this.mutex.Unlock()
		fmt.Println("Get conn from Pool\n")
		return conn, nil
	}

	clientConn, err := NewClientConnection()

	if err != nil {
		return nil, err
	}

	this.count++
	return clientConn, nil
}

func NewClientConnection() (*ClientConn, error) {
	clientConn := &ClientConn{}
	clientConn.data = make(chan interface{}, 10)
	clientConn.lastTime = time.Now().Unix()

	return clientConn, nil
}

func (this *ClientConn) Close() error {
	err := this.conn.Close()
	return err
}

func (this *ClientConn) getHex(data []byte, datalen int) string {

	var str string

	for i := 0; i < datalen; i++ {
		str += fmt.Sprintf("%.2x ", data[i])
	}

	return str
}

func (this *ClientConn) read(pool *ClientConnPool) {

	buf := make([]byte, 4096)
	var str string

	for {

		datalen, err := this.conn.Read(buf)
		if err != nil {
			log.Printf("Close socket %s,%s\n", time.Now().Format("2006-01-02 15:04:05"), this.conn.RemoteAddr().String())
			this.Close()
			pool.PutConn(this)
			break
		}

		//buffer := bytes.NewBuffer(buf)
		this.rwMutex.Lock()
		this.lastTime = time.Now().Unix()
		this.rwMutex.Unlock()

		str = this.conn.RemoteAddr().String() + "-----" + this.getHex(buf, datalen)
		this.data <- str
	}
}

func (this *ClientConn) readChannel() {

	for {
		_, ok := <-this.data

		if !ok {
			log.Println("Channel is Close!!\n")
			break
		}

		select {
		case msg := <-this.data:
			log.Printf("Read Data from Channel:%s,%s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
		}
	}
}

func (this *ClientConn) Run(pool *ClientConnPool) {
	go this.read(pool)
	go this.readChannel()
}
