package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"log"

	"github.com/gorilla/websocket"
)

// WsReadFull 读指定长度
func WsReadFull(conn *websocket.Conn) (buf []byte, err error) {
	length := 0
	size := len(buf)
	for {
		conn.SetReadLimit(int64(size - length))
		_, b, e := conn.ReadMessage()
		if err != nil {
			err = e
			log.Print("WsReadFull : ", err)
			return
		}
		n := len(b)
		//copy(buf[length:], n)
		buf = append(buf, b[:]...)
		length += n
		if length == size {
			return
		}
	}
}
