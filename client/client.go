package main

import (
	"time"
)

func main() {

	pool := NewClientPool("10.24.233.78", 18899, 70000)

	str := []byte("123456789.abcdefg")

	for i := 0; i < 2000; i++ {
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
