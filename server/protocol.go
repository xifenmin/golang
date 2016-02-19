package main

import ()

type IMessage interface {
	toData() []byte
	fromData(buffer []byte) bool
}

type Message struct {
	magic   int
	seq     int
	version int
	body    interface{}
}

func (message *Message) toData() []byte {
	if message.body != nil {
		if m, ok := message.body.(IMessage); ok {
			return m.toData()
		}
	}
	return nil
}
