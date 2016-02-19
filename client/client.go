package main

import (
	"time"
)

func main() {

	pool := NewClientPool("192.168.104.168", 5779, 100000)

	str := []byte("123456789.abcdefg")

	for i := 0; i < 20000; i++ {
		_, err := pool.GetConn()
		if err != nil {
			return
		}
	}

	for {
		pool.SendData(str)
		time.Sleep(time.Second * 1)
	}
}
