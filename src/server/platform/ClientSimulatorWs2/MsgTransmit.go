package main

//
// Created by YangZhi
// 			4/9/2019
//

import (
	"io"
	"net"
)

//MsgTransmit 消息传输接口
type MsgTransmit interface {
	//接收数据
	OnRecvMessage(peer Session) (msg interface{}, err error)
	//发送数据
	OnSendMessage(peer Session, msg interface{}) error
}

//
func IsEOFOrReadError(err error) bool {
	if err == io.EOF {
		return true
	}
	ne, ok := err.(*net.OpError)
	return ok && ne.Op == "read"
}

//
func IsEOFOrWriteError(err error) bool {
	if err == io.EOF {
		return true
	}
	ne, ok := err.(*net.OpError)
	return ok && ne.Op == "write"
}
